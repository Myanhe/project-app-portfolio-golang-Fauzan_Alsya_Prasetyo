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

type mockPortfolioService struct {
	GetAllFunc func(ctx context.Context) ([]model.Portfolio, error)
	CreateFunc func(ctx context.Context, p *model.Portfolio) error
}

func (m *mockPortfolioService) GetAll(ctx context.Context) ([]model.Portfolio, error) {
	return m.GetAllFunc(ctx)
}
func (m *mockPortfolioService) GetByID(ctx context.Context, id int) (*model.Portfolio, error) {
	return nil, nil
}
func (m *mockPortfolioService) Create(ctx context.Context, p *model.Portfolio) error {
	return m.CreateFunc(ctx, p)
}
func (m *mockPortfolioService) Update(ctx context.Context, _ *model.Portfolio) error { return nil }
func (m *mockPortfolioService) Delete(ctx context.Context, _ int) error              { return nil }

func TestPortfolioHandler_GetProjects(t *testing.T) {
	svc := &mockPortfolioService{
		GetAllFunc: func(ctx context.Context) ([]model.Portfolio, error) {
			return []model.Portfolio{{ID: 1, Name: "A"}}, nil
		},
	}
	h := NewPortfolioHandler(svc, "")
	r := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	w := httptest.NewRecorder()

	h.GetProjects(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestPortfolioHandler_CreateProject(t *testing.T) {
	svc := &mockPortfolioService{
		CreateFunc: func(ctx context.Context, p *model.Portfolio) error {
			if p.Name == "" {
				return errors.New("invalid input")
			}
			return nil
		},
	}
	h := NewPortfolioHandler(svc, "")
	body, _ := json.Marshal(model.Portfolio{Name: "A", Description: "B"})
	r := httptest.NewRequest(http.MethodPost, "/api/projects", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateProject(w, r)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestPortfolioHandler_GetProjects_Error(t *testing.T) {
	svc := &mockPortfolioService{
		GetAllFunc: func(ctx context.Context) ([]model.Portfolio, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewPortfolioHandler(svc, "")
	r := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	w := httptest.NewRecorder()

	h.GetProjects(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestPortfolioHandler_CreateProject_BadRequest(t *testing.T) {
	svc := &mockPortfolioService{
		CreateFunc: func(ctx context.Context, p *model.Portfolio) error {
			return nil
		},
	}
	h := NewPortfolioHandler(svc, "")
	r := httptest.NewRequest(http.MethodPost, "/api/projects", bytes.NewReader([]byte("notjson")))
	w := httptest.NewRecorder()

	h.CreateProject(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestPortfolioHandler_CreateProject_ServiceError(t *testing.T) {
	svc := &mockPortfolioService{
		CreateFunc: func(ctx context.Context, p *model.Portfolio) error {
			return errors.New("service error")
		},
	}
	h := NewPortfolioHandler(svc, "")
	body, _ := json.Marshal(model.Portfolio{Name: "A", Description: "B"})
	r := httptest.NewRequest(http.MethodPost, "/api/projects", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateProject(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
