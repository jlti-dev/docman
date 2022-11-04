package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getDBConnection() (db *gorm.DB) {
	conn, err := gorm.Open(sqlite.Open("./data.sq3"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %s", err)
	}
	return conn
}
func initializeDB() (db *gorm.DB) {
	log.Println("Connecting to database")
	database := getDBConnection()
	log.Println("Connected to database")
	log.Println("Running migration")
	database.AutoMigrate(User{})
	database.AutoMigrate(PhysicalFile{})
	database.AutoMigrate(LogicalFile{})
	database.AutoMigrate(FileLink{})
	log.Println("Finished migration")
	return database
}
