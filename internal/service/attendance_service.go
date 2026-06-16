package service

import (
	"errors"
	"time"

	"github.com/Touy2004/palm-back-end/internal/model"
	"github.com/Touy2004/palm-back-end/internal/repository"
	"github.com/google/uuid"
)

type AttendanceService struct {
	attendanceRepo *repository.AttendanceRepository
	palmRepo       *repository.PalmRepository
	deviceRepo     *repository.DeviceRepository
	userRepo       *repository.UserRepository
	cryptoSvc      *CryptoService
}

func NewAttendanceService(
	attendanceRepo *repository.AttendanceRepository,
	palmRepo *repository.PalmRepository,
	deviceRepo *repository.DeviceRepository,
	userRepo *repository.UserRepository,
	cryptoSvc *CryptoService,
) *AttendanceService {
	return &AttendanceService{
		attendanceRepo: attendanceRepo,
		palmRepo:       palmRepo,
		deviceRepo:     deviceRepo,
		userRepo:       userRepo,
		cryptoSvc:      cryptoSvc,
	}
}

type ProcessAttendanceInput struct {
	DeviceCode     string
	ModelVersion   string
	EmbeddingDim   int
	Embedding      []float32
	LivenessPassed bool
	QualityScore   float64
	ThermalMin     float64
	ThermalMax     float64
	ThermalAvg     float64
}

type AttendanceResult struct {
	Action   string
	UserID   uuid.UUID
	FullName string
	Score    float64
	Time     time.Time
	Message  string
}

func (s *AttendanceService) ProcessPalmAttendance(input ProcessAttendanceInput) (*AttendanceResult, error) {
	device, err := s.deviceRepo.FindByCode(input.DeviceCode)
	if err != nil || device.Status != "active" {
		return nil, errors.New("invalid or inactive device")
	}

	if !input.LivenessPassed {
		return nil, errors.New("liveness check failed")
	}

	templates, err := s.palmRepo.FindAllActive()
	if err != nil || len(templates) == 0 {
		return nil, errors.New("no active palm templates found")
	}

	bestScore := 0.0
	var matchedTemplate *model.PalmTemplate

	for i := range templates {
		t := &templates[i]

		decryptedBytes, err := s.cryptoSvc.Decrypt(t.TemplateEncrypted, t.TemplateNonce)
		if err != nil {
			continue // Skip if decryption fails
		}

		dbEmbedding, err := BytesToFloat32Slice(decryptedBytes)
		if err != nil {
			continue
		}

		score := CosineSimilarity(input.Embedding, dbEmbedding)
		if score > bestScore {
			bestScore = score
			matchedTemplate = t
		}
	}

	if matchedTemplate == nil || bestScore < matchedTemplate.Threshold {
		return nil, errors.New("palm not recognized")
	}

	user, err := s.userRepo.FindByID(matchedTemplate.UserID.String())
	if err != nil || user.Status != "active" {
		return nil, errors.New("user inactive")
	}

	// Determine check-in or check-out
	todayLog, err := s.attendanceRepo.FindTodayByUserID(user.ID.String())

	action := ""
	message := ""
	now := time.Now().UTC()

	loc, locErr := time.LoadLocation("Asia/Bangkok")
	if locErr != nil {
		loc = time.FixedZone("ICT", 7*3600)
	}
	localTime := now.In(loc)

	if err != nil { // No record today, so Check In
		action = "check_in"
		message = "Check-in success"

		status := "present"
		if localTime.Hour() > 8 || (localTime.Hour() == 8 && localTime.Minute() > 15) {
			status = "late"
		}

		newLog := &model.AttendanceLog{
			UserID:          user.ID,
			DeviceID:        &device.ID,
			AttendanceDate:  now,
			CheckInTime:     &now,
			CheckInScore:    &bestScore,
			CheckInLiveness: &input.LivenessPassed,
			Status:          status,
		}
		_ = s.attendanceRepo.Create(newLog)
	} else if todayLog.CheckOutTime == nil { // Already checked in, so Check Out
		action = "check_out"
		message = "Check-out success"

		todayLog.CheckOutTime = &now
		todayLog.CheckOutScore = &bestScore
		todayLog.CheckOutLiveness = &input.LivenessPassed
		_ = s.attendanceRepo.Update(todayLog)
	} else { // Already checked out
		return nil, errors.New("already completed today")
	}

	return &AttendanceResult{
		Action:   action,
		UserID:   user.ID,
		FullName: user.FullName,
		Score:    bestScore,
		Time:     now,
		Message:  message,
	}, nil
}

func (s *AttendanceService) IdentifyPalm(input ProcessAttendanceInput) (*AttendanceResult, error) {
	device, err := s.deviceRepo.FindByCode(input.DeviceCode)
	if err != nil || device.Status != "active" {
		return nil, errors.New("invalid or inactive device")
	}

	if !input.LivenessPassed {
		return nil, errors.New("liveness check failed")
	}

	templates, err := s.palmRepo.FindAllActive()
	if err != nil || len(templates) == 0 {
		return nil, errors.New("no active palm templates found")
	}

	bestScore := 0.0
	var matchedTemplate *model.PalmTemplate

	for i := range templates {
		t := &templates[i]

		decryptedBytes, err := s.cryptoSvc.Decrypt(t.TemplateEncrypted, t.TemplateNonce)
		if err != nil {
			continue
		}

		dbEmbedding, err := BytesToFloat32Slice(decryptedBytes)
		if err != nil {
			continue
		}

		score := CosineSimilarity(input.Embedding, dbEmbedding)
		if score > bestScore {
			bestScore = score
			matchedTemplate = t
		}
	}

	if matchedTemplate == nil || bestScore < matchedTemplate.Threshold {
		return nil, errors.New("palm not recognized")
	}

	user, err := s.userRepo.FindByID(matchedTemplate.UserID.String())
	if err != nil || user.Status != "active" {
		return nil, errors.New("user inactive")
	}

	return &AttendanceResult{
		Action:   "identify",
		UserID:   user.ID,
		FullName: user.FullName,
		Score:    bestScore,
		Time:     time.Now().UTC(),
		Message:  "User identified",
	}, nil
}
