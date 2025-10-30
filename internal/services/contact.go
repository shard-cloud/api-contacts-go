package services

import (
	"fmt"
	"strings"

	"api-contacts-go/internal/models"

	"gorm.io/gorm"
)

type ContactService struct {
	db *gorm.DB
}

func NewContactService(db *gorm.DB) *ContactService {
	return &ContactService{db: db}
}

func (s *ContactService) GetContacts(page, limit int) ([]models.Contact, int64, error) {
	var contacts []models.Contact
	var total int64

	// Count total records
	if err := s.db.Model(&models.Contact{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := s.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&contacts).Error; err != nil {
		return nil, 0, err
	}

	return contacts, total, nil
}

func (s *ContactService) GetContact(id uint) (*models.Contact, error) {
	var contact models.Contact
	if err := s.db.First(&contact, id).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

func (s *ContactService) CreateContact(contact *models.Contact) error {
	return s.db.Create(contact).Error
}

func (s *ContactService) UpdateContact(id uint, req *models.UpdateContactRequest) (*models.Contact, error) {
	var contact models.Contact
	if err := s.db.First(&contact, id).Error; err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		contact.Name = *req.Name
	}
	if req.Email != nil {
		contact.Email = *req.Email
	}
	if req.Phone != nil {
		contact.Phone = *req.Phone
	}
	if req.Company != nil {
		contact.Company = *req.Company
	}

	if err := s.db.Save(&contact).Error; err != nil {
		return nil, err
	}

	return &contact, nil
}

func (s *ContactService) DeleteContact(id uint) error {
	return s.db.Delete(&models.Contact{}, id).Error
}

func (s *ContactService) SearchContacts(query string, page, limit int) ([]models.Contact, int64, error) {
	var contacts []models.Contact
	var total int64

	// Build search query
	searchQuery := s.db.Model(&models.Contact{}).Where(
		"LOWER(name) LIKE ? OR LOWER(email) LIKE ? OR LOWER(company) LIKE ?",
		"%"+strings.ToLower(query)+"%",
		"%"+strings.ToLower(query)+"%",
		"%"+strings.ToLower(query)+"%",
	)

	// Count total matching records
	if err := searchQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := searchQuery.Offset(offset).Limit(limit).Order("created_at DESC").Find(&contacts).Error; err != nil {
		return nil, 0, err
	}

	return contacts, total, nil
}
