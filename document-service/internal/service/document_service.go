package service

import (
	"errors"

	"github.com/yourusername/hospital-system-api/document-service/internal/domain"
	"github.com/yourusername/hospital-system-api/document-service/internal/repository"
	"github.com/yourusername/hospital-system-api/document-service/pkg/auth"
	"github.com/yourusername/hospital-system-api/document-service/pkg/elasticsearch"
)

var (
	ErrDocumentNotFound = errors.New("document not found")
	ErrUnauthorized     = errors.New("unauthorized access")
)

type DocumentService interface {
	CreateDocument(doc *domain.Document) error
	GetDocument(id uint) (*domain.Document, error)
	UpdateDocument(doc *domain.Document) error
	DeleteDocument(id uint) error
	GetPatientDocuments(patientID uint) ([]*domain.Document, error)
	SearchDocuments(query string) ([]*domain.Document, error)
}

type documentService struct {
	repo repository.DocumentRepository
	es   elasticsearch.Client
	auth auth.Client
}

func NewDocumentService(repo repository.DocumentRepository, es elasticsearch.Client, auth auth.Client) DocumentService {
	return &documentService{
		repo: repo,
		es:   es,
		auth: auth,
	}
}

func (s *documentService) CreateDocument(doc *domain.Document) error {
	if err := s.repo.Create(doc); err != nil {
		return err
	}

	// Index document in Elasticsearch
	if err := s.es.IndexDocument(doc); err != nil {
		// Log the error but don't fail the request
		// TODO: Add proper logging
	}

	return nil
}

func (s *documentService) GetDocument(id uint) (*domain.Document, error) {
	doc, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrDocumentNotFound
	}
	return doc, nil
}

func (s *documentService) UpdateDocument(doc *domain.Document) error {
	if err := s.repo.Update(doc); err != nil {
		return err
	}

	// Update document in Elasticsearch
	if err := s.es.IndexDocument(doc); err != nil {
		// Log the error but don't fail the request
		// TODO: Add proper logging
	}

	return nil
}

func (s *documentService) DeleteDocument(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// Delete document from Elasticsearch
	if err := s.es.DeleteDocument(id); err != nil {
		// Log the error but don't fail the request
		// TODO: Add proper logging
	}

	return nil
}

func (s *documentService) GetPatientDocuments(patientID uint) ([]*domain.Document, error) {
	return s.repo.GetByPatientID(patientID)
}

func (s *documentService) SearchDocuments(query string) ([]*domain.Document, error) {
	return s.es.SearchDocuments(query)
}
