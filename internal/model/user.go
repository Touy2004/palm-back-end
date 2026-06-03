package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	EmployeeCode string    `gorm:"type:varchar(50);unique" json:"employee_code"`
	FullName     string    `gorm:"type:varchar(150);not null" json:"full_name"`
	Phone        string    `gorm:"type:varchar(30);uniqueIndex;not null" json:"phone"`
	Email        string    `gorm:"type:varchar(150);uniqueIndex;not null" json:"email"`
	Department   string    `gorm:"type:varchar(100);not null" json:"department"`
	PasswordHash string    `gorm:"type:text" json:"-"`
	Role         string    `gorm:"type:varchar(30);default:'employee'" json:"role"`
	Status       string    `gorm:"type:varchar(30);default:'active'" json:"status"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}