package db

import (
	"context"
	"fmt"
	"testing"

	"cupboard/internal/model"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestFirestoreDB tests the FirestoreDB struct methods.
func TestFirestoreDB(t *testing.T) {
	// Initialize Firestore client
	ctx := context.Background()
	projectID := "test-project"
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	db := NewFirestoreDB(client)

	// Generate random household data
	household, err := model.GenerateHouseholds(1)
	if err != nil {
		t.Fatalf("Failed to generate household data: %v", err)
	}

	// Test AddHousehold
	err = db.AddHousehold(ctx, household[0])
	assert.NoError(t, err, "Failed to add household")

	// Test GetHouseholdByID
	retrievedHousehold, err := db.GetHouseholdByID(ctx, household[0].Id)
	assert.NoError(t, err, "Failed to retrieve household by ID")
	assert.Equal(t, household[0].Id, retrievedHousehold.Id, "Household ID mismatch")
	assert.Equal(t, household[0].Head.FirstName, retrievedHousehold.Head.FirstName, "Household head first name mismatch")
	assert.Equal(t, household[0].Head.LastName, retrievedHousehold.Head.LastName, "Household head last name mismatch")

	// Test GetHouseholds
	households, err := db.GetHouseholds(ctx)
	assert.NoError(t, err, "Failed to retrieve households")
	assert.GreaterOrEqual(t, len(households), 1, "No households retrieved")

	// Test DeleteHousehold
	err = db.DeleteHousehold(ctx, household[0].Id)
	assert.NoError(t, err, "Failed to delete household")

	// Verify deletion
	_, err = db.GetHouseholdByID(ctx, household[0].Id)
	assert.Error(t, err, "Household should not exist after deletion")
	assert.Equal(t, codes.NotFound, status.Code(err), "Expected not found error after deletion")
}

// TestFirestoreDBConcurrency tests concurrent operations on FirestoreDB.
func TestFirestoreDBConcurrency(t *testing.T) {
	// Initialize Firestore client
	ctx := context.Background()
	projectID := "test-project"
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	db := NewFirestoreDB(client)

	// Generate random household data
	households, err := model.GenerateHouseholds(10)
	if err != nil {
		t.Fatalf("Failed to generate household data: %v", err)
	}

	// Test concurrent AddHousehold
	errChan := make(chan error, len(households))
	for _, household := range households {
		go func(h model.Household) {
			errChan <- db.AddHousehold(ctx, h)
		}(household)
	}

	for i := 0; i < len(households); i++ {
		err := <-errChan
		assert.NoError(t, err, fmt.Sprintf("Failed to add household %d", i))
	}

	// Test concurrent GetHouseholdByID
	for _, household := range households {
		go func(h model.Household) {
			retrievedHousehold, err := db.GetHouseholdByID(ctx, h.Id)
			errChan <- err
			if err == nil {
				assert.Equal(t, h.Id, retrievedHousehold.Id, "Household ID mismatch")
				assert.Equal(t, h.Head.FirstName, retrievedHousehold.Head.FirstName, "Household head first name mismatch")
				assert.Equal(t, h.Head.LastName, retrievedHousehold.Head.LastName, "Household head last name mismatch")
			}
		}(household)
	}

	for i := 0; i < len(households); i++ {
		err := <-errChan
		assert.NoError(t, err, fmt.Sprintf("Failed to retrieve household %d by ID", i))
	}

	// Test concurrent DeleteHousehold
	for _, household := range households {
		go func(h model.Household) {
			errChan <- db.DeleteHousehold(ctx, h.Id)
		}(household)
	}

	for i := 0; i < len(households); i++ {
		err := <-errChan
		assert.NoError(t, err, fmt.Sprintf("Failed to delete household %d", i))
	}

	// Verify deletion
	for _, household := range households {
		go func(h model.Household) {
			_, err := db.GetHouseholdByID(ctx, h.Id)
			errChan <- err
		}(household)
	}

	for i := 0; i < len(households); i++ {
		err := <-errChan
		assert.Error(t, err, fmt.Sprintf("Household %d should not exist after deletion", i))
		assert.Equal(t, codes.NotFound, status.Code(err), fmt.Sprintf("Expected not found error after deletion for household %d", i))
	}
}
package db

import (
	"context"
	"fmt"
	"os"
	"testing"

	"cloud.google.com/go/firestore"
	"cupboard/internal/model"
	"github.com/oklog/ulid/v2"
	"google.golang.org/api/iterator"
)

func TestMain(m *testing.M) {
	os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8080")
	code := m.Run()
	os.Exit(code)
}

func TestFirestoreDB_PutAndGetPerson(t *testing.T) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "projectID")
	if err != nil {
		t.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	dbInstance := NewFirestoreDB(client)

	person, err := model.GeneratePerson()
	if err != nil {
		t.Fatalf("Failed to generate person: %v", err)
	}

	err = dbInstance.PutPerson(ctx, *person)
	if err != nil {
		t.Fatalf("Failed to put person: %v", err)
	}

	retrievedPerson, err := dbInstance.GetPerson(ctx, person.Id)
	if err != nil {
		t.Fatalf("Failed to get person: %v", err)
	}

	if retrievedPerson.Id != person.Id {
		t.Errorf("Expected person ID %s, got %s", person.Id, retrievedPerson.Id)
	}
}

