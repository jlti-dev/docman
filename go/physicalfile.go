package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type PhysicalFile struct {
	gorm.Model
	FileName   string    `gorm:"unique" json:"filename"`
	UploadTime time.Time `json:"uploadtime"`
	UploadedBy string    `json:"uploadmail"`
	Type       string    `json:"type"`
}

func uploadPhysicalFile(w http.ResponseWriter, r *http.Request) {
	userContext, err := getContext(r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if userContext.Role != ROLE_ADMIN {
		log.Printf("Benutzer %s hat versucht eine Datei anzulegen", userContext.Email)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	//Einlesen der Datei
	file, handler, err := r.FormFile("uploadFile")
	if err != nil {
		log.Printf("Error Retrieving the File: %s", err.Error())
		return
	}
	defer file.Close()

	//prüfen, ob Datei wohl schon existiert:
	if _, err := os.Stat("./data/" + handler.Filename); errors.Is(err, os.ErrNotExist) {
		log.Printf("Datei %s existiert noch nicht im Dateisystem", handler.Filename)
		
	  }else{
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	  }
	//neue Datei im os anlegen
	resFile, err := os.Create("./data/" + handler.Filename)
	if err != nil {
		fmt.Fprintln(w, err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return;
	}
	defer resFile.Close()

	if err == nil {
		io.Copy(resFile, file)
		defer resFile.Close()
		log.Printf("Successfully Uploaded Original File to os")
	} else {
		log.Printf("Fehler bei os: %s", err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	log.Printf("File %s is written to os, writing to DB", handler.Filename)

	//Datei jetzt in Datenbank aufnehmen
	var dbFile = PhysicalFile{}
	dbFile.FileName = handler.Filename
	dbFile.UploadTime = time.Now()
	dbFile.UploadedBy = userContext.Email
	dbFile.Type = "FILE"

	db.Create(&dbFile)
	if db.Error != nil {
		log.Println(db.Error.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	} else if( dbFile.ID == 0) {
		log.Printf("File %s already existed in db, overwritten OS")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}else{
		log.Printf("Created File with ability to link: %s", dbFile.FileName)
		//json.NewEncoder(w).Encode(dbFile)
		//return
	}
	
	//link erzeugen
	var logicalFile = LogicalFile{}
	logicalFile.LogicalName = dbFile.FileName
	logicalFile.AssignTime = time.Now()
	logicalFile.AssignedBy = userContext.Email
	logicalFile.AssignedTo = userContext.Email
	logicalFile.PhysicalName = dbFile.FileName
	logicalFile.Type = "LINK"
	db.Create(&logicalFile)
	log.Printf("User %s linked physical file %s as %s to user %s", userContext.Email, logicalFile.PhysicalName, logicalFile.LogicalName, logicalFile.AssignedTo)

	json.NewEncoder(w).Encode(logicalFile)
}
func getPhysicalFileByName(filename string) (*PhysicalFile, error) {
	log.Printf("PhysicalFile requested: %s", filename)
	var file *PhysicalFile
	db.Where("file_name = 	?", filename).First(&file)
	if file.FileName == "" {
		log.Printf("File for %s not found", filename)
		return nil, fmt.Errorf("user for %s not found", filename)
	}
	return file, nil
}

//func getAllPhysicalFiles() ([]PhysicalFile, error) {
//	log.Printf("PhysicalFile requested: all")
//	var files []PhysicalFile
//	db.Find(&files)
//	if len(files) == 0 {
//		log.Println(db.Error.Error())
//		return nil, fmt.Errorf("no PhysicalFiles found")
//	}
//	return files, nil
//}

func deletePhysicalFile(file *PhysicalFile) error {
	log.Printf("Delete for file requested: %s", file.FileName)
	err := os.Remove(file.FileName)

	db.Delete(file)
	if err != nil {
		log.Printf("os answers: %s", err.Error())
		return err
	}

	log.Println("Deleting cascading links")
	result := db.Where("PhysicalName = ?", file.FileName).Delete(&LogicalFile{})
	log.Printf("%d rows were effected", result.RowsAffected)
	if result.Error != nil {
		log.Println(result.Error.Error())
		return result.Error
	}
	return nil
}
