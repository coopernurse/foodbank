package db

import (
	"context"
	"os"
	"testing"

	"cupboard/internal/model"

	"cloud.google.com/go/firestore"
)

func TestMain(m *testing.M) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8080")
	}
	code := m.Run()
	os.Exit(code)
}

func newFirestoreDB(t *testing.T) *FirestoreDB {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "test-project")
	if err != nil {
		t.Fatalf("Failed to create Firestore client: %v", err)
	}
	t.Cleanup(func() { client.Close() })
	return NewFirestoreDB(client)
}

func TestFirestoreDB_PutAndGetPerson(t *testing.T) {
	dbInstance := newFirestoreDB(t)
	testPutAndGet(t, dbInstance, model.GeneratePerson, dbInstance.PutPerson, dbInstance.GetPerson)
}

func TestFirestoreDB_PutAndGetFoodBank(t *testing.T) {
	dbInstance := newFirestoreDB(t)
	testPutAndGet(t, dbInstance, model.GenerateFoodBank, dbInstance.PutFoodBank, dbInstance.GetFoodBank)
}

func TestFirestoreDB_PutAndGetFoodBankVisit(t *testing.T) {
	dbInstance := newFirestoreDB(t)
	testPutAndGet(t, dbInstance, model.GenerateFoodBankVisit, dbInstance.PutFoodBankVisit, dbInstance.GetFoodBankVisit)
}

func TestFirestoreDB_PutAndGetItem(t *testing.T) {
	dbInstance := newFirestoreDB(t)
	testPutAndGet(t, dbInstance, model.GenerateItem, dbInstance.PutItem, dbInstance.GetItem)
}

func testPutAndGet[T model.Entity](t *testing.T, dbInstance *FirestoreDB, generateFunc func() (*T, error), putFunc func(context.Context, T) error, getFunc func(context.Context, string) (*T, error)) {
	ctx := context.Background()

	entity, err := generateFunc()
	if err != nil {
		t.Fatalf("Failed to generate entity: %v", err)
	}

	err = putFunc(ctx, *entity)
	if err != nil {
		t.Fatalf("Failed to put entity: %v", err)
	}

	retrievedEntity, err := getFunc(ctx, entity.GetID())
	if err != nil {
		t.Fatalf("Failed to get entity: %v", err)
	}

	if retrievedEntity.GetID() != entity.GetID() {
		t.Errorf("Expected entity ID %s, got %s", entity.GetID(), retrievedEntity.GetID())
	}
}

func TestFirestoreDB_PutPersonsAndGetPersons(t *testing.T) {
	dbInstance := newFirestoreDB(t)
	testPutAndGets(t, dbInstance, model.GeneratePeople, dbInstance.PutPersons, dbInstance.GetPersons)
}

func TestFirestoreDB_PutFoodBanksAndGetFoodBanks(t *testing.T) {
	dbInstance := newFirestoreDB(t)
	testPutAndGets(t, dbInstance, model.GenerateFoodBanks, dbInstance.PutFoodBanks, dbInstance.GetFoodBank)
}

func TestFirestoreDB_PutFoodBankVisitsAndGetFoodBankVisits(t *testing.T) {
	dbInstance := newFirestoreDB(t)
	testPutAndGets(t, dbInstance, model.GenerateFoodBankVisits, dbInstance.PutFoodBankVisits, dbInstance.GetFoodBankVisit)
}

func TestFirestoreDB_PutItemsAndGetItems(t *testing.T) {
	dbInstance := newFirestoreDB(t)
	testPutAndGets(t, dbInstance, model.GenerateItems, dbInstance.PutItems, dbInstance.GetItem)
}

func testPutAndGets[T model.Entity](t *testing.T, dbInstance *FirestoreDB, generateFunc func(int) ([]T, error), putFunc func(context.Context, []T) error, getFunc func(context.Context, string) (*T, error)) {
	ctx := context.Background()

	entities, err := generateFunc(5)
	if err != nil {
		t.Fatalf("Failed to generate entities: %v", err)
	}

	err = putFunc(ctx, entities)
	if err != nil {
		t.Fatalf("Failed to put entities: %v", err)
	}

	for _, entity := range entities {
		retrievedEntity, err := getFunc(ctx, entity.GetID())
		if err != nil {
			t.Fatalf("Failed to get entity: %v", err)
		}

		if retrievedEntity.GetID() != entity.GetID() {
			t.Errorf("Expected entity ID %s, got %s", entity.GetID(), retrievedEntity.GetID())
		}
	}
}

