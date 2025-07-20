package service

import (
	"context"
	"log"
	"porto/model"
	"porto/repository"
	"porto/validation"
)

type ContactService interface {
	GetAll(ctx context.Context) ([]model.Contact, error)
	GetByID(ctx context.Context, id int) (*model.Contact, error)
	Create(ctx context.Context, c *model.Contact) error
	Delete(ctx context.Context, id int) error
}

type contactService struct {
	repo repository.ContactRepository
}

func NewContactService(repo repository.ContactRepository) ContactService {
	return &contactService{repo}
}

func (s *contactService) GetAll(ctx context.Context) ([]model.Contact, error) {
	return s.repo.GetAll(ctx)
}

func (s *contactService) GetByID(ctx context.Context, id int) (*model.Contact, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *contactService) Create(ctx context.Context, c *model.Contact) error {
	if err := validation.ValidateContact(c); err != nil {
		log.Printf("[ContactService] Create validation error: %v", err)
		return err
	}
	err := s.repo.Create(ctx, c)
	if err != nil {
		log.Printf("[ContactService] Create DB error: %v", err)
		return err
	}
	log.Printf("[ContactService] Created contact: %+v", c)
	return nil
}

func (s *contactService) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		log.Printf("[ContactService] Delete DB error: %v", err)
		return err
	}
	log.Printf("[ContactService] Deleted contact id: %d", id)
	return nil
}
