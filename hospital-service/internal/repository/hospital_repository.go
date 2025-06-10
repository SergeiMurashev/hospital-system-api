package repository

import (
	"context"

	"github.com/sergeimurashev/hospital-system-api/hospital-service/internal/domain"
	"gorm.io/gorm"
)

type HospitalRepository interface {
	Create(ctx context.Context, hospital *domain.Hospital) error
	GetByID(ctx context.Context, id uint64) (*domain.Hospital, error)
	Update(ctx context.Context, hospital *domain.Hospital) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, offset, limit int) ([]*domain.Hospital, error)
	Count(ctx context.Context) (int64, error)
}

type hospitalRepository struct {
	db *gorm.DB
}

func NewHospitalRepository(db *gorm.DB) HospitalRepository {
	return &hospitalRepository{
		db: db,
	}
}

func (r *hospitalRepository) Create(ctx context.Context, hospital *domain.Hospital) error {
	return r.db.WithContext(ctx).Create(hospital).Error
}

func (r *hospitalRepository) GetByID(ctx context.Context, id uint64) (*domain.Hospital, error) {
	var hospital domain.Hospital
	if err := r.db.WithContext(ctx).First(&hospital, id).Error; err != nil {
		return nil, err
	}
	return &hospital, nil
}

func (r *hospitalRepository) Update(ctx context.Context, hospital *domain.Hospital) error {
	return r.db.WithContext(ctx).Save(hospital).Error
}

func (r *hospitalRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&domain.Hospital{}, id).Error
}

func (r *hospitalRepository) List(ctx context.Context, offset, limit int) ([]*domain.Hospital, error) {
	var hospitals []*domain.Hospital
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&hospitals).Error; err != nil {
		return nil, err
	}
	return hospitals, nil
}

func (r *hospitalRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.Hospital{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
