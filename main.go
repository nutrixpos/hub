package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/hub/modules"
	"github.com/nutrixpos/hub/modules/hub"
	hub_services "github.com/nutrixpos/hub/modules/hub/services"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
)

func main() {

	// Initialize the logger using ZeroLog
	logger := logger.NewZeroLog()

	// Create the configuration using the Viper config backend
	conf := config.ConfigFactory("viper", "config.yaml", &logger)

	seeder_svc := hub_services.SeederService{
		Config: &conf,
	}

	// make sure that settings bootstrapping data exists, it's idempotent
	err := seeder_svc.SeedSettings()
	if err != nil {
		panic(err)
	}

	settings_svc := hub_services.SettingsService{
		Config: conf,
	}

	// Load settings from the database
	settings, err := settings_svc.GetSettings()

	if err != nil {
		// Log and panic if settings can't be loaded
		logger.Error(err.Error())
		panic("Can't load settings from DB")
	}

	// Log successful database connection
	logger.Info("Successfully connected to DB")

	router := mux.NewRouter()

	// Initialize the app manager with logger
	appmanager := modules.AppManager{
		Logger: &logger,
	}

	// Load the core module, register HTTP handlers and background workers, and save the module
	appmanager.LoadModule(&hub.HubModule{
		Logger:   &logger,
		Config:   conf,
		Settings: settings,
	}, "hub").RegisterHttpHandlers(router).RegisterBackgroundWorkers().Save()

	// Ignite the app manager to start all modules
	appmanager.Run()

	srv := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:8001",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
