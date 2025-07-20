package repository

import (
	"context"
	"database/sql"
	"porto/model"
)

type PortfolioRepository interface {
	GetAll(ctx context.Context) ([]model.Portfolio, error)
	GetByID(ctx context.Context, id int) (*model.Portfolio, error)
	Create(ctx context.Context, p *model.Portfolio) error
	Update(ctx context.Context, p *model.Portfolio) error
	Delete(ctx context.Context, id int) error
}

type portfolioRepository struct {
	db *sql.DB
}

func NewPortfolioRepository(db *sql.DB) PortfolioRepository {
	return &portfolioRepository{db}
}

func (r *portfolioRepository) GetAll(ctx context.Context) ([]model.Portfolio, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, description, image_url, link FROM portfolios")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var portfolios []model.Portfolio
	for rows.Next() {
		var p model.Portfolio
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.ImageURL, &p.Link); err != nil {
			return nil, err
		}
		portfolios = append(portfolios, p)
	}
	return portfolios, nil
}

func (r *portfolioRepository) GetByID(ctx context.Context, id int) (*model.Portfolio, error) {
	var p model.Portfolio
	err := r.db.QueryRowContext(ctx, "SELECT id, name, description, image_url, link FROM portfolios WHERE id=$1", id).Scan(&p.ID, &p.Name, &p.Description, &p.ImageURL, &p.Link)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *portfolioRepository) Create(ctx context.Context, p *model.Portfolio) error {
	return r.db.QueryRowContext(ctx, "INSERT INTO portfolios (name, description, image_url, link) VALUES ($1, $2, $3, $4) RETURNING id", p.Name, p.Description, p.ImageURL, p.Link).Scan(&p.ID)
}

func (r *portfolioRepository) Update(ctx context.Context, p *model.Portfolio) error {
	_, err := r.db.ExecContext(ctx, "UPDATE portfolios SET name=$1, description=$2, image_url=$3, link=$4 WHERE id=$5", p.Name, p.Description, p.ImageURL, p.Link, p.ID)
	return err
}

func (r *portfolioRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM portfolios WHERE id=$1", id)
	return err
}
