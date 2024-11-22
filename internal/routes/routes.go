package routes

import (
	"cupboard/internal/db"
)

func NewRoutes(dbInstance *db.FirestoreDB, emailSender email.EmailSender) (*PersonsHandler, *FoodBanksHandler, *ItemsHandler, *VisitsHandler) {
	return NewPersonsHandler(dbInstance, emailSender), &FoodBanksHandler{DB: dbInstance}, &ItemsHandler{DB: dbInstance}, NewVisitsHandler(dbInstance, emailSender)
}
