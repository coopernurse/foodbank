package main

import (
	"context"
	"cupboard/internal/db"
	"cupboard/internal/email"
	"cupboard/internal/middleware"
	"cupboard/internal/model"
	"cupboard/internal/routes"
	"cupboard/internal/ui"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	echomid "github.com/labstack/echo/v4/middleware"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
)

var firestoreClient *firestore.Client

func main() {
	// Initialize leveled logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

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
		log.Fatal().Err(err).Msg("Failed to create Firestore client")
	}
	firestoreClient = client
	defer firestoreClient.Close()

	dbInstance := db.NewFirestoreDB(firestoreClient)

	// Initialize real email sender
	realEmailSender := &email.RealEmailSender{}

	// Initialize routes
	personsHandler, foodBanksHandler, itemsHandler, visitsHandler := routes.NewRoutes(dbInstance, realEmailSender)

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

	// Register Persons routes
	e.POST("/person", personsHandler.PutPerson)
	e.GET("/persons/search", personsHandler.SearchPersons)
	e.GET("/household/:id/persons", personsHandler.LoadHouseholdPersons)
	e.POST("/person/:id/reset-password", personsHandler.ResetPassword)
	e.POST("/person/:id/email-login-link", personsHandler.EmailLoginLink)
	e.GET("/person/:id/permissions", personsHandler.ResolvePermissions)

	// Register FoodBanks routes
	e.POST("/foodbank", foodBanksHandler.PutFoodBank)
	e.GET("/foodbanks", foodBanksHandler.LoadFoodBanks)
	e.POST("/foodbank/:id/assign-permissions", foodBanksHandler.AssignPersonPermissions)

	// Register Items routes
	e.POST("/item", itemsHandler.PutItem)

	// Register Visits routes
	e.GET("/household/:id/visits", visitsHandler.LoadHouseholdVisits)
	e.GET("/household/:id/visits/limits", visitsHandler.ComputeItemLimits)
	e.POST("/household/:id/visit/:visitId/item", visitsHandler.AddItemToVisit)
	e.POST("/visit", visitsHandler.PutVisit)

	// Email route for testing
	e.POST("/send-email", sendEmailHandler)

	// Initialize AuthHandler
	authHandler := routes.NewAuthHandler(dbInstance, realEmailSender)

	// Add the /login route
	e.POST("/login", authHandler.Login)

	// Protected routes
	authenticated := e.Group("/protected")
	authenticated.Use(middleware.AuthMiddleware)
	authenticated.GET("", func(c echo.Context) error {
		return c.String(http.StatusOK, "Protected route")
	})

	// Start server
	log.Info().Msg("Starting server on :8080")
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

	sender := email.RealEmailSender{}
	if err := sender.SendEmail(c.Request().Context(), req.To, req.Subject, req.Content); err != nil {
		log.Error().Err(err).Msg("Failed to send email")
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
		log.Error().Err(err).Msg("Failed to save household")
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
			log.Error().Err(err).Msg("Failed to retrieve households")
			return c.JSON(http.StatusInternalServerError,
				map[string]string{"error": fmt.Sprintf("Failed to retrieve households: %v", err)})
		}

		var household model.Household
		if err := doc.DataTo(&household); err != nil {
			log.Error().Err(err).Msg("Data parsing error")
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
