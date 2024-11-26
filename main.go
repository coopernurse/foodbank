package main

import (
	"context"
	"fmt"
	"foodbank/internal/config"
	"foodbank/internal/db"
	"foodbank/internal/email"
	"foodbank/internal/middleware"
	"foodbank/internal/model"
	"foodbank/internal/routes"
	"foodbank/internal/ui"
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
	client, err := firestore.NewClient(ctx, "uppervalleymend")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Firestore client")
	}
	firestoreClient = client
	defer firestoreClient.Close()

	dbInstance := db.NewFirestoreDB(firestoreClient)

	// Initialize Config
	config := config.Config{
		ServerURL:  os.Getenv("SERVER_URL"),
		SessionKey: os.Getenv("SESSION_KEY"),
		ProjectID:  "uppervalleymend",
	}
	if config.ServerURL == "" {
		log.Fatal().Msg("SERVER_URL environment variable is not set")
	}
	if config.SessionKey == "" {
		log.Fatal().Msg("SESSION_KEY environment variable is not set")
	}
	if len(config.SessionKey) != 32 {
		log.Fatal().Msg("SESSION_KEY must be 32 bytes long")
	}

	// Initialize real email sender
	realEmailSender := &email.RealEmailSender{}

	// Initialize routes
	personsHandler, foodBanksHandler, itemsHandler, visitsHandler, authHandler := routes.NewRoutes(dbInstance,
		realEmailSender, config)

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

	// Register routes for each handler
	personsHandler.RegisterRoutes(e)
	foodBanksHandler.RegisterRoutes(e)
	itemsHandler.RegisterRoutes(e)
	visitsHandler.RegisterRoutes(e)
	authHandler.RegisterRoutes(e)

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
