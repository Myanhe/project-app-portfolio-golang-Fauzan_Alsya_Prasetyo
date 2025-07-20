package validation

import (
	"porto/model"
	"testing"
)

func TestValidatePortfolio(t *testing.T) {
	cases := []struct {
		name      string
		portfolio model.Portfolio
		wantErr   bool
	}{
		{"valid", model.Portfolio{Name: "A", Description: "B"}, false},
		{"empty name", model.Portfolio{Name: "", Description: "B"}, true},
		{"empty desc", model.Portfolio{Name: "A", Description: ""}, true},
	}
	for _, c := range cases {
		err := ValidatePortfolio(&c.portfolio)
		if (err != nil) != c.wantErr {
			t.Errorf("%s: got error %v, wantErr %v", c.name, err, c.wantErr)
		}
	}
}

func TestValidateContact(t *testing.T) {
	cases := []struct {
		name    string
		contact model.Contact
		wantErr bool
	}{
		{"valid", model.Contact{Name: "A", Email: "a@mail.com", Message: "hi"}, false},
		{"empty name", model.Contact{Name: "", Email: "a@mail.com", Message: "hi"}, true},
		{"invalid email", model.Contact{Name: "A", Email: "a", Message: "hi"}, true},
		{"empty message", model.Contact{Name: "A", Email: "a@mail.com", Message: ""}, true},
	}
	for _, c := range cases {
		err := ValidateContact(&c.contact)
		if (err != nil) != c.wantErr {
			t.Errorf("%s: got error %v, wantErr %v", c.name, err, c.wantErr)
		}
	}
}
