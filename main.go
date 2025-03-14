package main

import (
	"context"
	"foodbank/internal/db"
	"foodbank/internal/ui"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
	echomid "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

	// Define routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	})

	signupPage := &ui.SignupPage{DB: dbInstance}
	householdListPage := &ui.HouseholdListPage{DB: dbInstance}
	householdDetailPage := &ui.HouseholdDetailPage{DB: dbInstance}

	e.GET("/signup", signupPage.GET)
	e.POST("/signup", signupPage.POST)
	e.GET("/households", householdListPage.GET)
	e.GET("/household/:id", householdDetailPage.GET)

	// Start server
	log.Info().Msg("Starting server on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}