func TestFirestoreDB_PutAndGetFoodBank(t *testing.T) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "projectID")
	if err != nil {
		t.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	dbInstance := NewFirestoreDB(client)

	foodBank, err := model.GenerateFoodBank()
	if err != nil {
		t.Fatalf("Failed to generate food bank: %v", err)
	}

	err = dbInstance.PutFoodBank(ctx, *foodBank)
	if err != nil {
		t.Fatalf("Failed to put food bank: %v", err)
	}

	retrievedFoodBank, err := dbInstance.GetFoodBank(ctx, foodBank.Id)
	if err != nil {
		t.Fatalf("Failed to get food bank: %v", err)
	}

	if retrievedFoodBank.Id != foodBank.Id {
		t.Errorf("Expected food bank ID %s, got %s", foodBank.Id, retrievedFoodBank.Id)
	}
}

func TestFirestoreDB_PutAndGetFoodBankVisit(t *testing.T) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "projectID")
	if err != nil {
		t.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	dbInstance := NewFirestoreDB(client)

	visit, err := model.GenerateFoodBankVisit()
	if err != nil {
		t.Fatalf("Failed to generate food bank visit: %v", err)
	}

	err = dbInstance.PutFoodBankVisit(ctx, *visit)
	if err != nil {
		t.Fatalf("Failed to put food bank visit: %v", err)
	}

	retrievedVisit, err := dbInstance.GetFoodBankVisit(ctx, visit.Id)
	if err != nil {
		t.Fatalf("Failed to get food bank visit: %v", err)
	}

	if retrievedVisit.Id != visit.Id {
		t.Errorf("Expected food bank visit ID %s, got %s", visit.Id, retrievedVisit.Id)
	}
}

func TestFirestoreDB_PutAndGetItem(t *testing.T) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "projectID")
	if err != nil {
		t.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	dbInstance := NewFirestoreDB(client)

	item, err := model.GenerateItem()
	if err != nil {
		t.Fatalf("Failed to generate item: %v", err)
	}

	err = dbInstance.PutItem(ctx, *item)
	if err != nil {
		t.Fatalf("Failed to put item: %v", err)
	}

	retrievedItem, err := dbInstance.GetItem(ctx, item.Id)
	if err != nil {
		t.Fatalf("Failed to get item: %v", err)
	}

	if retrievedItem.Id != item.Id {
		t.Errorf("Expected item ID %s, got %s", item.Id, retrievedItem.Id)
	}
}

func TestFirestoreDB_PutPersonsAndGetPersons(t *testing.T) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "projectID")
	if err != nil {
		t.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	dbInstance := NewFirestoreDB(client)

	persons, err := model.GeneratePeople(5)
	if err != nil {
		t.Fatalf("Failed to generate persons: %v", err)
	}

	err = dbInstance.PutPersons(ctx, persons)
	if err != nil {
		t.Fatalf("Failed to put persons: %v", err)
	}

	for _, person := range persons {
		retrievedPerson, err := dbInstance.GetPerson(ctx, person.Id)
		if err != nil {
			t.Fatalf("Failed to get person: %v", err)
		}

		if retrievedPerson.Id != person.Id {
			t.Errorf("Expected person ID %s, got %s", person.Id, retrievedPerson.Id)
		}
	}
}

func TestFirestoreDB_PutFoodBanksAndGetFoodBanks(t *testing.T) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "projectID")
	if err != nil {
		t.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	dbInstance := NewFirestoreDB(client)

	foodBanks, err := model.GenerateFoodBanks(5)
	if err != nil {
		t.Fatalf("Failed to generate food banks: %v", err)
	}

	err = dbInstance.PutFoodBanks(ctx, foodBanks)
	if err != nil {
		t.Fatalf("Failed to put food banks: %v", err)
	}

	for _, foodBank := range foodBanks {
		retrievedFoodBank, err := dbInstance.GetFoodBank(ctx, foodBank.Id)
		if err != nil {
			t.Fatalf("Failed to get food bank: %v", err)
		}

		if retrievedFoodBank.Id != foodBank.Id {
			t.Errorf("Expected food bank ID %s, got %s", foodBank.Id, retrievedFoodBank.Id)
		}
	}
}

