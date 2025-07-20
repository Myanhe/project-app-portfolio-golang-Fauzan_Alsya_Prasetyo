package service

import (
	"context"
	"errors"
	"porto/model"
	"testing"
)

type mockExperienceRepo struct {
	CreateFunc func(ctx context.Context, e *model.Experience) error
	UpdateFunc func(ctx context.Context, e *model.Experience) error
	DeleteFunc func(ctx context.Context, id int) error
}

func (m *mockExperienceRepo) GetAll(ctx context.Context) ([]model.Experience, error) { return nil, nil }
func (m *mockExperienceRepo) GetByID(ctx context.Context, id int) (*model.Experience, error) {
	return nil, nil
}
func (m *mockExperienceRepo) Create(ctx context.Context, e *model.Experience) error {
	return m.CreateFunc(ctx, e)
}
func (m *mockExperienceRepo) Update(ctx context.Context, e *model.Experience) error {
	return m.UpdateFunc(ctx, e)
}
func (m *mockExperienceRepo) Delete(ctx context.Context, id int) error {
	return m.DeleteFunc(ctx, id)
}

func TestExperienceService_Create(t *testing.T) {
	repo := &mockExperienceRepo{
		CreateFunc: func(ctx context.Context, e *model.Experience) error {
			if e.Title == "exists" {
				return errors.New("duplicate")
			}
			return nil
		},
	}
	svc := NewExperienceService(repo)

	// valid
	e := &model.Experience{Title: "A", Company: "B", StartDate: "2020"}
	err := svc.Create(context.Background(), e)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// invalid
	e2 := &model.Experience{Title: "", Company: "B", StartDate: "2020"}
	err = svc.Create(context.Background(), e2)
	if err == nil {
		t.Errorf("expected error for empty title")
	}

	// duplicate
	e3 := &model.Experience{Title: "exists", Company: "B", StartDate: "2020"}
	err = svc.Create(context.Background(), e3)
	if err == nil || err.Error() != "duplicate" {
		t.Errorf("expected duplicate error, got %v", err)
	}
}

func TestExperienceService_Create_Invalid(t *testing.T) {
	repo := &mockExperienceRepo{CreateFunc: func(ctx context.Context, e *model.Experience) error { return nil }}
	svc := NewExperienceService(repo)
	e := &model.Experience{Title: "", Company: "", StartDate: ""}
	err := svc.Create(context.Background(), e)
	if err == nil {
		t.Error("expected error for invalid input")
	}
}

func TestExperienceService_Update_Invalid(t *testing.T) {
	repo := &mockExperienceRepo{}
	svc := NewExperienceService(repo)
	e := &model.Experience{ID: 0, Title: "A", Company: "B", StartDate: "2020"}
	err := svc.Update(context.Background(), e)
	if err == nil {
		t.Error("expected error for missing id")
	}
}

func TestExperienceService_Update_Success(t *testing.T) {
	repo := &mockExperienceRepo{UpdateFunc: func(ctx context.Context, e *model.Experience) error { return nil }}
	svc := NewExperienceService(repo)
	e := &model.Experience{ID: 1, Title: "A", Company: "B", StartDate: "2020"}
	err := svc.Update(context.Background(), e)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestExperienceService_Update_RepoError(t *testing.T) {
	repo := &mockExperienceRepo{UpdateFunc: func(ctx context.Context, e *model.Experience) error { return errors.New("db error") }}
	svc := NewExperienceService(repo)
	e := &model.Experience{ID: 1, Title: "A", Company: "B", StartDate: "2020"}
	err := svc.Update(context.Background(), e)
	if err == nil {
		t.Error("expected error from repo")
	}
}

func TestExperienceService_Delete_Success(t *testing.T) {
	repo := &mockExperienceRepo{DeleteFunc: func(ctx context.Context, id int) error { return nil }}
	svc := NewExperienceService(repo)
	err := svc.Delete(context.Background(), 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestExperienceService_Delete_RepoError(t *testing.T) {
	repo := &mockExperienceRepo{DeleteFunc: func(ctx context.Context, id int) error { return errors.New("db error") }}
	svc := NewExperienceService(repo)
	err := svc.Delete(context.Background(), 1)
	if err == nil {
		t.Error("expected error from repo")
	}
}
