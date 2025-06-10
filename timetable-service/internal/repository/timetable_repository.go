package repository

import (
	"github.com/yourusername/hospital-system-api/timetable-service/internal/domain"
	"gorm.io/gorm"
)

type TimetableRepository interface {
	Create(timetable *domain.Timetable) error
	GetByID(id uint) (*domain.Timetable, error)
	Update(timetable *domain.Timetable) error
	Delete(id uint) error
	List(offset, limit int) ([]*domain.Timetable, error)
	GetAppointments(timetableID uint) ([]*domain.Appointment, error)
	CreateAppointment(appointment *domain.Appointment) error
	DeleteAppointment(id uint) error
}

type timetableRepository struct {
	db *gorm.DB
}

func NewTimetableRepository(db *gorm.DB) TimetableRepository {
	return &timetableRepository{db: db}
}

func (r *timetableRepository) Create(timetable *domain.Timetable) error {
	return r.db.Create(timetable).Error
}

func (r *timetableRepository) GetByID(id uint) (*domain.Timetable, error) {
	var timetable domain.Timetable
	if err := r.db.First(&timetable, id).Error; err != nil {
		return nil, err
	}
	return &timetable, nil
}

func (r *timetableRepository) Update(timetable *domain.Timetable) error {
	return r.db.Save(timetable).Error
}

func (r *timetableRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Timetable{}, id).Error
}

func (r *timetableRepository) List(offset, limit int) ([]*domain.Timetable, error) {
	var timetables []*domain.Timetable
	if err := r.db.Offset(offset).Limit(limit).Find(&timetables).Error; err != nil {
		return nil, err
	}
	return timetables, nil
}

func (r *timetableRepository) GetAppointments(timetableID uint) ([]*domain.Appointment, error) {
	var appointments []*domain.Appointment
	if err := r.db.Where("timetable_id = ?", timetableID).Find(&appointments).Error; err != nil {
		return nil, err
	}
	return appointments, nil
}

func (r *timetableRepository) CreateAppointment(appointment *domain.Appointment) error {
	return r.db.Create(appointment).Error
}

func (r *timetableRepository) DeleteAppointment(id uint) error {
	return r.db.Delete(&domain.Appointment{}, id).Error
}
