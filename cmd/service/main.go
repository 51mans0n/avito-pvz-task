package main

import (
	"fmt"
	"github.com/51mans0n/avito-pvz-task/internal/api"
	"github.com/51mans0n/avito-pvz-task/internal/db"
	grpcserver "github.com/51mans0n/avito-pvz-task/internal/grpc"
	"github.com/51mans0n/avito-pvz-task/internal/logging"
	"github.com/51mans0n/avito-pvz-task/internal/metrics"
	pvz_v1 "github.com/51mans0n/avito-pvz-task/pkg/proto/pvz/v1"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

func main() {

	logging.Init(false) // dev режим
	defer logging.Sync()

	fmt.Println("Starting HTTP service on :8080...")

	database, err := db.InitDB()
	if err != nil {
		logging.S().Fatalf("failed to init DB: %v", err)
	}
	defer database.Close()

	repo := db.NewRepo(database)

	go func() {
		lis, _ := net.Listen("tcp", ":3000")
		if err != nil {
			logging.S().Fatalw("listen gRPC", "err", err)
		}

		g := grpc.NewServer()
		pvz_v1.RegisterPVZServiceServer(g, grpcserver.New(repo))

		logging.S().Infow("gRPC :3000 started")
		if err := g.Serve(lis); err != nil {
			logging.S().Fatalw("gRPC:", err)
		}
	}()

	metrics.MustRegister()
	r := chi.NewRouter()
	r.Use(logging.RequestLogger)
	r.Use(metrics.PromMiddleware)

	// prometheus
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	// Health-check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Dummy login (не требует авторизации)
	r.Post("/dummyLogin", api.DummyLoginHandler)

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

	logging.S().Infow("HTTP started", "addr", ":8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logging.S().Fatalw("serve HTTP", "err", err)
	}
}