func TestFirestoreDB_PutFoodBankVisitsAndGetFoodBankVisits(t *testing.T) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "projectID")
	if err != nil {
		t.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	dbInstance := NewFirestoreDB(client)

	visits, err := model.GenerateFoodBankVisits(5)
	if err != nil {
		t.Fatalf("Failed to generate food bank visits: %v", err)
	}

	err = dbInstance.PutFoodBankVisits(ctx, visits)
	if err != nil {
		t.Fatalf("Failed to put food bank visits: %v", err)
	}

	for _, visit := range visits {
		retrievedVisit, err := dbInstance.GetFoodBankVisit(ctx, visit.Id)
		if err != nil {
			t.Fatalf("Failed to get food bank visit: %v", err)
		}

		if retrievedVisit.Id != visit.Id {
			t.Errorf("Expected food bank visit ID %s, got %s", visit.Id, retrievedVisit.Id)
		}
	}
}

func TestFirestoreDB_PutItemsAndGetItems(t *testing.T) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "projectID")
	if err != nil {
		t.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	dbInstance := NewFirestoreDB(client)

	items, err := model.GenerateItems(5)
	if err != nil {
		t.Fatalf("Failed to generate items: %v", err)
	}

	err = dbInstance.PutItems(ctx, items)
	if err != nil {
		t.Fatalf("Failed to put items: %v", err)
	}

	for _, item := range items {
		retrievedItem, err := dbInstance.GetItem(ctx, item.Id)
		if err != nil {
			t.Fatalf("Failed to get item: %v", err)
		}

		if retrievedItem.Id != item.Id {
			t.Errorf("Expected item ID %s, got %s", item.Id, retrievedItem.Id)
		}
	}
}

func TestFirestoreDB_DeletePerson(t *testing.T) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "projectID")
	if err != nil {
		t.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	dbInstance := NewFirestoreDB(client)

	person, err := model.GeneratePerson()
	if err != nil {
		t.Fatalf("Failed to generate person: %v", err)
	}

	err = dbInstance.PutPerson(ctx, *person)
	if err != nil {
		t.Fatalf("Failed to put person: %v", err)
	}

	err = dbInstance.DeletePerson(ctx, person.Id)
	if err != nil {
		t.Fatalf("Failed to delete person: %v", err)
	}

	_, err = dbInstance.GetPerson(ctx, person.Id)
	if err == nil {
		t.Errorf("Expected error when retrieving deleted person, got nil")
	}
}

func TestFirestoreDB_DeleteFoodBank(t *testing.T) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "projectID")
	if err != nil {
		t.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	dbInstance := NewFirestoreDB(client)

	foodBank, err := model.GenerateFoodBank()
	if err != nil {
		t.Fatalf("Failed to generate food bank: %v", err)
	}

	err = dbInstance.PutFoodBank(ctx, *foodBank)
	if err != nil {
		t.Fatalf("Failed to put food bank: %v", err)
	}

	err = dbInstance.DeleteFoodBank(ctx, foodBank.Id)
	if err != nil {
		t.Fatalf("Failed to delete food bank: %v", err)
	}

	_, err = dbInstance.GetFoodBank(ctx, foodBank.Id)
	if err == nil {
		t.Errorf("Expected error when retrieving deleted food bank, got nil")
	}
}

func TestFirestoreDB_DeleteFoodBankVisit(t *testing.T) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "projectID")
	if err != nil {
		t.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	dbInstance := NewFirestoreDB(client)

	visit, err := model.GenerateFoodBankVisit()
	if err != nil {
		t.Fatalf("Failed to generate food bank visit: %v", err)
	}

	err = dbInstance.PutFoodBankVisit(ctx, *visit)
	if err != nil {
		t.Fatalf("Failed to put food bank visit: %v", err)
	}

	err = dbInstance.DeleteFoodBankVisit(ctx, visit.Id)
	if err != nil {
		t.Fatalf("Failed to delete food bank visit: %v", err)
	}

	_, err = dbInstance.GetFoodBankVisit(ctx, visit.Id)
	if err == nil {
		t.Errorf("Expected error when retrieving deleted food bank visit, got nil")
	}
}

