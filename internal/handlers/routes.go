package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	contactHandler := NewContactHandler(db)

	// Contact routes
	contacts := app.Group("/contacts")
	contacts.Get("/", contactHandler.GetContacts)
	contacts.Get("/search", contactHandler.SearchContacts)
	contacts.Get("/:id", contactHandler.GetContact)
	contacts.Post("/", contactHandler.CreateContact)
	contacts.Put("/:id", contactHandler.UpdateContact)
	contacts.Delete("/:id", contactHandler.DeleteContact)
}
