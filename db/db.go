package db

import (
	"bookstore/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB // Define global variable

// InitDatabase initializes and returns a new database instance
func InitDatabase() (*gorm.DB, error) {
	var err error
	DB, err = gorm.Open(sqlite.Open("bookstore.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return nil, err
	}

	// Auto-migrate the Book table
	err = DB.AutoMigrate(&models.Book{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
		return nil, err
	}

	log.Println("Database connected and migrated successfully.")
	return DB, nil
}
