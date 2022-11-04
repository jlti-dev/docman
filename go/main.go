package main

import (
	"log"
	"os"

	"gorm.io/gorm"
)

var db *gorm.DB
var secretkey string
var admin_email string
var admin_pw string

func main() {
	log.Println("Application started")
	db = initializeDB()
	log.Println("Setting defaults")
	setGlobalDefaults()
	checkAndCreateAdminUser(db, admin_email, admin_pw)

	runApi(8080)

}
func setGlobalDefaults() {
	log.Println("Setting global defaults")
	secretkey = os.Getenv("DEFAULT_KEY")
	admin_email = os.Getenv("DEFAULT_MAIL")
	admin_pw = os.Getenv("DEFAULT_PW")
	log.Println("Setting global defaults done")
}
