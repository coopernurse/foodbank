package db

import (
	"context"
	"fmt"

	"cupboard/internal/model"

	"cloud.google.com/go/firestore"
	"github.com/oklog/ulid/v2"
	"google.golang.org/api/iterator"
)

// FirestoreDB encapsulates the Firestore client.
type FirestoreDB struct {
	Client *firestore.Client
}

// NewFirestoreDB creates a new instance of FirestoreDB.
func NewFirestoreDB(client *firestore.Client) *FirestoreDB {
	return &FirestoreDB{Client: client}
}

// GetHouseholds retrieves all households ordered by ID in descending order.
func (db *FirestoreDB) GetHouseholds(ctx context.Context) ([]model.Household, error) {
	var households []model.Household
	iter := db.Client.Collection("households").OrderBy("Id", firestore.Desc).Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error retrieving households: %w", err)
		}

		var household model.Household
		if err := doc.DataTo(&household); err != nil {
			return nil, fmt.Errorf("error parsing household data: %w", err)
		}
		households = append(households, household)
	}

	return households, nil
}

// GetHouseholdByID retrieves a specific household by its ID.
func (db *FirestoreDB) GetHouseholdByID(ctx context.Context, id string) (*model.Household, error) {
	doc, err := db.Client.Collection("households").Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving household with ID %s: %w", id, err)
	}

	var household model.Household
	if err := doc.DataTo(&household); err != nil {
		return nil, fmt.Errorf("error parsing household data for ID %s: %w", id, err)
	}

	return &household, nil
}

// AddHousehold adds a new household to Firestore.
func (db *FirestoreDB) AddHousehold(ctx context.Context, household model.Household) error {
	// Generate a ULID if ID is not set
	if household.Id == "" {
		household.Id = ulid.Make().String()
	}

	_, err := db.Client.Collection("households").Doc(household.Id).Set(ctx, household)
	if err != nil {
		return fmt.Errorf("error saving household: %w", err)
	}
	return nil
}

// DeleteHousehold deletes a specific household by its ID.
func (db *FirestoreDB) DeleteHousehold(ctx context.Context, id string) error {
	_, err := db.Client.Collection("households").Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("error deleting household with ID %s: %w", id, err)
	}
	return nil
}
