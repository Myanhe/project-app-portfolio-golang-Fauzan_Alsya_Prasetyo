package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"

	"porto/handler"
	"porto/repository"
	"porto/service"
)

func main() {
	// Load environment variables or config here
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/portfolio?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Repository, Service, Handler wiring
	portfolioRepo := repository.NewPortfolioRepository(db)
	experienceRepo := repository.NewExperienceRepository(db)
	contactRepo := repository.NewContactRepository(db)

	portfolioService := service.NewPortfolioService(portfolioRepo)
	experienceService := service.NewExperienceService(experienceRepo)
	contactService := service.NewContactService(contactRepo)

	portfolioHandler := handler.NewPortfolioHandler(portfolioService, "webview")
	experienceHandler := handler.NewExperienceHandler(experienceService)
	contactHandler := handler.NewContactHandler(contactService)

	r := chi.NewRouter()

	// Portfolio endpoints
	r.Get("/api/projects", portfolioHandler.GetProjects)
	r.Post("/api/projects", portfolioHandler.CreateProject)
	r.Put("/api/projects", portfolioHandler.UpdateProject)
	r.Delete("/api/projects/{id}", portfolioHandler.DeleteProject)

	// Experience endpoints
	r.Get("/api/experiences", experienceHandler.GetExperiences)
	r.Post("/api/experiences", experienceHandler.CreateExperience)
	r.Put("/api/experiences", experienceHandler.UpdateExperience)
	r.Delete("/api/experiences/{id}", experienceHandler.DeleteExperience)

	// Contact endpoints
	r.Get("/api/contacts", contactHandler.GetContacts)
	r.Post("/api/contacts", contactHandler.CreateContact)
	r.Delete("/api/contacts/{id}", contactHandler.DeleteContact)

	// Static file serving (optional)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("webview/static"))))

	// HTML template routes
	r.Get("/portfolio", portfolioHandler.RenderPortfolioPage)
	r.Get("/about", portfolioHandler.RenderAboutPage)

	log.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
