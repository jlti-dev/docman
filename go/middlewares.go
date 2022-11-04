package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Printf("[API] %s: Begin of %s\n", r.Method, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
		log.Printf("[API] %s: %s\n", r.Method, r.RequestURI)
	})
}
func authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/api/signin" {
			//login ist unauthorisiert
			next.ServeHTTP(w, r)
			return
		}else if ! strings.HasPrefix(r.RequestURI, "/api/"){
			//nur die API ist passwort gesichert, der Rest nicht
			next.ServeHTTP(w,r)
			return
		}
		log.Printf("Authorization Check")
		//Existiert ein token?
		if r.Header["Token"] == nil {
			log.Println("Token fehlt")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		_, err := getTokenData(r.Header["Token"][0])
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		log.Print("Authorized")
		next.ServeHTTP(w, r)
	})
}
func getTokenData(headerString string) (*User, error) {
	//Token prüfen
	token, err := jwt.Parse(headerString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error in parsing token")
		}
		return []byte(secretkey), nil
	})

	//Token konnte nicht geparsed werden
	if err != nil {
		log.Println("Token konnte nicht geparsed werden")
		return nil, fmt.Errorf("Token konnte nicht geparsed werden")
	}

	//Token wurde manipuliert
	if !token.Valid {
		log.Println("Token wurde manipuliert")
		return nil, fmt.Errorf("Token wurde manipuliert")
	}
	//Claims korrekt?
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return getUserByMail(fmt.Sprintf("%s", claims["email"]))
	}
	return nil, fmt.Errorf("claims are wrong")
}
