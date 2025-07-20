package service

import (
	"context"
	"errors"
	"porto/model"
	"testing"
)

type mockPortfolioRepo struct {
	CreateFunc func(ctx context.Context, p *model.Portfolio) error
}

func (m *mockPortfolioRepo) GetAll(ctx context.Context) ([]model.Portfolio, error) { return nil, nil }
func (m *mockPortfolioRepo) GetByID(ctx context.Context, id int) (*model.Portfolio, error) {
	return nil, nil
}
func (m *mockPortfolioRepo) Create(ctx context.Context, p *model.Portfolio) error {
	return m.CreateFunc(ctx, p)
}
func (m *mockPortfolioRepo) Update(ctx context.Context, p *model.Portfolio) error { return nil }
func (m *mockPortfolioRepo) Delete(ctx context.Context, id int) error             { return nil }

func TestPortfolioService_Create(t *testing.T) {
	repo := &mockPortfolioRepo{
		CreateFunc: func(ctx context.Context, p *model.Portfolio) error {
			if p.Name == "exists" {
				return errors.New("duplicate")
			}
			return nil
		},
	}
	svc := NewPortfolioService(repo)

	// valid
	p := &model.Portfolio{Name: "A", Description: "B"}
	err := svc.Create(context.Background(), p)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// invalid
	p2 := &model.Portfolio{Name: "", Description: "B"}
	err = svc.Create(context.Background(), p2)
	if err == nil {
		t.Errorf("expected error for empty name")
	}

	// duplicate
	p3 := &model.Portfolio{Name: "exists", Description: "B"}
	err = svc.Create(context.Background(), p3)
	if err == nil || err.Error() != "duplicate" {
		t.Errorf("expected duplicate error, got %v", err)
	}
}

func TestPortfolioService_Create_Invalid(t *testing.T) {
	repo := &mockPortfolioRepo{CreateFunc: func(ctx context.Context, p *model.Portfolio) error { return nil }}
	svc := NewPortfolioService(repo)
	p := &model.Portfolio{Name: "", Description: ""}
	err := svc.Create(context.Background(), p)
	if err == nil {
		t.Error("expected error for invalid input")
	}
}

func TestPortfolioService_Update_Invalid(t *testing.T) {
	repo := &mockPortfolioRepo{}
	svc := NewPortfolioService(repo)
	p := &model.Portfolio{ID: 0, Name: "A", Description: "B"}
	err := svc.Update(context.Background(), p)
	if err == nil {
		t.Error("expected error for missing id")
	}
}
