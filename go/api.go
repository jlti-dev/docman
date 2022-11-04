package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func runApi(port int) {
	log.Println("Initializing API")
	router := mux.NewRouter()

	log.Println("Adding Logging")
	router.Use(loggingMiddleware)

	log.Println("Adding authentication")
	router.Use(authenticationMiddleware) //Filtert die Route "/signin selbst raus"

	log.Println("Adding Route for login")
	router.HandleFunc("/api/signin", signIn).Methods("POST")

	log.Println("Adding Route for account Management")
	router.HandleFunc("/api/account", getSingleAccount).Methods("GET")
	router.HandleFunc("/api/account/{email}", getSingleAccount).Methods("GET")
	router.HandleFunc("/api/accounts", getAllAccounts).Methods("GET")
	router.HandleFunc("/api/account", changeAccount).Methods("PUT")
	router.HandleFunc("/api/account/{email}", changeAccount).Methods("PUT")
	router.HandleFunc("/api/account", createAccount).Methods("POST")

	log.Println("Adding Route for files")
	router.HandleFunc("/api/file", uploadPhysicalFile).Methods("POST")
	router.HandleFunc("/api/file/{filename}", downloadFile).Methods("GET")
	router.HandleFunc("/api/file", linkFile).Methods("PUT")
	router.HandleFunc("/api/file", unlinkFile).Methods("DELETE")
	router.HandleFunc("/api/file", getAllFiles).Methods("GET")

	log.Println("Adding Route for document")
	router.HandleFunc("/api/link/{mail}", getAllLinks).Methods("GET")

	//Static content
	log.Println("Adding route for static files")
	staticDir := "/docker/public"
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(staticDir))))

	log.Println("Adding OPTIONS Method")
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	})
	log.Println("Starting normal API operation")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
func getContext(r *http.Request) (*User, error) {
	return getTokenData(r.Header["Token"][0])
}
