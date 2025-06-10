package service

import (
	"context"

	"github.com/sergeimurashev/hospital-system-api/hospital-service/internal/domain"
	"github.com/sergeimurashev/hospital-system-api/hospital-service/internal/repository"
)

type HospitalService interface {
	Create(ctx context.Context, req domain.CreateHospitalRequest) (*domain.Hospital, error)
	GetByID(ctx context.Context, id uint64) (*domain.Hospital, error)
	Update(ctx context.Context, req domain.UpdateHospitalRequest) (*domain.Hospital, error)
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, offset, limit int) ([]*domain.Hospital, int64, error)
	GetRooms(ctx context.Context, hospitalID uint64) ([]*domain.Room, error)
}

type hospitalService struct {
	hospitalRepo repository.HospitalRepository
	roomRepo     repository.RoomRepository
}

func NewHospitalService(hospitalRepo repository.HospitalRepository, roomRepo repository.RoomRepository) HospitalService {
	return &hospitalService{
		hospitalRepo: hospitalRepo,
		roomRepo:     roomRepo,
	}
}

func (s *hospitalService) Create(ctx context.Context, req domain.CreateHospitalRequest) (*domain.Hospital, error) {
	hospital := &domain.Hospital{
		Name:    req.Name,
		Address: req.Address,
		Phone:   req.Phone,
	}

	if err := s.hospitalRepo.Create(ctx, hospital); err != nil {
		return nil, err
	}

	rooms := make([]*domain.Room, len(req.Rooms))
	for i, roomName := range req.Rooms {
		room := &domain.Room{
			Name:       roomName,
			HospitalID: hospital.ID,
		}
		if err := s.roomRepo.Create(ctx, room); err != nil {
			return nil, err
		}
		rooms[i] = room
	}
	hospital.Rooms = rooms

	return hospital, nil
}

func (s *hospitalService) GetByID(ctx context.Context, id uint64) (*domain.Hospital, error) {
	hospital, err := s.hospitalRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	rooms, err := s.roomRepo.GetByHospitalID(ctx, id)
	if err != nil {
		return nil, err
	}
	hospital.Rooms = rooms

	return hospital, nil
}

func (s *hospitalService) Update(ctx context.Context, req domain.UpdateHospitalRequest) (*domain.Hospital, error) {
	hospital, err := s.hospitalRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	hospital.Name = req.Name
	hospital.Address = req.Address
	hospital.Phone = req.Phone

	if err := s.hospitalRepo.Update(ctx, hospital); err != nil {
		return nil, err
	}

	// Delete existing rooms
	if err := s.roomRepo.DeleteByHospitalID(ctx, hospital.ID); err != nil {
		return nil, err
	}

	// Create new rooms
	rooms := make([]*domain.Room, len(req.Rooms))
	for i, roomName := range req.Rooms {
		room := &domain.Room{
			Name:       roomName,
			HospitalID: hospital.ID,
		}
		if err := s.roomRepo.Create(ctx, room); err != nil {
			return nil, err
		}
		rooms[i] = room
	}
	hospital.Rooms = rooms

	return hospital, nil
}

func (s *hospitalService) Delete(ctx context.Context, id uint64) error {
	if err := s.roomRepo.DeleteByHospitalID(ctx, id); err != nil {
		return err
	}
	return s.hospitalRepo.Delete(ctx, id)
}

func (s *hospitalService) List(ctx context.Context, offset, limit int) ([]*domain.Hospital, int64, error) {
	hospitals, err := s.hospitalRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.hospitalRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	for _, hospital := range hospitals {
		rooms, err := s.roomRepo.GetByHospitalID(ctx, hospital.ID)
		if err != nil {
			return nil, 0, err
		}
		hospital.Rooms = rooms
	}

	return hospitals, total, nil
}

func (s *hospitalService) GetRooms(ctx context.Context, hospitalID uint64) ([]*domain.Room, error) {
	return s.roomRepo.GetByHospitalID(ctx, hospitalID)
}
