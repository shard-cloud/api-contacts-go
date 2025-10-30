package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"api-contacts-go/internal/config"
	"api-contacts-go/internal/database"
	"api-contacts-go/internal/handlers"
	"api-contacts-go/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database")
	}

	// Auto migrate
	db.AutoMigrate(&models.Contact{})

	return db
}

func setupTestApp(db *gorm.DB) *fiber.App {
	app := fiber.New()
	api := app.Group("/api/v1")
	handlers.SetupRoutes(api, db)
	return app
}

func TestHealthCheck(t *testing.T) {
	app := fiber.New()
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestCreateContact(t *testing.T) {
	db := setupTestDB()
	app := setupTestApp(db)

	contact := models.CreateContactRequest{
		Name:    "Test User",
		Email:   "test@example.com",
		Phone:   "+55 11 99999-9999",
		Company: "Test Company",
	}

	jsonData, _ := json.Marshal(contact)
	req := httptest.NewRequest("POST", "/api/v1/contacts", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	var response models.ContactResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, contact.Name, response.Name)
	assert.Equal(t, contact.Email, response.Email)
}

func TestGetContacts(t *testing.T) {
	db := setupTestDB()
	app := setupTestApp(db)

	// Create test contact
	contact := &models.Contact{
		Name:    "Test User",
		Email:   "test@example.com",
		Phone:   "+55 11 99999-9999",
		Company: "Test Company",
	}
	db.Create(contact)

	req := httptest.NewRequest("GET", "/api/v1/contacts", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var response models.PaginatedResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), response.Total)
	assert.Len(t, response.Data, 1)
}

func TestGetContact(t *testing.T) {
	db := setupTestDB()
	app := setupTestApp(db)

	// Create test contact
	contact := &models.Contact{
		Name:    "Test User",
		Email:   "test@example.com",
		Phone:   "+55 11 99999-9999",
		Company: "Test Company",
	}
	db.Create(contact)

	req := httptest.NewRequest("GET", "/api/v1/contacts/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var response models.ContactResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, contact.Name, response.Name)
}

func TestUpdateContact(t *testing.T) {
	db := setupTestDB()
	app := setupTestApp(db)

	// Create test contact
	contact := &models.Contact{
		Name:    "Test User",
		Email:   "test@example.com",
		Phone:   "+55 11 99999-9999",
		Company: "Test Company",
	}
	db.Create(contact)

	updateData := models.UpdateContactRequest{
		Name: stringPtr("Updated Name"),
	}

	jsonData, _ := json.Marshal(updateData)
	req := httptest.NewRequest("PUT", "/api/v1/contacts/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var response models.ContactResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", response.Name)
}

func TestDeleteContact(t *testing.T) {
	db := setupTestDB()
	app := setupTestApp(db)

	// Create test contact
	contact := &models.Contact{
		Name:    "Test User",
		Email:   "test@example.com",
		Phone:   "+55 11 99999-9999",
		Company: "Test Company",
	}
	db.Create(contact)

	req := httptest.NewRequest("DELETE", "/api/v1/contacts/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)

	// Verify contact is deleted
	var count int64
	db.Model(&models.Contact{}).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestSearchContacts(t *testing.T) {
	db := setupTestDB()
	app := setupTestApp(db)

	// Create test contacts
	contacts := []models.Contact{
		{Name: "João Silva", Email: "joao@example.com"},
		{Name: "Maria Santos", Email: "maria@example.com"},
		{Name: "Pedro Oliveira", Email: "pedro@example.com"},
	}

	for _, contact := range contacts {
		db.Create(&contact)
	}

	req := httptest.NewRequest("GET", "/api/v1/contacts/search?q=João", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var response models.PaginatedResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), response.Total)
	assert.Len(t, response.Data, 1)
	assert.Equal(t, "João Silva", response.Data[0].Name)
}

func stringPtr(s string) *string {
	return &s
}
