package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"porto/model"
	"porto/service"
	"strconv"

	"github.com/gorilla/mux"
)

type ExperienceHandler struct {
	Service service.ExperienceService
}

func NewExperienceHandler(s service.ExperienceService) *ExperienceHandler {
	return &ExperienceHandler{Service: s}
}

func (h *ExperienceHandler) GetExperiences(w http.ResponseWriter, r *http.Request) {
	exps, err := h.Service.GetAll(r.Context())
	if err != nil {
		log.Printf("[ExperienceHandler] GetExperiences error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("[ExperienceHandler] GetExperiences success, count: %d", len(exps))
	json.NewEncoder(w).Encode(exps)
}

func (h *ExperienceHandler) CreateExperience(w http.ResponseWriter, r *http.Request) {
	var e model.Experience
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		log.Printf("[ExperienceHandler] CreateExperience decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.Service.Create(r.Context(), &e); err != nil {
		log.Printf("[ExperienceHandler] CreateExperience service error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	log.Printf("[ExperienceHandler] CreateExperience success: %+v", e)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(e)
}

func (h *ExperienceHandler) UpdateExperience(w http.ResponseWriter, r *http.Request) {
	var e model.Experience
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		log.Printf("[ExperienceHandler] UpdateExperience decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.Service.Update(r.Context(), &e); err != nil {
		log.Printf("[ExperienceHandler] UpdateExperience service error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	log.Printf("[ExperienceHandler] UpdateExperience success: %+v", e)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(e)
}

func (h *ExperienceHandler) DeleteExperience(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("[ExperienceHandler] DeleteExperience invalid id: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.Service.Delete(r.Context(), id); err != nil {
		log.Printf("[ExperienceHandler] DeleteExperience service error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("[ExperienceHandler] DeleteExperience success, id: %d", id)
	w.WriteHeader(http.StatusNoContent)
}
