package validation

import (
	"errors"
	"porto/model"
	"regexp"
	"strings"
)

func ValidatePortfolio(p *model.Portfolio) error {
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("portfolio name is required")
	}
	if strings.TrimSpace(p.Description) == "" {
		return errors.New("portfolio description is required")
	}
	return nil
}

func ValidateExperience(e *model.Experience) error {
	if strings.TrimSpace(e.Title) == "" {
		return errors.New("experience title is required")
	}
	if strings.TrimSpace(e.Company) == "" {
		return errors.New("experience company is required")
	}
	if strings.TrimSpace(e.StartDate) == "" {
		return errors.New("start date is required")
	}
	return nil
}

func ValidateContact(c *model.Contact) error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("contact name is required")
	}
	if strings.TrimSpace(c.Email) == "" {
		return errors.New("contact email is required")
	}
	if !isValidEmail(c.Email) {
		return errors.New("invalid email format")
	}
	if strings.TrimSpace(c.Message) == "" {
		return errors.New("contact message is required")
	}
	return nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
