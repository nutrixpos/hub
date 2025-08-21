package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/hub/common/config"
	"github.com/nutrixpos/hub/modules/hub/models"
	"github.com/nutrixpos/pos/common/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type KoptanService struct {
	Config   config.Config
	Settings models.Settings
	Logger   logger.ILogger
}

func (ks *KoptanService) GetSuggestions(tenant_id string, page_number int, page_size int) (suggestions []models.KoptanSuggestion, totalRecords int, err error) {
	// Implementation of fetching suggestions from the database or any other source
	// This is a placeholder for the actual implementation
	suggestions = make([]models.KoptanSuggestion, 0)

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", ks.Config.Databases[0].Host, ks.Config.Databases[0].Port))

	deadline := 5 * time.Second
	if ks.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return suggestions, totalRecords, err
	}
	// connected to db

	collection := client.Database(ks.Config.Databases[0].Database).Collection("koptan_suggestions")
	filter := bson.M{"tenant_id": tenant_id}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return suggestions, totalRecords, err
	}
	totalRecords = int(count)

	findOptions := options.Find()
	skip := (page_number - 1) * page_size
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(page_size))

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		ks.Logger.Error(err.Error())
		return suggestions, totalRecords, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var suggestion models.KoptanSuggestion
		err := cursor.Decode(&suggestion)
		if err != nil {
			ks.Logger.Error(err.Error())
			return suggestions, totalRecords, err
		}
		suggestions = append(suggestions, suggestion)
	}

	if err := cursor.Err(); err != nil {
		ks.Logger.Error(err.Error())
		return suggestions, totalRecords, err
	}

	return suggestions, totalRecords, err
}
