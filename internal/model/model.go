package model

import (
	"regexp"
)

type Entity interface {
	GetID() string
}

type ValidationError struct {
	Field   string
	Type    string
	Message string
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

func (p Person) Validate() ValidationErrors {
	var errors ValidationErrors

	if p.FirstName == "" {
		errors = append(errors, ValidationError{Field: "firstName", Type: "missing", Message: "field_missing"})
	}
	if p.LastName == "" {
		errors = append(errors, ValidationError{Field: "lastName", Type: "missing", Message: "field_missing"})
	}
	if p.Email == "" {
		errors = append(errors, ValidationError{Field: "email", Type: "missing", Message: "field_missing"})
	} else if !isValidEmail(p.Email) {
		errors = append(errors, ValidationError{Field: "email", Type: "invalid", Message: "invalid_email"})
	}
	// Add more validation rules as needed

	return errors
}

func (fb FoodBank) Validate() ValidationErrors {
	var errors ValidationErrors

	if fb.Name == "" {
		errors = append(errors, ValidationError{Field: "name", Type: "missing", Message: "field_missing"})
	}
	if fb.Address.Street1 == "" {
		errors = append(errors, ValidationError{Field: "address.street1", Type: "missing", Message: "field_missing"})
	}
	if fb.Address.City == "" {
		errors = append(errors, ValidationError{Field: "address.city", Type: "missing", Message: "field_missing"})
	}
	if fb.Address.State == "" {
		errors = append(errors, ValidationError{Field: "address.state", Type: "missing", Message: "field_missing"})
	}
	if fb.Address.Zip == "" {
		errors = append(errors, ValidationError{Field: "address.zip", Type: "missing", Message: "field_missing"})
	}
	if fb.Address.Country == "" {
		errors = append(errors, ValidationError{Field: "address.country", Type: "missing", Message: "field_missing"})
	}
	// Add more validation rules as needed

	return errors
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

type Person struct {
	Id           string `json:"id"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	Street       string `json:"street"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postalCode"`
	Phone        string `json:"phone"`
	Gender       string `json:"gender"`
	DOB          string `json:"dob"`
	Race         string `json:"race"`
	Language     string `json:"language"`
	Relationship string `json:"relationship"`
}

func (p Person) GetID() string {
	return p.Id
}

type FoodBank struct {
	Id      string  `json:"id"`
	Name    string  `json:"name"`
	Address Address `json:"address"`
}

func (fb FoodBank) GetID() string {
	return fb.Id
}

type Address struct {
	Street1 string `json:"street1"`
	Street2 string `json:"street2"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Country string `json:"country"`
}

type FoodBankVisit struct {
	Id         string `json:"id"`
	Date       string `json:"date"`
	PersonId   string `json:"personId"`
	FoodBankId string `json:"foodBankId"`
	Notes      string `json:"notes"`
}

func (fbv FoodBankVisit) GetID() string {
	return fbv.Id
}

type Item struct {
	Id         string `json:"id"`
	FoodBankId string `json:"foodBankId"`
	Name       string `json:"name"`
	Points     int    `json:"points"`
}

func (i Item) GetID() string {
	return i.Id
}
