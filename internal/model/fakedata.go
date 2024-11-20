package model

import (
	"github.com/brianvoe/gofakeit"
)

func randomStringFromSlice(slice []string) string {
	return slice[gofakeit.Number(0, len(slice)-1)]
}

func Race() string {
	races := []string{
		"White",
		"Black or African American",
		"American Indian or Alaska Native",
		"Asian",
		"Native Hawaiian or Other Pacific Islander",
		"Some other race",
	}
	return randomStringFromSlice(races)
}

func Language() string {
	languages := []string{
		"English",
		"Spanish",
		"Chinese",
		"Tagalog",
		"French",
		"Vietnamese",
		"German",
		"Korean",
		"Russian",
		"Arabic",
		"Other languages",
	}
	return randomStringFromSlice(languages)
}

func Relationship() string {
	relationships := []string{
		"Spouse",
		"Child",
		"Sibling",
		"Parent",
		"Grandparent",
		"Grandchild",
		"Other relative",
		"Non-relative",
	}
	return randomStringFromSlice(relationships)
}

func GenerateHouseholds(n int) ([]Household, error) {
	households := make([]Household, n)
	for i := range households {
		person, err := GeneratePerson()
		if err != nil {
			return nil, err
		}
		members, err := GeneratePeople(gofakeit.Number(1, 5)) // generate 1-5 members
		if err != nil {
			return nil, err
		}
		households[i] = Household{
			Id:      gofakeit.UUID(),
			Head:    *person,
			Members: members,
		}
	}
	return households, nil
}

func GeneratePeople(n int) ([]Person, error) {
	people := make([]Person, n)
	for i := range people {
		people[i] = Person{
			FirstName:    gofakeit.FirstName(),
			LastName:     gofakeit.LastName(),
			Email:        gofakeit.Email(),
			Street:       gofakeit.Street(),
			City:         gofakeit.City(),
			State:        gofakeit.State(),
			PostalCode:   gofakeit.Zip(),
			Phone:        gofakeit.Phone(),
			Gender:       gofakeit.Gender(),
			DOB:          gofakeit.Date().Format("2006-01-02"),
			Race:         Race(),
			Language:     Language(),
			Relationship: Relationship(),
		}
	}
	return people, nil
}

func GeneratePerson() (*Person, error) {
	person := Person{
		FirstName:    gofakeit.FirstName(),
		LastName:     gofakeit.LastName(),
		Email:        gofakeit.Email(),
		Street:       gofakeit.Street(),
		City:         gofakeit.City(),
		State:        gofakeit.State(),
		PostalCode:   gofakeit.Zip(),
		Phone:        gofakeit.Phone(),
		Gender:       gofakeit.Gender(),
		DOB:          gofakeit.Date().Format("2006-01-02"),
		Race:         Race(),
		Language:     Language(),
		Relationship: Relationship(),
	}
	return &person, nil
}
