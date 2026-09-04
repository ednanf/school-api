package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	"github.com/ednanf/school-api/internal/repository"
	transportHttp "github.com/ednanf/school-api/internal/transport/http"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("[ERROR]", err)
		return
	}

	port := os.Getenv("API_PORT")
	dbUsername := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Database DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUsername, dbPassword, dbHost, dbPort, dbName)

	// Open DB connection pool
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MariaDB: %v", err)
	}
	defer db.Close()

	log.Printf("[SYSTEM] Connected to %s...", dbName)

	// Configure pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	// Initialize a Router
	r := chi.NewRouter()

	// Initialize validator to be passed to handlers
	validate := validator.New(validator.WithRequiredStructEnabled())

	// Middlewares
	// RequestID injects a unique ID into each request for tracing
	r.Use(middleware.RequestID)
	// Logger prints request logs to the console
	r.Use(middleware.Logger)
	// Recoverer catches panics so the server doesn't crash from a single bad request
	r.Use(middleware.Recoverer)

	// Initialize repository (struct that encapsulates all database access logic) and handlers
	studentRepo := repository.NewStudentRepository(db)
	studentHandler := transportHttp.NewStudentHandler(studentRepo, validate)

	classRepo := repository.NewClassRepository(db)
	classHandler := transportHttp.NewClassHandler(classRepo, validate)

	// Mount the routes under a versioned API prefix
	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/students", studentHandler.StudentRoutes())
		r.Mount("/classes", classHandler.ClassRoutes())
	})

	log.Printf("[SYSTEM] Server running on port %s...", port)

	// Start the server
	err = http.ListenAndServe(port, r)
	if err != nil {
		fmt.Printf("[ERROR] Server failed to start: %v\n", err)
	}
}
