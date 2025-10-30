package main

import (
	"fmt"
	"log"
	"os"

	"api-contacts-go/internal/config"
	"api-contacts-go/internal/database"
	"api-contacts-go/internal/models"

	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Seed data
	if err := seedContacts(db); err != nil {
		log.Fatal("Failed to seed contacts:", err)
	}

	fmt.Println("✅ Seed completed successfully!")
}

func seedContacts(db *gorm.DB) error {
	// Check if contacts already exist
	var count int64
	db.Model(&models.Contact{}).Count(&count)
	if count > 0 {
		fmt.Println("Contacts already exist, skipping seed")
		return nil
	}

	contacts := []models.Contact{
		{
			Name:    "João Silva",
			Email:   "joao.silva@example.com",
			Phone:   "+55 11 99999-1111",
			Company: "Tech Corp",
		},
		{
			Name:    "Maria Santos",
			Email:   "maria.santos@example.com",
			Phone:   "+55 11 99999-2222",
			Company: "Design Studio",
		},
		{
			Name:    "Pedro Oliveira",
			Email:   "pedro.oliveira@example.com",
			Phone:   "+55 11 99999-3333",
			Company: "Marketing Agency",
		},
		{
			Name:    "Ana Costa",
			Email:   "ana.costa@example.com",
			Phone:   "+55 11 99999-4444",
			Company: "Consulting Group",
		},
		{
			Name:    "Carlos Ferreira",
			Email:   "carlos.ferreira@example.com",
			Phone:   "+55 11 99999-5555",
			Company: "Startup Inc",
		},
		{
			Name:    "Lucia Rodrigues",
			Email:   "lucia.rodrigues@example.com",
			Phone:   "+55 11 99999-6666",
			Company: "Finance Corp",
		},
		{
			Name:    "Roberto Alves",
			Email:   "roberto.alves@example.com",
			Phone:   "+55 11 99999-7777",
			Company: "Healthcare Ltd",
		},
		{
			Name:    "Fernanda Lima",
			Email:   "fernanda.lima@example.com",
			Phone:   "+55 11 99999-8888",
			Company: "Education Center",
		},
		{
			Name:    "Marcos Pereira",
			Email:   "marcos.pereira@example.com",
			Phone:   "+55 11 99999-9999",
			Company: "Real Estate",
		},
		{
			Name:    "Juliana Rocha",
			Email:   "juliana.rocha@example.com",
			Phone:   "+55 11 99999-0000",
			Company: "Retail Store",
		},
	}

	// Insert contacts
	for _, contact := range contacts {
		if err := db.Create(&contact).Error; err != nil {
			return fmt.Errorf("failed to create contact %s: %w", contact.Name, err)
		}
	}

	fmt.Printf("✅ Created %d contacts\n", len(contacts))
	return nil
}
