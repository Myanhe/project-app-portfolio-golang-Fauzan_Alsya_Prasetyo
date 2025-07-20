package service

import (
	"context"
	"errors"
	"log"
	"porto/model"
	"porto/repository"
	"porto/validation"
)

type ExperienceService interface {
	GetAll(ctx context.Context) ([]model.Experience, error)
	GetByID(ctx context.Context, id int) (*model.Experience, error)
	Create(ctx context.Context, e *model.Experience) error
	Update(ctx context.Context, e *model.Experience) error
	Delete(ctx context.Context, id int) error
}

type experienceService struct {
	repo repository.ExperienceRepository
}

func NewExperienceService(repo repository.ExperienceRepository) ExperienceService {
	return &experienceService{repo}
}

func (s *experienceService) GetAll(ctx context.Context) ([]model.Experience, error) {
	return s.repo.GetAll(ctx)
}

func (s *experienceService) GetByID(ctx context.Context, id int) (*model.Experience, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *experienceService) Create(ctx context.Context, e *model.Experience) error {
	if err := validation.ValidateExperience(e); err != nil {
		log.Printf("[ExperienceService] Create validation error: %v", err)
		return err
	}
	err := s.repo.Create(ctx, e)
	if err != nil {
		log.Printf("[ExperienceService] Create DB error: %v", err)
		return err
	}
	log.Printf("[ExperienceService] Created experience: %+v", e)
	return nil
}

func (s *experienceService) Update(ctx context.Context, e *model.Experience) error {
	if e.ID == 0 {
		log.Printf("[ExperienceService] Update error: id is required")
		return errors.New("id is required")
	}
	if err := validation.ValidateExperience(e); err != nil {
		log.Printf("[ExperienceService] Update validation error: %v", err)
		return err
	}
	err := s.repo.Update(ctx, e)
	if err != nil {
		log.Printf("[ExperienceService] Update DB error: %v", err)
		return err
	}
	log.Printf("[ExperienceService] Updated experience: %+v", e)
	return nil
}

func (s *experienceService) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		log.Printf("[ExperienceService] Delete DB error: %v", err)
		return err
	}
	log.Printf("[ExperienceService] Deleted experience id: %d", id)
	return nil
}
