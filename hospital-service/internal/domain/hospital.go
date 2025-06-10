package domain

import (
	"time"

	"gorm.io/gorm"
)

type Hospital struct {
	ID        uint64         `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `json:"name"`
	Address   string         `json:"address"`
	Phone     string         `json:"phone"`
	Rooms     []*Room        `gorm:"foreignKey:HospitalID" json:"rooms"`
}

type Room struct {
	ID         uint64         `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	Name       string         `json:"name"`
	HospitalID uint64         `json:"hospital_id"`
}

type CreateHospitalRequest struct {
	Name    string   `json:"name" binding:"required"`
	Address string   `json:"address" binding:"required"`
	Phone   string   `json:"phone" binding:"required"`
	Rooms   []string `json:"rooms" binding:"required"`
}

type UpdateHospitalRequest struct {
	ID      uint64   `json:"id" binding:"required"`
	Name    string   `json:"name" binding:"required"`
	Address string   `json:"address" binding:"required"`
	Phone   string   `json:"phone" binding:"required"`
	Rooms   []string `json:"rooms" binding:"required"`
}
