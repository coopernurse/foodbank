package routes

import (
	"cupboard/internal/db"
)

func NewRoutes(dbInstance *db.FirestoreDB) (*PersonsHandler, *FoodBanksHandler, *ItemsHandler, *VisitsHandler) {
	return &PersonsHandler{DB: dbInstance}, &FoodBanksHandler{DB: dbInstance}, &ItemsHandler{DB: dbInstance}, &VisitsHandler{DB: dbInstance}
}
