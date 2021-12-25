package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/daparadoks/go-fuel-api/internal/consumption"
	"github.com/daparadoks/go-fuel-api/internal/database"
	"github.com/daparadoks/go-fuel-api/internal/member"
	transportHTTP "github.com/daparadoks/go-fuel-api/internal/transport/http"

	log "github.com/sirupsen/logrus"
)

// App - the struct which contains things like pointers
// to database connections
type App struct {
	Name    string
	Version string
}

// Run - sets up our application
func (app *App) Run() error {
	log.SetFormatter(&log.JSONFormatter{})
	log.WithFields(
		log.Fields{
			"AppName":    app.Name,
			"AppVersion": app.Version,
		}).Info("Setting Up Our APP")

	var err error
	db, err := database.NewDatabase()
	if err != nil {
		return err
	}
	err = database.MigrateDB(db)
	if err != nil {
		return err
	}
	rand.Seed(time.Now().UnixNano())

	memberService := member.NewService(db)
	consumptionService := consumption.NewService(db)

	handler := transportHTTP.NewHandler(memberService, consumptionService)
	handler.SetupRoutes()

	if err := http.ListenAndServe(":3434", handler.Router); err != nil {
		fmt.Println("Failed to set up server: " + err.Error())
		return err
	}

	return nil
}

func main() {
	app := App{
		Name:    "Comment API",
		Version: "1.0",
	}
	if err := app.Run(); err != nil {
		log.Error("Error starting up our REST API")
		log.Fatal(err)
	}
}
