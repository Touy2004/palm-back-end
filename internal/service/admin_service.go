package service

import (
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
