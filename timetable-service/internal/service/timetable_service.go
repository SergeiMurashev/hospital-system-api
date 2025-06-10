package service

import (
	"errors"
	"time"

	"github.com/yourusername/hospital-system-api/timetable-service/internal/domain"
	"github.com/yourusername/hospital-system-api/timetable-service/internal/repository"
	"github.com/yourusername/hospital-system-api/timetable-service/pkg/auth"
)

var (
	ErrTimetableNotFound = errors.New("timetable not found")
	ErrInvalidTimeRange  = errors.New("invalid time range")
	ErrTimeSlotTaken     = errors.New("time slot is already taken")
)

type TimetableService interface {
	CreateTimetable(timetable *domain.Timetable) error
	GetTimetable(id uint) (*domain.Timetable, error)
	UpdateTimetable(timetable *domain.Timetable) error
	DeleteTimetable(id uint) error
	ListTimetables(offset, limit int) ([]*domain.Timetable, error)
	GetAppointments(timetableID uint) ([]*domain.Appointment, error)
	CreateAppointment(timetableID uint, userID uint, time time.Time) error
	DeleteAppointment(id uint) error
}

type timetableService struct {
	repo repository.TimetableRepository
	auth auth.Client
}

func NewTimetableService(repo repository.TimetableRepository, auth auth.Client) TimetableService {
	return &timetableService{
		repo: repo,
		auth: auth,
	}
}

func (s *timetableService) CreateTimetable(timetable *domain.Timetable) error {
	if timetable.From.After(timetable.To) {
		return ErrInvalidTimeRange
	}
	return s.repo.Create(timetable)
}

func (s *timetableService) GetTimetable(id uint) (*domain.Timetable, error) {
	timetable, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrTimetableNotFound
	}
	return timetable, nil
}

func (s *timetableService) UpdateTimetable(timetable *domain.Timetable) error {
	if timetable.From.After(timetable.To) {
		return ErrInvalidTimeRange
	}
	return s.repo.Update(timetable)
}

func (s *timetableService) DeleteTimetable(id uint) error {
	return s.repo.Delete(id)
}

func (s *timetableService) ListTimetables(offset, limit int) ([]*domain.Timetable, error) {
	return s.repo.List(offset, limit)
}

func (s *timetableService) GetAppointments(timetableID uint) ([]*domain.Appointment, error) {
	return s.repo.GetAppointments(timetableID)
}

func (s *timetableService) CreateAppointment(timetableID uint, userID uint, time time.Time) error {
	timetable, err := s.repo.GetByID(timetableID)
	if err != nil {
		return ErrTimetableNotFound
	}

	if time.Before(timetable.From) || time.After(timetable.To) {
		return ErrInvalidTimeRange
	}

	appointments, err := s.repo.GetAppointments(timetableID)
	if err != nil {
		return err
	}

	for _, appointment := range appointments {
		if appointment.AppointmentTime.Equal(time) {
			return ErrTimeSlotTaken
		}
	}

	appointment := &domain.Appointment{
		TimetableID:     timetableID,
		UserID:          userID,
		AppointmentTime: time,
	}

	return s.repo.CreateAppointment(appointment)
}

func (s *timetableService) DeleteAppointment(id uint) error {
	return s.repo.DeleteAppointment(id)
}
