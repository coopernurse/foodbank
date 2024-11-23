package routes

import (
	"cupboard/internal/db"
	"cupboard/internal/email"
)

func NewRoutes(dbInstance *db.FirestoreDB, emailSender email.EmailSender) (*PersonsHandler,
	*FoodBanksHandler, *ItemsHandler, *VisitsHandler, *AuthHandler) {
	return NewPersonsHandler(dbInstance, emailSender),
		&FoodBanksHandler{DB: dbInstance},
		&ItemsHandler{DB: dbInstance},
		NewVisitsHandler(dbInstance, emailSender), NewAuthHandler(dbInstance, emailSender)
}
