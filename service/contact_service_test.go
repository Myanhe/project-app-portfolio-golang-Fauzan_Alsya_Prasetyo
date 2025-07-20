package service

import (
	"context"
	"errors"
	"porto/model"
	"testing"
)

type mockContactRepo struct {
	CreateFunc func(ctx context.Context, c *model.Contact) error
}

func (m *mockContactRepo) GetAll(ctx context.Context) ([]model.Contact, error) { return nil, nil }
func (m *mockContactRepo) GetByID(ctx context.Context, id int) (*model.Contact, error) {
	return nil, nil
}
func (m *mockContactRepo) Create(ctx context.Context, c *model.Contact) error {
	return m.CreateFunc(ctx, c)
}
func (m *mockContactRepo) Delete(ctx context.Context, id int) error { return nil }

func TestContactService_Create(t *testing.T) {
	repo := &mockContactRepo{
		CreateFunc: func(ctx context.Context, c *model.Contact) error {
			if c.Email == "exists@mail.com" {
				return errors.New("duplicate")
			}
			return nil
		},
	}
	svc := NewContactService(repo)

	// valid
	c := &model.Contact{Name: "A", Email: "a@mail.com", Message: "hi"}
	err := svc.Create(context.Background(), c)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// invalid
	c2 := &model.Contact{Name: "", Email: "a@mail.com", Message: "hi"}
	err = svc.Create(context.Background(), c2)
	if err == nil {
		t.Errorf("expected error for empty name")
	}

	// duplicate
	c3 := &model.Contact{Name: "A", Email: "exists@mail.com", Message: "hi"}
	err = svc.Create(context.Background(), c3)
	if err == nil || err.Error() != "duplicate" {
		t.Errorf("expected duplicate error, got %v", err)
	}
}

func TestContactService_Create_Invalid(t *testing.T) {
	repo := &mockContactRepo{CreateFunc: func(ctx context.Context, c *model.Contact) error { return nil }}
	svc := NewContactService(repo)
	c := &model.Contact{Name: "", Email: "", Message: ""}
	err := svc.Create(context.Background(), c)
	if err == nil {
		t.Error("expected error for invalid input")
	}
}
