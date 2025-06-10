package repository

import (
	"github.com/sergeimurashev/hospital-system-api/document-service/internal/domain"
	"gorm.io/gorm"
)

type DocumentRepository interface {
	Create(doc *domain.Document) error
	GetByID(id uint) (*domain.Document, error)
	Update(doc *domain.Document) error
	Delete(id uint) error
	GetByPatientID(patientID uint) ([]*domain.Document, error)
}

type documentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) DocumentRepository {
	return &documentRepository{db: db}
}

func (r *documentRepository) Create(doc *domain.Document) error {
	return r.db.Create(doc).Error
}

func (r *documentRepository) GetByID(id uint) (*domain.Document, error) {
	var doc domain.Document
	if err := r.db.First(&doc, id).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *documentRepository) Update(doc *domain.Document) error {
	return r.db.Save(doc).Error
}

func (r *documentRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Document{}, id).Error
}

func (r *documentRepository) GetByPatientID(patientID uint) ([]*domain.Document, error) {
	var docs []*domain.Document
	if err := r.db.Where("patient_id = ?", patientID).Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}
