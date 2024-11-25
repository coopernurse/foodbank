package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"foodbank/internal/db"
	"foodbank/internal/model"

	"cloud.google.com/go/firestore"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Define command line flags
	firstName := flag.String("firstname", "", "First name of the person")
	lastName := flag.String("lastname", "", "Last name of the person")
	email := flag.String("email", "", "Email of the person")
	password := flag.String("password", "", "Password of the person")

	flag.Parse()

	// Validate required flags
	if *firstName == "" || *lastName == "" || *email == "" || *password == "" {
		log.Fatal("All flags are required: firstname, lastname, email, password")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Create a new person
	person := model.Person{
		PersonCommon: model.PersonCommon{
			Id:        ulid.Make().String(),
			FirstName: *firstName,
			LastName:  *lastName,
			Email:     *email,
		},
		PasswordHash: string(hashedPassword),
	}

	// Initialize Firestore client
	ctx := context.Background()
	projectID := "uppervalleymend"
	firestoreClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create Firestore client")
	}
	defer firestoreClient.Close()

	// Create FirestoreDB instance
	dbInstance := db.NewFirestoreDB(firestoreClient)

	// Save the person to Firestore
	err = dbInstance.PutPerson(ctx, person)
	if err != nil {
		log.Fatalf("Failed to save person: %v", err)
	}

	fmt.Println("Person created successfully")
}
