package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"porto/model"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPortfolioRepository_GetAll(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPortfolioRepository(db)

	// success
	rows := sqlmock.NewRows([]string{"id", "name", "description", "image_url", "link"}).
		AddRow(1, "A", "desc", "img", "link")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description, image_url, link FROM portfolios")).
		WillReturnRows(rows)
	result, err := repo.GetAll(context.Background())
	if err != nil || len(result) != 1 {
		t.Errorf("expected 1 result, got %v, err %v", result, err)
	}

	// error
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description, image_url, link FROM portfolios")).
		WillReturnError(sql.ErrConnDone)
	_, err = repo.GetAll(context.Background())
	if err == nil {
		t.Error("expected error")
	}
}

func TestPortfolioRepository_GetByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPortfolioRepository(db)

	// success
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description, image_url, link FROM portfolios WHERE id=$1")).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "image_url", "link"}).AddRow(1, "A", "desc", "img", "link"))
	_, err := repo.GetByID(context.Background(), 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// error
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, description, image_url, link FROM portfolios WHERE id=$1")).
		WithArgs(2).
		WillReturnError(sql.ErrNoRows)
	_, err = repo.GetByID(context.Background(), 2)
	if err == nil {
		t.Error("expected error")
	}
}

func TestPortfolioRepository_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPortfolioRepository(db)

	// success
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO portfolios (name, description, image_url, link) VALUES ($1, $2, $3, $4) RETURNING id")).
		WithArgs("A", "desc", "img", "link").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	p := &model.Portfolio{Name: "A", Description: "desc", ImageURL: "img", Link: "link"}
	err := repo.Create(context.Background(), p)
	if err != nil || p.ID != 1 {
		t.Errorf("expected id 1, got %v, err %v", p.ID, err)
	}

	// error
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO portfolios (name, description, image_url, link) VALUES ($1, $2, $3, $4) RETURNING id")).
		WithArgs("B", "desc", "img", "link").
		WillReturnError(sql.ErrConnDone)
	p2 := &model.Portfolio{Name: "B", Description: "desc", ImageURL: "img", Link: "link"}
	err = repo.Create(context.Background(), p2)
	if err == nil {
		t.Error("expected error")
	}
}

func TestPortfolioRepository_Update(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPortfolioRepository(db)

	// success
	mock.ExpectExec(regexp.QuoteMeta("UPDATE portfolios SET name=$1, description=$2, image_url=$3, link=$4 WHERE id=$5")).
		WithArgs("A", "desc", "img", "link", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	p := &model.Portfolio{ID: 1, Name: "A", Description: "desc", ImageURL: "img", Link: "link"}
	err := repo.Update(context.Background(), p)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// error
	mock.ExpectExec(regexp.QuoteMeta("UPDATE portfolios SET name=$1, description=$2, image_url=$3, link=$4 WHERE id=$5")).
		WithArgs("B", "desc", "img", "link", 2).
		WillReturnError(sql.ErrConnDone)
	p2 := &model.Portfolio{ID: 2, Name: "B", Description: "desc", ImageURL: "img", Link: "link"}
	err = repo.Update(context.Background(), p2)
	if err == nil {
		t.Error("expected error")
	}
}

func TestPortfolioRepository_Delete(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewPortfolioRepository(db)

	// success
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM portfolios WHERE id=$1")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err := repo.Delete(context.Background(), 1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// error
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM portfolios WHERE id=$1")).
		WithArgs(2).
		WillReturnError(sql.ErrConnDone)
	err = repo.Delete(context.Background(), 2)
	if err == nil {
		t.Error("expected error")
	}
}
