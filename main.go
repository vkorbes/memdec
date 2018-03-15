package main

import (
	"log"
	"net/http"

	"github.com/ellenkorbes/memdec/ctrl"
	"github.com/ellenkorbes/memdec/db"
)

func main() {
	d := db.Init("mongodb://user:password@yourdatabase.com:12345/dbname")
	defer d.Close()
	ctrl := ctrl.NewController(d)
	mux := http.NewServeMux()
	mux.HandleFunc("/listall", ctrl.ListAllDecks)
	mux.HandleFunc("/create", ctrl.Create)
	mux.HandleFunc("/info/", ctrl.Info)
	mux.HandleFunc("/nextcard/", ctrl.NextCard)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
