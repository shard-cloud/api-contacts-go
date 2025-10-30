package models

import (
	"time"

	"gorm.io/gorm"
)

type Contact struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null" validate:"required,min=2,max=100"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Phone     string         `json:"phone" gorm:"size:20" validate:"omitempty,min=10,max=20"`
	Company   string         `json:"company" gorm:"size:100" validate:"omitempty,max=100"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateContactRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=100"`
	Email   string `json:"email" validate:"required,email"`
	Phone   string `json:"phone" validate:"omitempty,min=10,max=20"`
	Company string `json:"company" validate:"omitempty,max=100"`
}

type UpdateContactRequest struct {
	Name    *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Email   *string `json:"email,omitempty" validate:"omitempty,email"`
	Phone   *string `json:"phone,omitempty" validate:"omitempty,min=10,max=20"`
	Company *string `json:"company,omitempty" validate:"omitempty,max=100"`
}

type ContactResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Company   string    `json:"company"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PaginatedResponse struct {
	Data       []ContactResponse `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

func (c *Contact) ToResponse() ContactResponse {
	return ContactResponse{
		ID:        c.ID,
		Name:      c.Name,
		Email:     c.Email,
		Phone:     c.Phone,
		Company:   c.Company,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
