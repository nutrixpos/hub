package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/nutrixpos/hub/modules/hub/models"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TenantService struct {
	Config   config.Config
	Logger   logger.ILogger
	Settings models.Settings
}

// GetSettings returns the settings from the database
func (ts *TenantService) AddAPIKey(tenant_id string, title string, expirationDate time.Time) (err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", ts.Config.Databases[0].Host, ts.Config.Databases[0].Port))

	db_connection_deadline := 5 * time.Second
	if ts.Config.Env == "dev" {
		db_connection_deadline = 1000 * time.Second
	}

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), db_connection_deadline)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return
	}

	// Connected successfully

	api_key, err := GenerateAPIKey(32)
	if err != nil {
		return err
	}

	// terminal for pos terminals like restaurant or vending macine
	api_key = fmt.Sprintf("terminal-%s", api_key)

	collection := client.Database(ts.Config.Databases[0].Database).Collection("api_keys")
	_, err = collection.InsertOne(ctx, bson.M{
		"tenant_id":       tenant_id,
		"api_key":         api_key,
		"title":           title,
		"expiration_date": expirationDate,
	})
	if err != nil {
		return err
	}

	return nil
}

// GenerateAPIKey creates a new API key
func GenerateAPIKey(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(randomBytes), nil
}
