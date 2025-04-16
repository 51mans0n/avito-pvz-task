package main

import (
	"fmt"
	"github.com/51mans0n/avito-pvz-task/internal/api"
	"github.com/51mans0n/avito-pvz-task/internal/db"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting HTTP service on :8080...")

	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("failed to init DB: %v", err)
	}
	defer database.Close()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	http.HandleFunc("/dummyLogin", api.DummyLoginHandler)

	// TODO: тут потом будем подключать роутеры с swagger.yaml

	log.Fatal(http.ListenAndServe(":8080", nil))
}
