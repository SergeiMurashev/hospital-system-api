package domain

import (
	"time"

	"gorm.io/gorm"
)

type Document struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Date       time.Time      `gorm:"not null" json:"date"`
	PatientID  uint           `gorm:"not null" json:"patient_id"`
	HospitalID uint           `gorm:"not null" json:"hospital_id"`
	DoctorID   uint           `gorm:"not null" json:"doctor_id"`
	Room       string         `gorm:"not null" json:"room"`
	Data       string         `gorm:"type:text;not null" json:"data"`
}

type CreateDocumentRequest struct {
	Date       time.Time `json:"date" binding:"required"`
	PatientID  uint      `json:"patient_id" binding:"required"`
	HospitalID uint      `json:"hospital_id" binding:"required"`
	DoctorID   uint      `json:"doctor_id" binding:"required"`
	Room       string    `json:"room" binding:"required"`
	Data       string    `json:"data" binding:"required"`
}

type UpdateDocumentRequest struct {
	Date       time.Time `json:"date" binding:"required"`
	PatientID  uint      `json:"patient_id" binding:"required"`
	HospitalID uint      `json:"hospital_id" binding:"required"`
	DoctorID   uint      `json:"doctor_id" binding:"required"`
	Room       string    `json:"room" binding:"required"`
	Data       string    `json:"data" binding:"required"`
}
