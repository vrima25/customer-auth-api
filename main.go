package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/vrima25/go-auth-service/config"
	"github.com/vrima25/go-auth-service/controller"
	"github.com/vrima25/go-auth-service/middleware"
	"github.com/vrima25/go-auth-service/repository"
	"github.com/vrima25/go-auth-service/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on system environtment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("connected to database")

	customerRepo := repository.NewCustomerRepository(db)
	authService := service.NewAuthService(customerRepo, cfg.JWTSecret)
	authController := controller.NewAuthController(authService)
	healthController := controller.NewHealthController(db)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthController.Health)
	mux.HandleFunc("/api/customers/register", authController.Register)
	mux.HandleFunc("/api/customers/login", authController.Login)
	mux.Handle(
		"/api/customers/me",
		middleware.AuthMiddleware(cfg.JWTSecret)(http.HandlerFunc(authController.Me)),
	)

	log.Printf("server running on %s", cfg.Addr())

	handler := corsMiddleware(mux)
	if err := http.ListenAndServe(cfg.Addr(), handler); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
