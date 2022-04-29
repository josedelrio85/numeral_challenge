package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	numeral "github.com/josedelrio85/numeral_challenge/pkg"
)

// var envvars numeral.EnvVars
var handler numeral.Handler

func init() {
	REQUEST_SCHEMA, err := numeral.GetSetting("REQUEST_SCHEMA")
	if err != nil {
		log.Fatal(err)
	}
	BANK_FOLDER, err := numeral.GetSetting("BANK_FOLDER")
	if err != nil {
		log.Fatal(err)
	}
	SQLITE_DB_FILE_LOCATION, err := numeral.GetSetting("SQLITE_DB_FILE_LOCATION")
	if err != nil {
		log.Fatal(err)
	}

	handler = numeral.Handler{
		EnvVars: numeral.EnvVars{
			RequestSchema:    REQUEST_SCHEMA,
			BankFolder:       BANK_FOLDER,
			SqliteDbLocation: SQLITE_DB_FILE_LOCATION,
		},
	}
}

func main() {
	log.Println("Numeral Go Coding Challenge starting...")

	db := numeral.CreateDbInstance()
	err := db.Open(handler.EnvVars.SqliteDbLocation)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer db.Db.Close()

	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("error creating the table. err: %s", err)
	}

	handler.Database = db

	router := mux.NewRouter()
	router.Use(numeral.AuthMiddleware)

	router.PathPrefix("/payment/receive").Handler(handler.AddEntity()).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":4567", router))
}
