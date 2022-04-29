package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	numeral "github.com/josedelrio85/numeral_challenge/pkg"
)

func main() {
	log.Println("Numeral Go Coding Challenge starting...")

	db := numeral.CreateDbInstance()
	err := db.Open()
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer db.Db.Close()

	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("error creating the table. err: %s", err)
	}

	handler := numeral.Handler{}
	handler.Database = db

	router := mux.NewRouter()
	router.Use(numeral.AuthMiddleware)

	router.PathPrefix("/payment/receive").Handler(handler.AddEntity()).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":4567", router))
}
