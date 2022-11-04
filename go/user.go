package main

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"log"
)

type User struct {
	gorm.Model
	Email     string    `gorm:"unique" json:"email"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	LastLogin time.Time `json:"login"`
}

func checkAndCreateAdminUser(db *gorm.DB, email string, password string) {
	log.Println("Prüfe auf Administrativen Zugriff")
	_, err := getUserByMail(email)
	if err != nil {
		log.Printf("Admin with email %s does not exist", email)
		createUser(email, "Administrator", password, ROLE_ADMIN)
	} else {
		log.Printf("Admin with email %s already exists overwriting password", email)
		_, err = changeUser(email, "", password)
		if err != nil {
			log.Println(err)
		}
	}
}
func createUser(email, name, password, role string) (*User, error) {
	_, err := getUserByMail(email)
	if err == nil {
		return nil, fmt.Errorf("user with mail %s already exists", email)
	}
	var dbUser User
	dbUser.Email = email
	dbUser.Name = name
	dbUser.Password, _ = GenerateHashPassword(password)
	dbUser.Role = role

	db.Create(&dbUser)
	log.Printf("Created User %s with role %s", dbUser.Email, dbUser.Role)
	return &dbUser, db.Error
}
func changeUser(email, name, password string) (*User, error) {
	user, err := getUserByMail(email)
	if err != nil {
		return nil, err
	}
	change := false
	pwHash, _ := GenerateHashPassword(password)
	if user.Name == name && user.Password == pwHash {
		log.Printf("User with email %s would not be changed", email)
		return user, nil
	}
	if user.Password != pwHash && password != "" {
		log.Printf("Updating password for %s", email)
		user.Password = pwHash
		change = true
	}
	if user.Name != name && name != "" {
		log.Printf("Updating name for %s", email)
		user.Name = name
		change = true
	}
	if !change {
		log.Printf("User with email %s would not be changed", email)
		return user, nil
	}
	db.Save(&user)
	return user, db.Error
}
func setUserLogin(email string) (*User, error) {
	dbUser, err := getUserByMail(email)
	if err != nil {
		return nil, err
	}
	log.Printf("Setze letzte Login für Benutzer")
	dbUser.LastLogin = time.Now()
	db.Save(&dbUser)
	return dbUser, db.Error
}
func getUserByMail(email string) (*User, error) {
	log.Printf("User requested: %s", email)
	var dbUser *User
	db.Where("email = 	?", email).First(&dbUser)
	if dbUser.Email == "" {
		log.Printf("user for %s not found", email)
		return nil, fmt.Errorf("user for %s not found", email)
	}
	return dbUser, nil
}
func getAllUser() ([]User, error) {
	log.Printf("User requested: all")
	var dbUser []User
	db.Find(&dbUser)
	if len(dbUser) == 0 {
		log.Println(db.Error.Error())
		return nil, fmt.Errorf("no users found")
	}
	return dbUser, nil

}

func validateUser(userIn *User) (*User, error) {
	switch "" {
	case userIn.Email:
		return nil, fmt.Errorf("no mail")
	//case userIn.Password:
	//	return nil, fmt.Errorf("no password")
	//case userIn.Role:
	//	return nil, fmt.Errorf("no role")
	//case userIn.Name:
	//Ohne Name können wir leben
	//	return nil, fmt.Errorf("No name")
	default:
		//nothing!
	}
	if userIn.Role != ROLE_ADMIN && userIn.Role != ROLE_USER {
		//Unbekannte Rolle
		return nil, fmt.Errorf("unknown role")
	}
	return userIn, nil
}

//take password as input and generate new hash password from it
func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
