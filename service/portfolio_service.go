package service

import (
	"context"
	"errors"
	"log"
	"porto/model"
	"porto/repository"
	"porto/validation"
)

type PortfolioService interface {
	GetAll(ctx context.Context) ([]model.Portfolio, error)
	GetByID(ctx context.Context, id int) (*model.Portfolio, error)
	Create(ctx context.Context, p *model.Portfolio) error
	Update(ctx context.Context, p *model.Portfolio) error
	Delete(ctx context.Context, id int) error
}

type portfolioService struct {
	repo repository.PortfolioRepository
}

func NewPortfolioService(repo repository.PortfolioRepository) PortfolioService {
	return &portfolioService{repo}
}

func (s *portfolioService) GetAll(ctx context.Context) ([]model.Portfolio, error) {
	return s.repo.GetAll(ctx)
}

func (s *portfolioService) GetByID(ctx context.Context, id int) (*model.Portfolio, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *portfolioService) Create(ctx context.Context, p *model.Portfolio) error {
	if err := validation.ValidatePortfolio(p); err != nil {
		log.Printf("[PortfolioService] Create validation error: %v", err)
		return err
	}
	err := s.repo.Create(ctx, p)
	if err != nil {
		log.Printf("[PortfolioService] Create DB error: %v", err)
		return err
	}
	log.Printf("[PortfolioService] Created portfolio: %+v", p)
	return nil
}

func (s *portfolioService) Update(ctx context.Context, p *model.Portfolio) error {
	if p.ID == 0 {
		log.Printf("[PortfolioService] Update error: id is required")
		return errors.New("id is required")
	}
	if err := validation.ValidatePortfolio(p); err != nil {
		log.Printf("[PortfolioService] Update validation error: %v", err)
		return err
	}
	err := s.repo.Update(ctx, p)
	if err != nil {
		log.Printf("[PortfolioService] Update DB error: %v", err)
		return err
	}
	log.Printf("[PortfolioService] Updated portfolio: %+v", p)
	return nil
}

func (s *portfolioService) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		log.Printf("[PortfolioService] Delete DB error: %v", err)
		return err
	}
	log.Printf("[PortfolioService] Deleted portfolio id: %d", id)
	return nil
}
