package routes

import (
	"cupboard/internal/db"
)

func NewRoutes(dbInstance *db.FirestoreDB, emailSender email.EmailSender) (*PersonsHandler, *FoodBanksHandler, *ItemsHandler, *VisitsHandler) {
	return routes.NewPersonsHandler(dbInstance, emailSender), &FoodBanksHandler{DB: dbInstance}, &ItemsHandler{DB: dbInstance}, &VisitsHandler{DB: dbInstance}
}
