package model

import (
	"github.com/brianvoe/gofakeit"
	"github.com/oklog/ulid/v2"
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

func GeneratePerson() (*Person, error) {
	person := Person{
		PersonCommon: PersonCommon{
			Id:           ulid.Make().String(),
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
		},
	}
	return &person, nil
}

func GenerateFoodBank() (*FoodBank, error) {
	foodBank := FoodBank{
		Id:   ulid.Make().String(),
		Name: gofakeit.Company(),
		Address: Address{
			Street1: gofakeit.Street(),
			City:    gofakeit.City(),
			State:   gofakeit.State(),
			Zip:     gofakeit.Zip(),
			Country: gofakeit.Country(),
		},
	}
	return &foodBank, nil
}

func GenerateFoodBankVisit() (*FoodBankVisit, error) {
	foodBankVisit := FoodBankVisit{
		Id:         ulid.Make().String(),
		Date:       gofakeit.Date().Format("2006-01-02"),
		PersonId:   ulid.Make().String(),
		FoodBankId: ulid.Make().String(),
		Notes:      gofakeit.Sentence(5),
	}
	return &foodBankVisit, nil
}

func GenerateItem() (*Item, error) {
	item := Item{
		Id:         ulid.Make().String(),
		FoodBankId: ulid.Make().String(),
		Name:       gofakeit.Word(),
		Points:     gofakeit.Number(1, 100),
	}
	return &item, nil
}

func GeneratePeople(n int) ([]Person, error) {
	people := make([]Person, n)
	for i := range people {
		people[i] = Person{
			PersonCommon: PersonCommon{
				Id:           ulid.Make().String(),
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
			},
		}
	}
	return people, nil
}

func GenerateFoodBanks(n int) ([]FoodBank, error) {
	foodBanks := make([]FoodBank, n)
	for i := range foodBanks {
		foodBank, err := GenerateFoodBank()
		if err != nil {
			return nil, err
		}
		foodBanks[i] = *foodBank
	}
	return foodBanks, nil
}

func GenerateFoodBankVisits(n int) ([]FoodBankVisit, error) {
	foodBankVisits := make([]FoodBankVisit, n)
	for i := range foodBankVisits {
		foodBankVisit, err := GenerateFoodBankVisit()
		if err != nil {
			return nil, err
		}
		foodBankVisits[i] = *foodBankVisit
	}
	return foodBankVisits, nil
}

func GenerateItems(n int) ([]Item, error) {
	items := make([]Item, n)
	for i := range items {
		item, err := GenerateItem()
		if err != nil {
			return nil, err
		}
		items[i] = *item
	}
	return items, nil
}
