package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"porto/model"
	"porto/service"
	"strconv"

	"github.com/gorilla/mux"
)

type PortfolioHandler struct {
	Service     service.PortfolioService
	TemplateDir string
}

func NewPortfolioHandler(s service.PortfolioService, templateDir string) *PortfolioHandler {
	return &PortfolioHandler{Service: s, TemplateDir: templateDir}
}

func (h *PortfolioHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.Service.GetAll(r.Context())
	if err != nil {
		log.Printf("[PortfolioHandler] GetProjects error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("[PortfolioHandler] GetProjects success, count: %d", len(projects))
	json.NewEncoder(w).Encode(projects)
}

func (h *PortfolioHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var p model.Portfolio
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("[PortfolioHandler] CreateProject decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.Service.Create(r.Context(), &p); err != nil {
		log.Printf("[PortfolioHandler] CreateProject service error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	log.Printf("[PortfolioHandler] CreateProject success: %+v", p)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (h *PortfolioHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	var p model.Portfolio
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Printf("[PortfolioHandler] UpdateProject decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.Service.Update(r.Context(), &p); err != nil {
		log.Printf("[PortfolioHandler] UpdateProject service error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	log.Printf("[PortfolioHandler] UpdateProject success: %+v", p)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

func (h *PortfolioHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("[PortfolioHandler] DeleteProject invalid id: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.Service.Delete(r.Context(), id); err != nil {
		log.Printf("[PortfolioHandler] DeleteProject service error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("[PortfolioHandler] DeleteProject success, id: %d", id)
	w.WriteHeader(http.StatusNoContent)
}

// Render halaman daftar portfolio (HTML dinamis)
func (h *PortfolioHandler) RenderPortfolioPage(w http.ResponseWriter, r *http.Request) {
	projects, err := h.Service.GetAll(r.Context())
	if err != nil {
		log.Printf("[PortfolioHandler] RenderPortfolioPage error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tmplPath := filepath.Join(h.TemplateDir, "portfolio.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("[PortfolioHandler] template parse error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, projects)
	if err != nil {
		log.Printf("[PortfolioHandler] template execute error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type AboutData struct {
	Name     string
	Bio      string
	ImageURL string
}

func (h *PortfolioHandler) RenderAboutPage(w http.ResponseWriter, r *http.Request) {
	data := AboutData{
		Name:     "Fauzan Alsya Prasetyo",
		Bio:      "Saya adalah seorang software engineer yang berfokus pada pengembangan aplikasi web dengan arsitektur clean architecture dan Go.",
		ImageURL: "img/about-us.png",
	}
	tmplPath := filepath.Join(h.TemplateDir, "about.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Template error: " + err.Error()))
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Render error: " + err.Error()))
	}
}
