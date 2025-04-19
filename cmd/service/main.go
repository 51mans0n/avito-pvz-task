package main

import (
	"fmt"
	"github.com/51mans0n/avito-pvz-task/internal/api"
	"github.com/51mans0n/avito-pvz-task/internal/db"
	"github.com/go-chi/chi/v5"
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

	repo := db.NewRepo(database)

	r := chi.NewRouter()

	// Health-check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Dummy login (не требует авторизации)
	r.Get("/dummyLogin", api.DummyLoginHandler)

	// Endpoints, защищённые AuthMiddleware
	r.Group(func(sub chi.Router) {
		sub.Use(api.AuthMiddleware) // каждый запрос внутри sub будет проходить AuthMiddleware

		// /pvz
		sub.Route("/pvz", func(rpvz chi.Router) {
			// POST /pvz -> Create
			rpvz.Post("/", api.CreatePVZHandler(repo))

			// GET /pvz -> List
			rpvz.Get("/", api.GetPVZListHandler(repo))
			rpvz.Post("/{pvzId}/delete_last_product", api.DeleteLastProductHandler(repo))
			rpvz.Post("/{pvzId}/close_last_reception", api.CloseLastReceptionHandler(repo))
		})

		// /receptions
		sub.Post("/receptions", api.CreateReceptionHandler(repo))

		// /products
		sub.Post("/products", api.CreateProductHandler(repo))
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
