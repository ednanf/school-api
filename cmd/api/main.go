package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	"github.com/ednanf/school-api/internal/domain"
	"github.com/ednanf/school-api/internal/pkg/hasher"
	"github.com/ednanf/school-api/internal/pkg/mailer"
	"github.com/ednanf/school-api/internal/pkg/token"
	"github.com/ednanf/school-api/internal/repository"
	"github.com/ednanf/school-api/internal/service"
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

	// Read JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("[FATAL] JWT_SECRET environment variable is missing")
	}

	// Read MailPit secrets from environment
	mailHost := os.Getenv("MAILPIT_HOST")
	mailPortStr := os.Getenv("MAILPIT_PORT")
	mailFrom := os.Getenv("MAIL_FROM")

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

	// Initialize MailPit service
	mailPort := 1025 // Default Mailpit SMTP port
	if p, err := strconv.Atoi(mailPortStr); err == nil {
		mailPort = p
	}

	mailService := mailer.NewMailpitService(mailHost, mailPort, mailFrom)

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Initialize repositories, services, and handlers
	studentRepo := repository.NewStudentRepository(db)
	studentService := service.NewStudentService(studentRepo)
	studentHandler := transportHttp.NewStudentHandler(studentService, validate)

	classRepo := repository.NewClassRepository(db)
	classService := service.NewClassService(classRepo)
	classHandler := transportHttp.NewClassHandler(classService, validate)

	subjectRepo := repository.NewSubjectRepository(db)
	subjectService := service.NewSubjectService(subjectRepo)
	subjectHandler := transportHttp.NewSubjectHandler(subjectService, validate)

	teacherRepo := repository.NewTeacherRepository(db)
	teacherService := service.NewTeacherService(teacherRepo)
	teacherHandler := transportHttp.NewTeacherHandler(teacherService, validate)

	teacherAssignmentRepo := repository.NewTeacherAssignmentRepository(db)
	teacherAssignmentService := service.NewTeacherAssignmentService(teacherAssignmentRepo)
	teacherAssignmentHandler := transportHttp.NewTeacherAssignmentHandler(teacherAssignmentService, validate)

	departmentRepo := repository.NewDepartmentRepository(db)
	departmentService := service.NewDepartmentService(departmentRepo)
	departmentHandler := transportHttp.NewDepartmentHandler(departmentService, validate)

	staffPositionRepo := repository.NewStaffPositionRepository(db)
	staffPositionService := service.NewStaffPositionService(staffPositionRepo)
	staffPositionHandler := transportHttp.NewStaffPositionHandler(staffPositionService, validate)

	staffRepo := repository.NewStaffRepository(db)
	staffService := service.NewStaffService(staffRepo)
	staffHandler := transportHttp.NewStaffHandler(staffService, validate)

	hasher := hasher.NewArgon2Hasher()
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, hasher)
	userHandler := transportHttp.NewUserHandler(userService, validate)

	jwtIssuer := "school-api"
	tokenService := token.NewJWTService(jwtSecret, jwtIssuer)
	authService := service.NewAuthService(userRepo, hasher, tokenService, mailService, 24*time.Hour)
	authHandler := transportHttp.NewAuthHandler(authService, validate)

	// Mount the routes under a versioned API prefix
	r.Route("/api/v1", func(r chi.Router) {
		// Mount ALL auth endpoints on /auth once
		r.Mount("/auth", authHandler.Routes(tokenService))

		// Protected Resource Routes (Requires valid JWT)
		r.Group(func(r chi.Router) {
			r.Use(transportHttp.AuthMiddleware(tokenService))

			// Admin-only user management
			r.Group(func(r chi.Router) {
				r.Use(transportHttp.RequireRole(domain.RoleAdmin))
				r.Mount("/users", userHandler.UserRoutes())
			})

			// Standard authenticated resource endpoints
			r.Mount("/students", studentHandler.StudentRoutes())
			r.Mount("/classes", classHandler.ClassRoutes())
			r.Mount("/subjects", subjectHandler.SubjectRoutes())
			r.Mount("/teachers", teacherHandler.TeacherRoutes())
			r.Mount("/teacher_assignments", teacherAssignmentHandler.TeacherAssignmentRoutes())
			r.Mount("/departments", departmentHandler.DepartmentRoutes())
			r.Mount("/staff_positions", staffPositionHandler.StaffPositionRoutes())
			r.Mount("/staff", staffHandler.StaffRoutes())
		})
	})

	log.Printf("[SYSTEM] Server running on port %s...", port)

	// Start the server
	err = http.ListenAndServe(port, r)
	if err != nil {
		fmt.Printf("[ERROR] Server failed to start: %v\n", err)
	}
}
