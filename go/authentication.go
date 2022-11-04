package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Token struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	Token string `json:"token"`
}

//Generate JWT token
func GenerateJWT(email, role string) (string, error) {
	var mySigningKey = []byte(secretkey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Fatalf("Something went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func signIn(w http.ResponseWriter, r *http.Request) {
	var authDetails Authentication

	//JSON lesbar?
	err := json.NewDecoder(r.Body).Decode(&authDetails)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	//Benutzer lesbar?
	dbUser, err := getUserByMail(authDetails.Email)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	//Benutzer existiert also, prüfen, ob Passwort identisch ist:
	if !CheckPasswordHash(authDetails.Password, dbUser.Password) {
		log.Println("Password mismatch")
		//Password hashes passen nicht
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	//Token generieren
	validToken, err := GenerateJWT(dbUser.Email, dbUser.Role)
	if err != nil {
		log.Println("Token invalid")
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	var token Token
	token.Email = dbUser.Email
	token.Role = dbUser.Role
	token.Token = validToken
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)

	//update des Benutzerstamms, Fehler sind hier irrelevant
	setUserLogin(token.Email)
}

//compare plain password with hash password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Println(err.Error())
	}
	return err == nil
}
