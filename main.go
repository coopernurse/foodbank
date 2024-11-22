package main

import (
	"context"
	"cupboard/internal/db"
	"cupboard/internal/email"
	"cupboard/internal/middleware"
	"cupboard/internal/model"
	"cupboard/internal/ui"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	echomid "github.com/labstack/echo/v4/middleware"
	"github.com/oklog/ulid/v2"
	"google.golang.org/api/iterator"
)

var firestoreClient *firestore.Client

func main() {
	// Initialize Echo and middlewares
	e := echo.New()
	e.Use(echomid.Logger())
	e.Use(echomid.Recover())

	e.Static("/static", "static")

	// Initialize Firestore client
	ctx := context.Background()
	projectID := "uppervalleymend"
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	firestoreClient = client
	defer firestoreClient.Close()

	dbInstance := db.NewFirestoreDB(firestoreClient)

	// Define routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	})
	e.GET("/households", getHouseholds)

	signupPage := &ui.SignupPage{DB: dbInstance}
	householdListPage := &ui.HouseholdListPage{DB: dbInstance}
	householdDetailPage := &ui.HouseholdDetailPage{DB: dbInstance}

	e.GET("/signup", signupPage.GET)
	e.POST("/signup", signupPage.POST)
	e.GET("/households", householdListPage.GET)
	e.GET("/household/:id", householdDetailPage.GET)

	// Email route for testing
	e.POST("/send-email", sendEmailHandler)

	// Protected routes
	authenticated := e.Group("/protected")
	authenticated.Use(middleware.AuthMiddleware)
	authenticated.GET("", func(c echo.Context) error {
		return c.String(http.StatusOK, "Protected route")
	})

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// POST /send-email
func sendEmailHandler(c echo.Context) error {
	type EmailRequest struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
		Content string `json:"content"`
	}

	var req EmailRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
	}

	if err := email.SendEmail(c.Request().Context(), req.To, req.Subject, req.Content); err != nil {
		fmt.Printf("sendEmail ERR: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send email"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Email sent successfully"})
}

// POST /household
func postHousehold(c echo.Context) error {
	var household model.Household
	if err := c.Bind(&household); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
	}

	if household.Id == "" {
		household.Id = ulid.Make().String()
	}

	// Store household data in Firestore
	_, err := firestoreClient.Collection("households").Doc(household.Id).Set(context.Background(), household)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"error": fmt.Sprintf("Failed to save household: %v", err)})
	}
	return c.JSON(http.StatusCreated, household)
}

// GET /households
func getHouseholds(c echo.Context) error {
	households := []map[string]interface{}{}
	iter := firestoreClient.Collection("households").OrderBy("Id", firestore.Asc).Documents(context.Background())

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				map[string]string{"error": fmt.Sprintf("Failed to retrieve households: %v", err)})
		}

		var household model.Household
		if err := doc.DataTo(&household); err != nil {
			return c.JSON(http.StatusInternalServerError,
				map[string]string{"error": fmt.Sprintf("Data parsing error: %v", err)})
		}

		// Add selected fields for each household
		households = append(households, map[string]interface{}{
			"id":        household.Id,
			"firstName": household.Head.FirstName,
			"lastName":  household.Head.LastName,
			"dob":       household.Head.DOB,
			"members":   strconv.Itoa(len(household.Members)),
		})
	}

	return c.JSON(http.StatusOK, households)
}
