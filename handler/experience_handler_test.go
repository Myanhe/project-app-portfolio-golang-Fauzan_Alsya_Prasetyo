package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"porto/model"
)

type mockExperienceService struct {
	GetAllFunc func(ctx context.Context) ([]model.Experience, error)
	CreateFunc func(ctx context.Context, e *model.Experience) error
}

func (m *mockExperienceService) GetAll(ctx context.Context) ([]model.Experience, error) {
	return m.GetAllFunc(ctx)
}
func (m *mockExperienceService) GetByID(ctx context.Context, id int) (*model.Experience, error) {
	return nil, nil
}
func (m *mockExperienceService) Create(ctx context.Context, e *model.Experience) error {
	return m.CreateFunc(ctx, e)
}
func (m *mockExperienceService) Update(ctx context.Context, _ *model.Experience) error { return nil }
func (m *mockExperienceService) Delete(ctx context.Context, _ int) error               { return nil }

func TestExperienceHandler_GetExperiences(t *testing.T) {
	svc := &mockExperienceService{
		GetAllFunc: func(ctx context.Context) ([]model.Experience, error) {
			return []model.Experience{{ID: 1, Title: "A"}}, nil
		},
	}
	h := NewExperienceHandler(svc)
	r := httptest.NewRequest(http.MethodGet, "/api/experiences", nil)
	w := httptest.NewRecorder()

	h.GetExperiences(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestExperienceHandler_CreateExperience(t *testing.T) {
	svc := &mockExperienceService{
		CreateFunc: func(ctx context.Context, e *model.Experience) error {
			if e.Title == "" {
				return errors.New("invalid input")
			}
			return nil
		},
	}
	h := NewExperienceHandler(svc)
	body, _ := json.Marshal(model.Experience{Title: "A", Company: "B", StartDate: "2020"})
	r := httptest.NewRequest(http.MethodPost, "/api/experiences", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateExperience(w, r)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestExperienceHandler_GetExperiences_Error(t *testing.T) {
	svc := &mockExperienceService{
		GetAllFunc: func(ctx context.Context) ([]model.Experience, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewExperienceHandler(svc)
	r := httptest.NewRequest(http.MethodGet, "/api/experiences", nil)
	w := httptest.NewRecorder()

	h.GetExperiences(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestExperienceHandler_CreateExperience_BadRequest(t *testing.T) {
	svc := &mockExperienceService{
		CreateFunc: func(ctx context.Context, e *model.Experience) error {
			return nil
		},
	}
	h := NewExperienceHandler(svc)
	r := httptest.NewRequest(http.MethodPost, "/api/experiences", bytes.NewReader([]byte("notjson")))
	w := httptest.NewRecorder()

	h.CreateExperience(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestExperienceHandler_CreateExperience_ServiceError(t *testing.T) {
	svc := &mockExperienceService{
		CreateFunc: func(ctx context.Context, e *model.Experience) error {
			return errors.New("service error")
		},
	}
	h := NewExperienceHandler(svc)
	body, _ := json.Marshal(model.Experience{Title: "A", Company: "B", StartDate: "2020"})
	r := httptest.NewRequest(http.MethodPost, "/api/experiences", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateExperience(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
