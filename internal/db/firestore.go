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
func (db *FirestoreDB) GetHouseholds(ctx context.Context, pageSize int, startAfter string) ([]model.Household, string, error) {
	var households []model.Household
	var query *firestore.Query

	if startAfter == "" {
		query = db.Client.Collection("households").OrderBy("Id", firestore.Desc).Limit(pageSize)
	} else {
		lastDoc, err := db.Client.Collection("households").Doc(startAfter).Get(ctx)
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving last document for pagination: %w", err)
		}
		query = db.Client.Collection("households").OrderBy("Id", firestore.Desc).Limit(pageSize).StartAfter(lastDoc.Data())
	}

	iter := query.Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, "", fmt.Errorf("error retrieving households: %w", err)
		}

		var household model.Household
		if err := doc.DataTo(&household); err != nil {
			return nil, "", fmt.Errorf("error parsing household data: %w", err)
		}
		households = append(households, household)
	}

	nextPageToken := ""
	if iter.Next() != iterator.Done {
		nextPageToken = households[len(households)-1].Id
	}

	return households, nextPageToken, nil
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

func (db *FirestoreDB) PutPerson(ctx context.Context, person model.Person) error {
	_, err := db.Client.Collection("persons").Doc(person.Id).Set(ctx, person)
	if err != nil {
		return fmt.Errorf("error saving person: %w", err)
	}
	return nil
}

func (db *FirestoreDB) PutPersons(ctx context.Context, persons []model.Person) error {
	batch := db.Client.Batch()
	for _, person := range persons {
		if person.Id == "" {
			person.Id = ulid.Make().String()
		}
		batch.Set(db.Client.Collection("persons").Doc(person.Id), person)
	}
	_, err := batch.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error saving persons: %w", err)
	}
	return nil
}

func (db *FirestoreDB) GetPerson(ctx context.Context, id string) (*model.Person, error) {
	doc, err := db.Client.Collection("persons").Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving person with ID %s: %w", id, err)
	}

	var person model.Person
	if err := doc.DataTo(&person); err != nil {
		return nil, fmt.Errorf("error parsing person data for ID %s: %w", id, err)
	}

	return &person, nil
}

func (db *FirestoreDB) DeletePerson(ctx context.Context, id string) error {
	_, err := db.Client.Collection("persons").Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("error deleting person with ID %s: %w", id, err)
	}
	return nil
}

func (db *FirestoreDB) DeletePersons(ctx context.Context, ids []string) error {
	batch := db.Client.Batch()
	for _, id := range ids {
		batch.Delete(db.Client.Collection("persons").Doc(id))
	}
	_, err := batch.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error deleting persons: %w", err)
	}
	return nil
}

func (db *FirestoreDB) PutFoodBank(ctx context.Context, foodBank model.FoodBank) error {
	_, err := db.Client.Collection("foodbanks").Doc(foodBank.Id).Set(ctx, foodBank)
	if err != nil {
		return fmt.Errorf("error saving food bank: %w", err)
	}
	return nil
}

func (db *FirestoreDB) PutFoodBanks(ctx context.Context, foodBanks []model.FoodBank) error {
	batch := db.Client.Batch()
	for _, foodBank := range foodBanks {
		if foodBank.Id == "" {
			foodBank.Id = ulid.Make().String()
		}
		batch.Set(db.Client.Collection("foodbanks").Doc(foodBank.Id), foodBank)
	}
	_, err := batch.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error saving food banks: %w", err)
	}
	return nil
}

func (db *FirestoreDB) GetFoodBank(ctx context.Context, id string) (*model.FoodBank, error) {
	doc, err := db.Client.Collection("foodbanks").Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving food bank with ID %s: %w", id, err)
	}

	var foodBank model.FoodBank
	if err := doc.DataTo(&foodBank); err != nil {
		return nil, fmt.Errorf("error parsing food bank data for ID %s: %w", id, err)
	}

	return &foodBank, nil
}

