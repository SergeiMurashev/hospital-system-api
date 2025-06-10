package domain

import (
	"time"

	"gorm.io/gorm"
)

type Timetable struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	HospitalID uint           `gorm:"not null" json:"hospital_id"`
	DoctorID   uint           `gorm:"not null" json:"doctor_id"`
	From       time.Time      `gorm:"not null" json:"from"`
	To         time.Time      `gorm:"not null" json:"to"`
	Room       string         `gorm:"not null" json:"room"`
}

type Appointment struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	TimetableID     uint           `gorm:"not null" json:"timetable_id"`
	UserID          uint           `gorm:"not null" json:"user_id"`
	AppointmentTime time.Time      `gorm:"not null" json:"appointment_time"`
}

type CreateTimetableRequest struct {
	HospitalID uint      `json:"hospital_id" binding:"required"`
	DoctorID   uint      `json:"doctor_id" binding:"required"`
	From       time.Time `json:"from" binding:"required"`
	To         time.Time `json:"to" binding:"required"`
	Room       string    `json:"room" binding:"required"`
}

type UpdateTimetableRequest struct {
	HospitalID uint      `json:"hospital_id" binding:"required"`
	DoctorID   uint      `json:"doctor_id" binding:"required"`
	From       time.Time `json:"from" binding:"required"`
	To         time.Time `json:"to" binding:"required"`
	Room       string    `json:"room" binding:"required"`
}

type CreateAppointmentRequest struct {
	Time time.Time `json:"time" binding:"required"`
}
