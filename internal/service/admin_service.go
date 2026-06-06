package service

import (
	"fmt"
	"time"

	"github.com/Touy2004/palm-back-end/internal/model"
	"github.com/Touy2004/palm-back-end/internal/repository"
	"github.com/Touy2004/palm-back-end/pkg/hash"
)

type AdminService struct {
	adminRepo      *repository.AdminRepository
	userRepo       *repository.UserRepository
	deviceRepo     *repository.DeviceRepository
	attendanceRepo *repository.AttendanceRepository
	palmRepo       *repository.PalmRepository
}

func NewAdminService(
	adminRepo *repository.AdminRepository,
	userRepo *repository.UserRepository,
	deviceRepo *repository.DeviceRepository,
	attendanceRepo *repository.AttendanceRepository,
	palmRepo *repository.PalmRepository,
) *AdminService {
	return &AdminService{
		adminRepo:      adminRepo,
		userRepo:       userRepo,
		deviceRepo:     deviceRepo,
		attendanceRepo: attendanceRepo,
		palmRepo:       palmRepo,
	}
}

// User methods
func (s *AdminService) CreateUser(input RegisterInput) (*model.User, error) {
	hashed, err := hash.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		FullName:     input.FullName,
		EmployeeCode: input.EmployeeCode,
		Email:        input.Email,
		Department:   input.Department,
		Phone:        input.Phone,
		PasswordHash: hashed,
		Role:         input.Role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AdminService) GetUsers() ([]model.User, error) {
	return s.userRepo.FindAll()
}

func (s *AdminService) GetUserByID(id string) (*model.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *AdminService) UpdateUser(id string, input map[string]interface{}) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if v, ok := input["full_name"].(string); ok {
		user.FullName = v
	}
	if v, ok := input["employee_code"].(string); ok {
		user.EmployeeCode = v
	}
	if v, ok := input["email"].(string); ok {
		user.Email = v
	}
	if v, ok := input["department"].(string); ok {
		user.Department = v
	}
	if v, ok := input["phone"].(string); ok {
		user.Phone = v
	}
	if v, ok := input["role"].(string); ok {
		user.Role = v
	}
	if v, ok := input["status"].(string); ok {
		user.Status = v
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AdminService) DeleteUser(id string) error {
	return s.userRepo.Delete(id)
}

func (s *AdminService) SearchUsers(query string) ([]model.User, error) {
	return s.userRepo.Search(query)
}

// Device methods
func (s *AdminService) GetDevices() ([]model.Device, error) {
	return s.deviceRepo.FindAll()
}

func (s *AdminService) CreateDevice(device *model.Device) error {
	return s.deviceRepo.Create(device)
}

func (s *AdminService) UpdateDevice(id string, input map[string]interface{}) (*model.Device, error) {
	device, err := s.deviceRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if v, ok := input["device_code"].(string); ok {
		device.DeviceCode = v
	}
	if v, ok := input["device_name"].(string); ok {
		device.DeviceName = v
	}
	if v, ok := input["location_name"].(string); ok {
		device.LocationName = v
	}
	if v, ok := input["status"].(string); ok {
		device.Status = v
	}

	if err := s.deviceRepo.Update(device); err != nil {
		return nil, err
	}
	return device, nil
}

// Attendance methods
func (s *AdminService) GetAttendanceHistory(page, limit int, startDate, endDate string) ([]model.AttendanceLog, int64, error) {
	return s.attendanceRepo.FindAll(page, limit, startDate, endDate)
}

func (s *AdminService) GetUserAttendanceHistory(userID string, page, limit int, startDate, endDate string) ([]model.AttendanceLog, int64, error) {
	return s.attendanceRepo.FindByUserID(userID, page, limit, startDate, endDate)
}

func (s *AdminService) GetDashboardSummary() (*repository.DashboardSummary, error) {
	return s.adminRepo.GetDashboardSummary()
}

func (s *AdminService) GetUserPalmTemplates(userID string) ([]model.PalmTemplate, error) {
	return s.palmRepo.FindByUserID(userID)
}

func (s *AdminService) DeleteUserPalmTemplate(userID, templateID string) error {
	return s.palmRepo.Delete(templateID, userID)
}

func calculateWorkDays(start, end time.Time) int {
	days := 0
	current := start
	for !current.After(end) {
		if current.Weekday() != time.Saturday && current.Weekday() != time.Sunday {
			days++
		}
		current = current.AddDate(0, 0, 1)
	}
	return days
}

func (s *AdminService) GetReports(month, department string) ([]model.ReportRow, error) {
	// 1. Parse month (e.g., "2026-06")
	layout := "2006-01"
	startMonth, err := time.Parse(layout, month)
	if err != nil {
		startMonth = time.Now()
	}

	// Calculate start and end dates
	startDate := time.Date(startMonth.Year(), startMonth.Month(), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, -1) // Last day of month

	// For workdays calculation, cap the end date to today if we're querying the current month
	workEnd := endDate
	now := time.Now()
	if now.Before(endDate) {
		workEnd = now
	}
	
	totalWorkDays := calculateWorkDays(startDate, workEnd)

	// Format for DB query
	startDateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")

	// 2. Fetch Users
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// 3. Fetch Logs
	logs, err := s.attendanceRepo.FindAllByDateRange(startDateStr, endDateStr)
	if err != nil {
		return nil, err
	}

	// Map logs by user
	logsByUser := make(map[string][]model.AttendanceLog)
	for _, l := range logs {
		logsByUser[l.UserID.String()] = append(logsByUser[l.UserID.String()], l)
	}

	// 4. Aggregate
	var reports []model.ReportRow
	for _, u := range users {
		if u.Status != "active" {
			continue
		}
		if department != "" && department != "All departments" && u.Department != department {
			continue
		}

		userLogs := logsByUser[u.ID.String()]
		var present, late, incomplete int
		var totalDurationMinutes float64

		for _, l := range userLogs {
			if l.Status == "present" {
				present++
			} else if l.Status == "late" {
				late++
			} else if l.Status == "incomplete" {
				incomplete++
			}

			// Calculate duration if check_in and check_out exist
			if l.CheckInTime != nil && l.CheckOutTime != nil {
				duration := l.CheckOutTime.Sub(*l.CheckInTime)
				totalDurationMinutes += duration.Minutes()
			}
		}

		absent := totalWorkDays - (present + late + incomplete)
		if absent < 0 {
			absent = 0
		}

		completedDays := present + late
		avgHoursStr := "—"
		if completedDays > 0 {
			avgMins := totalDurationMinutes / float64(completedDays)
			h := int(avgMins) / 60
			m := int(avgMins) % 60
			avgHoursStr = fmt.Sprintf("%dh %dm", h, m)
		}

		reports = append(reports, model.ReportRow{
			ID:         u.ID.String(),
			UserID:     u.ID.String(),
			Present:    present,
			Late:       late,
			Incomplete: incomplete,
			Absent:     absent,
			AvgHours:   avgHoursStr,
		})
	}

	return reports, nil
}
