package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type LogicalFile struct {
	gorm.Model
	LogicalName  string    `json:"filename"`
	AssignTime   time.Time `json:"assigntime"`
	AssignedBy   string    `json:"assigneemail"`
	AssignedTo   string
	PhysicalName string    `json:"physicalname"`
	LastDownload time.Time
	Type         string `json:"type"`
}
type FileLink struct {
	LogicalName  string `json:"logicalname"`
	PhysicalName string `json:"physicalname"`
	AssignedTo   string `json:"assignedto"`
}

func getAllFiles(w http.ResponseWriter, r *http.Request) {
	userContext, err := getContext(r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if(userContext.Role == ROLE_ADMIN){
		var pf []PhysicalFile
		db.Find(&pf)
		json.NewEncoder(w).Encode(pf)
	}else{
		var lf []LogicalFile
		db.Where("assigned_to = ?", userContext.Email).Find(&lf)
	
		json.NewEncoder(w).Encode(lf)
	}
}
func getAllLinks(w http.ResponseWriter, r *http.Request) {
	userContext, err := getContext(r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	mail := mux.Vars(r)["mail"]
	if mail == "" {
		log.Println("Kein mail zur linksuche angegeben")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	
	if(userContext.Role == ROLE_ADMIN){ //Das hier dürfen nur admins aufrufen
		log.Printf("Suche nach Links für Benutzer %s", mail)
		var lf []LogicalFile
		db.Where("assigned_to = ?", mail).Find(&lf)
	
		json.NewEncoder(w).Encode(lf)
	}else{
		log.Println("Endpunkt ist nur für Admins!")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
}
func unlinkFile(w http.ResponseWriter, r *http.Request) {
	userContext, err := getContext(r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if userContext.Role != ROLE_ADMIN {
		log.Printf("Benutzer %s wollte Dateien verlinken", userContext.Email)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	link := FileLink{}

	json.NewDecoder(r.Body).Decode(&link)

	log.Printf("User %s unlinks physical file %s as %s from user %s", userContext.Email, link.PhysicalName, link.LogicalName, link.AssignedTo)
	if link.LogicalName == "" || link.PhysicalName == "" || link.AssignedTo == "" {
		log.Println("Fehler in Link Struktur")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	//Existiert der Benutzer?
	_, err = getUserByMail(link.AssignedTo)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	//Existiert die Physische Datei?
	pf, err := getPhysicalFileByName(link.PhysicalName)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if link.AssignedTo == userContext.Email {
		//Admin versucht eine physische Datei zu löschen
		log.Printf("Admin %s deletes file %s and all links!", userContext.Email, link.PhysicalName)
		deletePhysicalFile(pf)
	} else {
		//Admin entfernt eine Datei von einem Benutzer (also den Link)
		lf, err := getLogicalFileByMailAndName(link.AssignedTo, link.LogicalName)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		db.Delete(lf)
	}
}

func linkFile(w http.ResponseWriter, r *http.Request) {
	userContext, err := getContext(r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if userContext.Role != ROLE_ADMIN {
		log.Printf("Benutzer %s wollte Dateien verlinken", userContext.Email)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	link := FileLink{}

	json.NewDecoder(r.Body).Decode(&link)

	log.Printf("User %s links physical file %s as %s to user %s", userContext.Email, link.PhysicalName, link.LogicalName, link.AssignedTo)
	if link.LogicalName == "" || link.PhysicalName == "" || link.AssignedTo == "" {
		log.Println("Fehler in Link Struktur")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	//Existiert der Benutzer?
	_, err = getUserByMail(link.AssignedTo)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	//Existiert die Physische Datei?
	_, err = getPhysicalFileByName(link.PhysicalName)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	//link erzeugen
	var logicalFile = LogicalFile{}
	logicalFile.LogicalName = link.LogicalName
	logicalFile.AssignTime = time.Now()
	logicalFile.AssignedBy = userContext.Email
	logicalFile.AssignedTo = link.AssignedTo
	logicalFile.PhysicalName = link.PhysicalName
	logicalFile.Type = "LINK"

	db.Create(&logicalFile)
	log.Printf("User %s linked physical file %s as %s to user %s", userContext.Email, link.PhysicalName, link.LogicalName, link.AssignedTo)
}
func getLogicalFileByMailAndName(email, file string) (*LogicalFile, error) {
	var lf *LogicalFile
	db.Where("assigned_to = ?", email).Where("logical_name = ?", file).First(&lf)
	if lf.AssignedTo != email || lf.LogicalName != file {
		log.Printf("Benutzer %s hat keine Logische Datei %s", email, file)
		return nil, fmt.Errorf("%s hat keine Logische Datei %s", email, file)
	} else {
		return lf, nil
	}
}
func downloadFile(w http.ResponseWriter, r *http.Request) {
	userContext, err := getContext(r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	filename := mux.Vars(r)["filename"]
	if filename == "" {
		log.Println("Kein filename zum Download angegeben")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	dbFile, err := getLogicalFileByMailAndName(userContext.Email, filename)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	physFile, err := getPhysicalFileByName(dbFile.PhysicalName)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	osFile, err := os.Open("./data/" + physFile.FileName)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer osFile.Close()

	w.Header().Set("Content-Type", "text/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=\"%s\"", filename))

	log.Printf("Sending physical File: ./data/%s", physFile.FileName)
	io.Copy(w, osFile)
}
