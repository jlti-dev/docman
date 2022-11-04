package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func getSingleAccount(w http.ResponseWriter, r *http.Request) {
	userContext, err := getContext(r)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	user := userContext
	email := mux.Vars(r)["email"]
	if email != "" {
		if userContext.Role != ROLE_ADMIN {
			log.Printf("Benutzer %s hat versucht %s zu lesen", userContext.Email, email)
		} else {
			//wir sind ADMIN, dürfen also alles lesen
			log.Printf("Benutzer %s hat Benutzer %s gelesen", userContext.Email, email)
			user, err = getUserByMail(email)
			if err != nil {
				log.Println(err)
				user = userContext
			}
		}
	}
	json.NewEncoder(w).Encode(user)

}
func changeAccount(w http.ResponseWriter, r *http.Request) {
	userContext, err := getContext(r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	userIn := User{}

	json.NewDecoder(r.Body).Decode(&userIn)
	validUser, err := validateUser(&userIn)
	if err != nil {
		log.Printf("Validierungsfehler: %s", err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if validUser.Email != userContext.Email {
		if userContext.Role != ROLE_ADMIN {
			//Kein Admin, aber fremder Benutzer!
			log.Printf("Benutzer %s hat versucht %s zu ändern", userContext.Email, validUser.Email)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
	}
	log.Printf("Benutzer %s darf Benutzer %s ändern", userContext.Email, userIn.Email)

	userDB, err := changeUser(validUser.Email, validUser.Name, validUser.Password)
	if err != nil {
		log.Printf("Fehler beim Ändern des Benutzers: %s", err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(userDB)
}
func getAllAccounts(w http.ResponseWriter, r *http.Request) {
	userContext, err := getContext(r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if userContext.Role != ROLE_ADMIN {
		log.Printf("Benutzer %s hat versucht alle Accounts zu lesen", userContext.Email)
		getSingleAccount(w, r)
		return
	}
	userList, err := getAllUser()
	if err != nil {
		//Kann eigentlich nicht passieren, da zumindest der Adminuser aus dem userContext
		//existieren muss
		log.Println(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(userList)
}
func createAccount(w http.ResponseWriter, r *http.Request) {
	userContext, err := getContext(r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if userContext.Role != ROLE_ADMIN {
		log.Printf("Benutzer %s hat versucht einen Account anzulegen", userContext.Email)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	userIn := User{}

	json.NewDecoder(r.Body).Decode(&userIn)
	validUser, err := validateUser(&userIn)
	if err != nil {
		log.Printf("Validierungsfehler: %s", err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	log.Printf("Benutzer %s versucht Benutzer %s anzulegen", userContext.Email, validUser.Email)
	_, err = getUserByMail(validUser.Email)
	if err == nil {
		//Benutzer existiert
		log.Printf("Benutzer mit Email %s existiert schon", validUser.Email)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	dbUser, err := createUser(validUser.Email, validUser.Name, validUser.Password, validUser.Role)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(dbUser)

}
