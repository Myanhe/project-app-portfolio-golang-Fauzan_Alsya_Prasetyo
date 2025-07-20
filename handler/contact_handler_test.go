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

type mockContactService struct {
	GetAllFunc func(ctx context.Context) ([]model.Contact, error)
	CreateFunc func(ctx context.Context, c *model.Contact) error
}

func (m *mockContactService) GetAll(ctx context.Context) ([]model.Contact, error) {
	return m.GetAllFunc(ctx)
}
func (m *mockContactService) GetByID(ctx context.Context, id int) (*model.Contact, error) {
	return nil, nil
}
func (m *mockContactService) Create(ctx context.Context, c *model.Contact) error {
	return m.CreateFunc(ctx, c)
}
func (m *mockContactService) Delete(ctx context.Context, _ int) error { return nil }

func TestContactHandler_GetContacts(t *testing.T) {
	svc := &mockContactService{
		GetAllFunc: func(ctx context.Context) ([]model.Contact, error) {
			return []model.Contact{{ID: 1, Name: "A"}}, nil
		},
	}
	h := NewContactHandler(svc)
	r := httptest.NewRequest(http.MethodGet, "/api/contacts", nil)
	w := httptest.NewRecorder()

	h.GetContacts(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestContactHandler_CreateContact(t *testing.T) {
	svc := &mockContactService{
		CreateFunc: func(ctx context.Context, c *model.Contact) error {
			if c.Name == "" {
				return errors.New("invalid input")
			}
			return nil
		},
	}
	h := NewContactHandler(svc)
	body, _ := json.Marshal(model.Contact{Name: "A", Email: "a@mail.com", Message: "hi"})
	r := httptest.NewRequest(http.MethodPost, "/api/contacts", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateContact(w, r)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestContactHandler_GetContacts_Error(t *testing.T) {
	svc := &mockContactService{
		GetAllFunc: func(ctx context.Context) ([]model.Contact, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewContactHandler(svc)
	r := httptest.NewRequest(http.MethodGet, "/api/contacts", nil)
	w := httptest.NewRecorder()

	h.GetContacts(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestContactHandler_CreateContact_BadRequest(t *testing.T) {
	svc := &mockContactService{
		CreateFunc: func(ctx context.Context, c *model.Contact) error {
			return nil
		},
	}
	h := NewContactHandler(svc)
	r := httptest.NewRequest(http.MethodPost, "/api/contacts", bytes.NewReader([]byte("notjson")))
	w := httptest.NewRecorder()

	h.CreateContact(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestContactHandler_CreateContact_ServiceError(t *testing.T) {
	svc := &mockContactService{
		CreateFunc: func(ctx context.Context, c *model.Contact) error {
			return errors.New("service error")
		},
	}
	h := NewContactHandler(svc)
	body, _ := json.Marshal(model.Contact{Name: "A", Email: "a@mail.com", Message: "hi"})
	r := httptest.NewRequest(http.MethodPost, "/api/contacts", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateContact(w, r)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