func TestFirestoreDB_DeletePerson(t *testing.T) {
	dbInstance := newFirestoreDB(t)
	testDelete(t, dbInstance, model.GeneratePerson, dbInstance.PutPerson, dbInstance.DeletePerson, dbInstance.GetPerson)
}

func TestFirestoreDB_DeleteFoodBank(t *testing.T) {
	dbInstance := newFirestoreDB(t)
	testDelete(t, dbInstance, model.GenerateFoodBank, dbInstance.PutFoodBank, dbInstance.DeleteFoodBank, dbInstance.GetFoodBank)
}

func TestFirestoreDB_DeleteFoodBankVisit(t *testing.T) {
	dbInstance := newFirestoreDB(t)
	testDelete(t, dbInstance, model.GenerateFoodBankVisit, dbInstance.PutFoodBankVisit, dbInstance.DeleteFoodBankVisit, dbInstance.GetFoodBankVisit)
}

func TestFirestoreDB_DeleteItem(t *testing.T) {
	dbInstance := newFirestoreDB(t)
	testDelete(t, dbInstance, model.GenerateItem, dbInstance.PutItem, dbInstance.DeleteItem, dbInstance.GetItem)
}

func testDelete[T model.Entity](t *testing.T, dbInstance *FirestoreDB, generateFunc func() (*T, error), putFunc func(context.Context, T) error, deleteFunc func(context.Context, string) error, getFunc func(context.Context, string) (*T, error)) {
	ctx := context.Background()

	entity, err := generateFunc()
	if err != nil {
		t.Fatalf("Failed to generate entity: %v", err)
	}

	err = putFunc(ctx, *entity)
	if err != nil {
		t.Fatalf("Failed to put entity: %v", err)
	}

	err = deleteFunc(ctx, entity.GetID())
	if err != nil {
		t.Fatalf("Failed to delete entity: %v", err)
	}

	_, err = getFunc(ctx, entity.GetID())
	if err == nil {
		t.Errorf("Expected error when retrieving deleted entity, got nil")
	}
}

func TestFirestoreDB_DeletePersons(t *testing.T) {
	dbInstance := newFirestoreDB(t)
	testDeletes(t, dbInstance, model.GeneratePeople, dbInstance.PutPersons, dbInstance.DeletePersons, dbInstance.GetPerson)
}

func TestFirestoreDB_DeleteFoodBanks(t *testing.T) {
	dbInstance := newFirestoreDB(t)
	testDeletes(t, dbInstance, model.GenerateFoodBanks, dbInstance.PutFoodBanks, dbInstance.DeleteFoodBanks, dbInstance.GetFoodBank)
}

func TestFirestoreDB_DeleteFoodBankVisits(t *testing.T) {
	dbInstance := newFirestoreDB(t)
	testDeletes(t, dbInstance, model.GenerateFoodBankVisits, dbInstance.PutFoodBankVisits, dbInstance.DeleteFoodBankVisits, dbInstance.GetFoodBankVisit)
}

func TestFirestoreDB_DeleteItems(t *testing.T) {
	dbInstance := newFirestoreDB(t)
	testDeletes(t, dbInstance, model.GenerateItems, dbInstance.PutItems, dbInstance.DeleteItems, dbInstance.GetItem)
}

func testDeletes[T model.Entity](t *testing.T, dbInstance *FirestoreDB, generateFunc func(int) ([]T, error), putFunc func(context.Context, []T) error, deleteFunc func(context.Context, []string) error, getFunc func(context.Context, string) (*T, error)) {
	ctx := context.Background()

	entities, err := generateFunc(5)
	if err != nil {
		t.Fatalf("Failed to generate entities: %v", err)
	}

	err = putFunc(ctx, entities)
	if err != nil {
		t.Fatalf("Failed to put entities: %v", err)
	}

	ids := make([]string, len(entities))
	for i, entity := range entities {
		ids[i] = entity.GetID()
	}

	err = deleteFunc(ctx, ids)
	if err != nil {
		t.Fatalf("Failed to delete entities: %v", err)
	}

	for _, id := range ids {
		_, err := getFunc(ctx, id)
		if err == nil {
			t.Errorf("Expected error when retrieving deleted entity with ID %s, got nil", id)
		}
	}
}
