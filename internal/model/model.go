package model

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

type FoodBank struct {
	Id      string  `json:"id"`
	Name    string  `json:"name"`
	Address Address `json:"address"`
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

type Item struct {
	Id         string `json:"id"`
	FoodBankId string `json:"foodBankId"`
	Name       string `json:"name"`
	Points     int    `json:"points"`
}
