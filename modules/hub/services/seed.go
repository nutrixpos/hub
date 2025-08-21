package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nutrixpos/hub/common/config"
	"github.com/nutrixpos/hub/modules/hub/models"
	"github.com/nutrixpos/pos/common/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SeederService struct {
	Config *config.Config
	Logger logger.ILogger
}

func (s *SeederService) Seed() error {
	err := s.SeedSettings()
	if err != nil {
		return err
	}

	return nil
}

func (s *SeederService) SeedSettings() error {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", s.Config.Databases[0].Host, s.Config.Databases[0].Port))

	deadline := 5 * time.Second
	if s.Config.Env == "dev" {
		deadline = 1000 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	// connected to db

	// Check if the "hubsync" collection exists
	db := client.Database(s.Config.Databases[0].Database)
	collectionNames, err := db.ListCollectionNames(ctx, bson.M{"name": "settings"})
	if err != nil {
		return err
	}

	if len(collectionNames) == 0 {

		err = db.CreateCollection(ctx, "settings")
		if err != nil {
			return err
		}

		// Insert a simple document into the "hubsync" collection
		hubsyncCollection := db.Collection("settings")
		_, err = hubsyncCollection.InsertOne(ctx, models.Settings{
			Id: primitive.NewObjectID().Hex(),
			Language: models.LanguageSettings{
				Code:     "en",
				Language: "English",
			},
		})
		if err != nil {
			return err
		}

	}

	return nil
}
