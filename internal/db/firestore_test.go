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