func TestFirestoreDB_DeleteItem(t *testing.T) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "projectID")
	if err != nil {
		t.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	dbInstance := NewFirestoreDB(client)

	item, err := model.GenerateItem()
	if err != nil {
		t.Fatalf("Failed to generate item: %v", err)
	}

	err = dbInstance.PutItem(ctx, *item)
	if err != nil {
		t.Fatalf("Failed to put item: %v", err)
	}

	err = dbInstance.DeleteItem(ctx, item.Id)
	if err != nil {
		t.Fatalf("Failed to delete item: %v", err)
	}

	_, err = dbInstance.GetItem(ctx, item.Id)
	if err == nil {
		t.Errorf("Expected error when retrieving deleted item, got nil")
	}
}

func TestFirestoreDB_DeletePersons(t *testing.T) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "projectID")
	if err != nil {
		t.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	dbInstance := NewFirestoreDB(client)

	persons, err := model.GeneratePeople(5)
	if err != nil {
		t.Fatalf("Failed to generate persons: %v", err)
	}

	err = dbInstance.PutPersons(ctx, persons)
	if err != nil {
		t.Fatalf("Failed to put persons: %v", err)
	}

	ids := make([]string, len(persons))
	for i, person := range persons {
		ids[i] = person.Id
	}

	err = dbInstance.DeletePersons(ctx, ids)
	if err != nil {
		t.Fatalf("Failed to delete persons: %v", err)
	}

	for _, id := range ids {
		_, err := dbInstance.GetPerson(ctx, id)
		if err == nil {
			t.Errorf("Expected error when retrieving deleted person with ID %s, got nil", id)
		}
	}
}

func TestFirestoreDB_DeleteFoodBanks(t *testing.T) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "projectID")
	if err != nil {
		t.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	dbInstance := NewFirestoreDB(client)

	foodBanks, err := model.GenerateFoodBanks(5)
	if err != nil {
		t.Fatalf("Failed to generate food banks: %v", err)
	}

	err = dbInstance.PutFoodBanks(ctx, foodBanks)
	if err != nil {
		t.Fatalf("Failed to put food banks: %v", err)
	}

	ids := make([]string, len(foodBanks))
	for i, foodBank := range foodBanks {
		ids[i] = foodBank.Id
	}

	err = dbInstance.DeleteFoodBanks(ctx, ids)
	if err != nil {
		t.Fatalf("Failed to delete food banks: %v", err)
	}

	for _, id := range ids {
		_, err := dbInstance.GetFoodBank(ctx, id)
		if err == nil {
			t.Errorf("Expected error when retrieving deleted food bank with ID %s, got nil", id)
		}
	}
}

func TestFirestoreDB_DeleteFoodBankVisits(t *testing.T) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "projectID")
	if err != nil {
		t.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	dbInstance := NewFirestoreDB(client)

	visits, err := model.GenerateFoodBankVisits(5)
	if err != nil {
		t.Fatalf("Failed to generate food bank visits: %v", err)
	}

	err = dbInstance.PutFoodBankVisits(ctx, visits)
	if err != nil {
		t.Fatalf("Failed to put food bank visits: %v", err)
	}

	ids := make([]string, len(visits))
	for i, visit := range visits {
		ids[i] = visit.Id
	}

	err = dbInstance.DeleteFoodBankVisits(ctx, ids)
	if err != nil {
		t.Fatalf("Failed to delete food bank visits: %v", err)
	}

	for _, id := range ids {
		_, err := dbInstance.GetFoodBankVisit(ctx, id)
		if err == nil {
			t.Errorf("Expected error when retrieving deleted food bank visit with ID %s, got nil", id)
		}
	}
}

func TestFirestoreDB_DeleteItems(t *testing.T) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "projectID")
	if err != nil {
		t.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	dbInstance := NewFirestoreDB(client)

	items, err := model.GenerateItems(5)
	if err != nil {
		t.Fatalf("Failed to generate items: %v", err)
	}

	err = dbInstance.PutItems(ctx, items)
	if err != nil {
		t.Fatalf("Failed to put items: %v", err)
	}

	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.Id
	}

	err = dbInstance.DeleteItems(ctx, ids)
	if err != nil {
		t.Fatalf("Failed to delete items: %v", err)
	}

	for _, id := range ids {
		_, err := dbInstance.GetItem(ctx, id)
		if err == nil {
			t.Errorf("Expected error when retrieving deleted item with ID %s, got nil", id)
		}
	}
}
