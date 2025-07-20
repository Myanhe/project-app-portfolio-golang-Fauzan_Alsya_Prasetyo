package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"porto/model"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestExperienceRepository_GetAll(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewExperienceRepository(db)

	// success
	rows := sqlmock.NewRows([]string{"id", "title", "company", "start_date", "end_date", "description"}).
		AddRow(1, "A", "B", "2020", "2021", "desc")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, company, start_date, end_date, description FROM experiences")).
		WillReturnRows(rows)
	result, err := repo.GetAll(context.Background())
	if err != nil || len(result) != 1 {
		t.Errorf("expected 1 result, got %v, err %v", result, err)
	}

	// error
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, company, start_date, end_date, description FROM experiences")).
		WillReturnError(sql.ErrConnDone)
	_, err = repo.GetAll(context.Background())
	if err == nil {
		t.Error("expected error")
	}
}

func TestExperienceRepository_GetByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewExperienceRepository(db)

	// success
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, company, start_date, end_date, description FROM experiences WHERE id=$1")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "company", "start_date", "end_date", "description"}).AddRow(1, "A", "B", "2020", "2021", "desc"))
	_, err := repo.GetByID(context.Background(), 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// error
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title, company, start_date, end_date, description FROM experiences WHERE id=$1")).
		WithArgs(2).
		WillReturnError(sql.ErrNoRows)
	_, err = repo.GetByID(context.Background(), 2)
	if err == nil {
		t.Error("expected error")
	}
}

func TestExperienceRepository_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewExperienceRepository(db)

	// success
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO experiences (title, company, start_date, end_date, description) VALUES ($1, $2, $3, $4, $5) RETURNING id")).
		WithArgs("A", "B", "2020", "2021", "desc").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	e := &model.Experience{Title: "A", Company: "B", StartDate: "2020", EndDate: "2021", Description: "desc"}
	err := repo.Create(context.Background(), e)
	if err != nil || e.ID != 1 {
		t.Errorf("expected id 1, got %v, err %v", e.ID, err)
	}

	// error
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO experiences (title, company, start_date, end_date, description) VALUES ($1, $2, $3, $4, $5) RETURNING id")).
		WithArgs("B", "C", "2020", "2021", "desc").
		WillReturnError(sql.ErrConnDone)
	e2 := &model.Experience{Title: "B", Company: "C", StartDate: "2020", EndDate: "2021", Description: "desc"}
	err = repo.Create(context.Background(), e2)
	if err == nil {
		t.Error("expected error")
	}
}

func TestExperienceRepository_Update(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewExperienceRepository(db)

	// success
	mock.ExpectExec(regexp.QuoteMeta("UPDATE experiences SET title=$1, company=$2, start_date=$3, end_date=$4, description=$5 WHERE id=$6")).
		WithArgs("A", "B", "2020", "2021", "desc", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	e := &model.Experience{ID: 1, Title: "A", Company: "B", StartDate: "2020", EndDate: "2021", Description: "desc"}
	err := repo.Update(context.Background(), e)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// error
	mock.ExpectExec(regexp.QuoteMeta("UPDATE experiences SET title=$1, company=$2, start_date=$3, end_date=$4, description=$5 WHERE id=$6")).
		WithArgs("B", "C", "2020", "2021", "desc", 2).
		WillReturnError(sql.ErrConnDone)
	e2 := &model.Experience{ID: 2, Title: "B", Company: "C", StartDate: "2020", EndDate: "2021", Description: "desc"}
	err = repo.Update(context.Background(), e2)
	if err == nil {
		t.Error("expected error")
	}
}

func TestExperienceRepository_Delete(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewExperienceRepository(db)

	// success
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM experiences WHERE id=$1")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err := repo.Delete(context.Background(), 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// error
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM experiences WHERE id=$1")).
		WithArgs(2).
		WillReturnError(sql.ErrConnDone)
	err = repo.Delete(context.Background(), 2)
	if err == nil {
		t.Error("expected error")
	}
}
