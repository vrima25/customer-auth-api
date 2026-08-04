package main

import (
	"database/sql"
	"log"
	"net/http"
	"ticket-triage-api/config"
	"ticket-triage-api/controller"
	"ticket-triage-api/middleware"
	"ticket-triage-api/repository"
	"ticket-triage-api/service"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main(){
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on system environtment variables")
	}

	cfg := config.Load()
	log.Printf("DEBUG — DSN: %s", cfg.DSN())

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is not set - check your .env file")
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

	mux := http.NewServeMux()
	mux.HandleFunc("/api/customers/register", authController.Register)
	mux.HandleFunc("/api/customers/login", authController.Login)
	mux.Handle(
		"/api/customers/me",
		middleware.AuthMiddleware(cfg.JWTSecret)(http.HandlerFunc(authController.Me)),
	)

	log.Println("server running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}