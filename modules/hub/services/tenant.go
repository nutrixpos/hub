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

type TenantService struct {
	Config   config.Config
	Logger   logger.ILogger
	Settings models.Settings
}

// GetSettings returns the settings from the database
func (ts *TenantService) GetTenantById(tenant_id string) (tenant models.Tenant, err error) {

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

	// Get the tenant by tenant_id from the database
	tenant_collection := client.Database(ts.Config.Databases[0].Database).Collection(ts.Config.Databases[0].Tables["sales"])
	err = tenant_collection.FindOne(ctx, bson.M{"tenant_id": tenant_id}).Decode(&tenant)
	return
}
