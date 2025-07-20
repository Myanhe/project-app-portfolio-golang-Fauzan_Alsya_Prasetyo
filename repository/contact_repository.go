package repository

import (
	"context"
	"database/sql"
	"porto/model"
)

type ContactRepository interface {
	GetAll(ctx context.Context) ([]model.Contact, error)
	GetByID(ctx context.Context, id int) (*model.Contact, error)
	Create(ctx context.Context, c *model.Contact) error
	Delete(ctx context.Context, id int) error
}

type contactRepository struct {
	db *sql.DB
}

func NewContactRepository(db *sql.DB) ContactRepository {
	return &contactRepository{db}
}

func (r *contactRepository) GetAll(ctx context.Context) ([]model.Contact, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, email, message FROM contacts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var contacts []model.Contact
	for rows.Next() {
		var c model.Contact
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Message); err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}

func (r *contactRepository) GetByID(ctx context.Context, id int) (*model.Contact, error) {
	var c model.Contact
	err := r.db.QueryRowContext(ctx, "SELECT id, name, email, message FROM contacts WHERE id=$1", id).Scan(&c.ID, &c.Name, &c.Email, &c.Message)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *contactRepository) Create(ctx context.Context, c *model.Contact) error {
	return r.db.QueryRowContext(ctx, "INSERT INTO contacts (name, email, message) VALUES ($1, $2, $3) RETURNING id", c.Name, c.Email, c.Message).Scan(&c.ID)
}

func (r *contactRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM contacts WHERE id=$1", id)
	return err
}
