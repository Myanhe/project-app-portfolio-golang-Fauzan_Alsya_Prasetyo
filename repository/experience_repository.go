package repository

import (
	"context"
	"database/sql"
	"porto/model"
)

type ExperienceRepository interface {
	GetAll(ctx context.Context) ([]model.Experience, error)
	GetByID(ctx context.Context, id int) (*model.Experience, error)
	Create(ctx context.Context, e *model.Experience) error
	Update(ctx context.Context, e *model.Experience) error
	Delete(ctx context.Context, id int) error
}

type experienceRepository struct {
	db *sql.DB
}

func NewExperienceRepository(db *sql.DB) ExperienceRepository {
	return &experienceRepository{db}
}

func (r *experienceRepository) GetAll(ctx context.Context) ([]model.Experience, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, title, company, start_date, end_date, description FROM experiences")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var experiences []model.Experience
	for rows.Next() {
		var e model.Experience
		if err := rows.Scan(&e.ID, &e.Title, &e.Company, &e.StartDate, &e.EndDate, &e.Description); err != nil {
			return nil, err
		}
		experiences = append(experiences, e)
	}
	return experiences, nil
}

func (r *experienceRepository) GetByID(ctx context.Context, id int) (*model.Experience, error) {
	var e model.Experience
	err := r.db.QueryRowContext(ctx, "SELECT id, title, company, start_date, end_date, description FROM experiences WHERE id=$1", id).Scan(&e.ID, &e.Title, &e.Company, &e.StartDate, &e.EndDate, &e.Description)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *experienceRepository) Create(ctx context.Context, e *model.Experience) error {
	return r.db.QueryRowContext(ctx, "INSERT INTO experiences (title, company, start_date, end_date, description) VALUES ($1, $2, $3, $4, $5) RETURNING id", e.Title, e.Company, e.StartDate, e.EndDate, e.Description).Scan(&e.ID)
}

func (r *experienceRepository) Update(ctx context.Context, e *model.Experience) error {
	_, err := r.db.ExecContext(ctx, "UPDATE experiences SET title=$1, company=$2, start_date=$3, end_date=$4, description=$5 WHERE id=$6", e.Title, e.Company, e.StartDate, e.EndDate, e.Description, e.ID)
	return err
}

func (r *experienceRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM experiences WHERE id=$1", id)
	return err
}
