package handlers

import (
	"strconv"

	"api-contacts-go/internal/models"
	"api-contacts-go/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ContactHandler struct {
	service   *services.ContactService
	validator *validator.Validate
}

func NewContactHandler(db *gorm.DB) *ContactHandler {
	return &ContactHandler{
		service:   services.NewContactService(db),
		validator: validator.New(),
	}
}

// GetContacts godoc
// @Summary Get all contacts
// @Description Get paginated list of contacts
// @Tags contacts
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse
// @Router /contacts [get]
func (h *ContactHandler) GetContacts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	contacts, total, err := h.service.GetContacts(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch contacts",
		})
	}

	// Convert to response format
	var contactResponses []models.ContactResponse
	for _, contact := range contacts {
		contactResponses = append(contactResponses, contact.ToResponse())
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return c.JSON(models.PaginatedResponse{
		Data:       contactResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

// GetContact godoc
// @Summary Get contact by ID
// @Description Get a specific contact by ID
// @Tags contacts
// @Accept json
// @Produce json
// @Param id path int true "Contact ID"
// @Success 200 {object} models.ContactResponse
// @Failure 404 {object} map[string]string
// @Router /contacts/{id} [get]
func (h *ContactHandler) GetContact(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid contact ID",
		})
	}

	contact, err := h.service.GetContact(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Contact not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch contact",
		})
	}

	return c.JSON(contact.ToResponse())
}

// CreateContact godoc
// @Summary Create a new contact
// @Description Create a new contact
// @Tags contacts
// @Accept json
// @Produce json
// @Param contact body models.CreateContactRequest true "Contact data"
// @Success 201 {object} models.ContactResponse
// @Failure 400 {object} map[string]string
// @Router /contacts [post]
func (h *ContactHandler) CreateContact(c *fiber.Ctx) error {
	var req models.CreateContactRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed",
			"details": err.Error(),
		})
	}

	contact := &models.Contact{
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Company: req.Company,
	}

	if err := h.service.CreateContact(contact); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create contact",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(contact.ToResponse())
}

// UpdateContact godoc
// @Summary Update a contact
// @Description Update an existing contact
// @Tags contacts
// @Accept json
// @Produce json
// @Param id path int true "Contact ID"
// @Param contact body models.UpdateContactRequest true "Contact data"
// @Success 200 {object} models.ContactResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /contacts/{id} [put]
func (h *ContactHandler) UpdateContact(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid contact ID",
		})
	}

	var req models.UpdateContactRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation failed",
			"details": err.Error(),
		})
	}

	contact, err := h.service.UpdateContact(uint(id), &req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Contact not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update contact",
		})
	}

	return c.JSON(contact.ToResponse())
}

// DeleteContact godoc
// @Summary Delete a contact
// @Description Delete a contact by ID
// @Tags contacts
// @Accept json
// @Produce json
// @Param id path int true "Contact ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /contacts/{id} [delete]
func (h *ContactHandler) DeleteContact(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid contact ID",
		})
	}

	if err := h.service.DeleteContact(uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Contact not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete contact",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// SearchContacts godoc
// @Summary Search contacts
// @Description Search contacts by name or email
// @Tags contacts
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedResponse
// @Router /contacts/search [get]
func (h *ContactHandler) SearchContacts(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search query is required",
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	contacts, total, err := h.service.SearchContacts(query, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search contacts",
		})
	}

	// Convert to response format
	var contactResponses []models.ContactResponse
	for _, contact := range contacts {
		contactResponses = append(contactResponses, contact.ToResponse())
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return c.JSON(models.PaginatedResponse{
		Data:       contactResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}