func (db *FirestoreDB) DeleteFoodBank(ctx context.Context, id string) error {
	_, err := db.Client.Collection("foodbanks").Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("error deleting food bank with ID %s: %w", id, err)
	}
	return nil
}

func (db *FirestoreDB) DeleteFoodBanks(ctx context.Context, ids []string) error {
	batch := db.Client.Batch()
	for _, id := range ids {
		batch.Delete(db.Client.Collection("foodbanks").Doc(id))
	}
	_, err := batch.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error deleting food banks: %w", err)
	}
	return nil
}

func (db *FirestoreDB) PutFoodBankVisit(ctx context.Context, visit model.FoodBankVisit) error {
	_, err := db.Client.Collection("foodbankvisits").Doc(visit.Id).Set(ctx, visit)
	if err != nil {
		return fmt.Errorf("error saving food bank visit: %w", err)
	}
	return nil
}

func (db *FirestoreDB) PutFoodBankVisits(ctx context.Context, visits []model.FoodBankVisit) error {
	batch := db.Client.Batch()
	for _, visit := range visits {
		if visit.Id == "" {
			visit.Id = ulid.Make().String()
		}
		batch.Set(db.Client.Collection("foodbankvisits").Doc(visit.Id), visit)
	}
	_, err := batch.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error saving food bank visits: %w", err)
	}
	return nil
}

func (db *FirestoreDB) GetFoodBankVisit(ctx context.Context, id string) (*model.FoodBankVisit, error) {
	doc, err := db.Client.Collection("foodbankvisits").Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving food bank visit with ID %s: %w", id, err)
	}

	var visit model.FoodBankVisit
	if err := doc.DataTo(&visit); err != nil {
		return nil, fmt.Errorf("error parsing food bank visit data for ID %s: %w", id, err)
	}

	return &visit, nil
}

func (db *FirestoreDB) DeleteFoodBankVisit(ctx context.Context, id string) error {
	_, err := db.Client.Collection("foodbankvisits").Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("error deleting food bank visit with ID %s: %w", id, err)
	}
	return nil
}

func (db *FirestoreDB) DeleteFoodBankVisits(ctx context.Context, ids []string) error {
	batch := db.Client.Batch()
	for _, id := range ids {
		batch.Delete(db.Client.Collection("foodbankvisits").Doc(id))
	}
	_, err := batch.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error deleting food bank visits: %w", err)
	}
	return nil
}

func (db *FirestoreDB) PutItem(ctx context.Context, item model.Item) error {
	_, err := db.Client.Collection("items").Doc(item.Id).Set(ctx, item)
	if err != nil {
		return fmt.Errorf("error saving item: %w", err)
	}
	return nil
}

func (db *FirestoreDB) PutItems(ctx context.Context, items []model.Item) error {
	batch := db.Client.Batch()
	for _, item := range items {
		if item.Id == "" {
			item.Id = ulid.Make().String()
		}
		batch.Set(db.Client.Collection("items").Doc(item.Id), item)
	}
	_, err := batch.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error saving items: %w", err)
	}
	return nil
}

func (db *FirestoreDB) GetItem(ctx context.Context, id string) (*model.Item, error) {
	doc, err := db.Client.Collection("items").Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving item with ID %s: %w", id, err)
	}

	var item model.Item
	if err := doc.DataTo(&item); err != nil {
		return nil, fmt.Errorf("error parsing item data for ID %s: %w", id, err)
	}

	return &item, nil
}

func (db *FirestoreDB) DeleteItem(ctx context.Context, id string) error {
	_, err := db.Client.Collection("items").Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("error deleting item with ID %s: %w", id, err)
	}
	return nil
}

func (db *FirestoreDB) DeleteItems(ctx context.Context, ids []string) error {
	batch := db.Client.Batch()
	for _, id := range ids {
		batch.Delete(db.Client.Collection("items").Doc(id))
	}
	_, err := batch.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error deleting items: %w", err)
	}
	return nil
}
