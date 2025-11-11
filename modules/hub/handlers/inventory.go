package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nutrixpos/hub/common/config"
	"github.com/nutrixpos/hub/modules/hub/models"
	"github.com/nutrixpos/pos/common/logger"
	core_handlers "github.com/nutrixpos/pos/modules/core/handlers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InventoryItemsDTO struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

func InventoryItemsPatch(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenant_id := "1"
		label := "dev"

		if config.Env != "dev" {
			token := r.Header.Get("X-Userinfo")
			if token == "" {
				http.Error(w, "X-Userinfo header is required", http.StatusBadRequest)
				return
			}

			decodedData, err := base64.StdEncoding.DecodeString(token)
			if err != nil {
				http.Error(w, "Failed to decode token", http.StatusBadRequest)
				logger.Error(fmt.Sprintf("ERROR: %v", err))
				return
			}

			var claims map[string]interface{}
			err = json.Unmarshal(decodedData, &claims)
			if err != nil {
				http.Error(w, "Failed to unmarshal token", http.StatusBadRequest)
				logger.Error(fmt.Sprintf("ERROR: %v", err))
				return
			}

			var ok bool
			tenant_id, ok = claims["tenant_id"].(string)
			if !ok || tenant_id == "" {
				http.Error(w, "tenant_id claim is required and must be a string", http.StatusBadRequest)
				logger.Error("ERROR: tenant_id claim is required and must be a string")
				return
			}

			label = claims["name"].(string)
			if label == "" {
				http.Error(w, "name claim is required", http.StatusBadRequest)
				logger.Error("ERROR: name claim is required")
				return
			}
		}

		request := struct {
			Data map[string]interface{} `json:"data"`
			Meta struct {
				Id string `json:"id"`
			}
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Set up MongoDB connection
		clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", config.Databases[0].Host, config.Databases[0].Port))
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}
		defer client.Disconnect(ctx)

		label = fmt.Sprintf("branch:%s", label)

		// Access the collection
		collection := client.Database(config.Databases[0].Database).Collection(config.Databases[0].Tables["sales"])

		// Define the filter
		filter := bson.M{
			"tenant_id": tenant_id,
		}

		// Find the document
		var existing models.Tenant
		err = collection.FindOne(ctx, filter).Decode(&existing)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "No document found with the specified tenant_id", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to retrieve document", http.StatusInternalServerError)
				logger.Error(fmt.Sprintf("ERROR: %v", err))
			}
			return
		}

		existing_item_index, err := getInventoryItemById(request.Meta.Id, existing.InventoryItems)
		new_items := existing.InventoryItems

		for key, value := range request.Data {
			switch key {
			case "alert_threshold":
				if threshold, ok := value.(float64); ok {
					new_items[existing_item_index].Settings.AlertThreshold = threshold
				} else {
					http.Error(w, "alert_threshold must be a number", http.StatusBadRequest)
					logger.Error("ERROR: alert_threshold must be a number")
					return
				}
			}

		}

		existing.InventoryItems = new_items

		// Update the document
		_, err = collection.ReplaceOne(ctx, filter, existing)
		if err != nil {
			http.Error(w, "Failed to update document", http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func InventoryItemsGet(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenant_id := "1"
		label := "dev"

		if config.Env != "dev" {
			token := r.Header.Get("X-Userinfo")
			if token == "" {
				http.Error(w, "X-Userinfo header is required", http.StatusBadRequest)
				return
			}

			decodedData, err := base64.StdEncoding.DecodeString(token)
			if err != nil {
				http.Error(w, "Failed to decode token", http.StatusBadRequest)
				logger.Error(fmt.Sprintf("ERROR: %v", err))
				return
			}

			// Create a map to hold the decoded JSON data
			var claims map[string]interface{}

			// Unmarshal the decoded data into the map
			err = json.Unmarshal(decodedData, &claims)
			if err != nil {
				http.Error(w, "Failed to unmarshal token", http.StatusBadRequest)
				logger.Error(fmt.Sprintf("ERROR: %v", err))
				return
			}

			var ok bool

			tenant_id, ok = claims["tenant_id"].(string)
			if !ok {
				http.Error(w, "tenant_id claim is required and must be a string", http.StatusBadRequest)
				logger.Error("ERROR: tenant_id claim is required and must be a string")
				return
			}

			if tenant_id == "" {
				http.Error(w, "tenant_id claim is required", http.StatusBadRequest)
				logger.Error("ERROR: tenant_id claim is required")
				return
			}

			label = claims["name"].(string)
			if label == "" {
				http.Error(w, "name claim is required", http.StatusBadRequest)
				logger.Error("ERROR: name claim is required")
				return
			}
		}

		filter := bson.D{{Key: "tenant_id", Value: tenant_id}}
		var result []models.Tenant
		client, err := mongo.Connect(r.Context(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", config.Databases[0].Host, config.Databases[0].Port)))
		if err != nil {
			http.Error(w, "Failed to connect to db", http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}
		defer client.Disconnect(r.Context())

		collection := client.Database(config.Databases[0].Database).Collection(config.Databases[0].Tables["sales"])
		cursor, err := collection.Find(context.TODO(), filter, &options.FindOptions{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = cursor.All(context.TODO(), &result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		inventoryItems := make([]models.InventoryItem, 0)
		total_records := 0

		if len(result) > 0 {

			total_records = len(result[0].InventoryItems)

			if result[0].Subscription.SubscriptionPlan == "free" {
				inventoryItems = append(inventoryItems, result[0].InventoryItems[0])
			} else {
				inventoryItems = result[0].InventoryItems
			}
		} else {
			http.Error(w, "No inventory items found", http.StatusNotFound)
			logger.Error("ERROR: No inventory items found")
			return
		}

		response := core_handlers.JSONApiOkResponse{
			Data: inventoryItems,
			Meta: core_handlers.JSONAPIMeta{
				TotalRecords: total_records,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func InventoryItemsPut(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tenant_id := "1"
		label := "dev"

		if config.Env != "dev" {
			token := r.Header.Get("X-Userinfo")
			if token == "" {
				http.Error(w, "X-Userinfo header is required", http.StatusBadRequest)
				return
			}

			decodedData, err := base64.StdEncoding.DecodeString(token)
			if err != nil {
				http.Error(w, "Failed to decode token", http.StatusBadRequest)
				logger.Error(fmt.Sprintf("ERROR: %v", err))
				return
			}

			// Create a map to hold the decoded JSON data
			var claims map[string]interface{}

			// Unmarshal the decoded data into the map
			err = json.Unmarshal(decodedData, &claims)
			if err != nil {
				http.Error(w, "Failed to unmarshal token", http.StatusBadRequest)
				logger.Error(fmt.Sprintf("ERROR: %v", err))
				return
			}

			var ok bool

			tenant_id, ok = claims["tenant_id"].(string)
			if !ok {
				http.Error(w, "tenant_id claim is required and must be a string", http.StatusBadRequest)
				logger.Error("ERROR: tenant_id claim is required and must be a string")
				return
			}

			if tenant_id == "" {
				http.Error(w, "tenant_id claim is required", http.StatusBadRequest)
				logger.Error("ERROR: tenant_id claim is required")
				return
			}

			label = claims["name"].(string)
			if label == "" {
				http.Error(w, "name claim is required", http.StatusBadRequest)
				logger.Error("ERROR: name claim is required")
				return
			}
		}

		client, err := mongo.Connect(r.Context(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", config.Databases[0].Host, config.Databases[0].Port)))
		if err != nil {
			http.Error(w, "Failed to connect to db", http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}
		defer client.Disconnect(r.Context())

		collection := client.Database(config.Databases[0].Database).Collection(config.Databases[0].Tables["sales"])

		// check if document with tenant_id exists, otherwise create it with empty object value
		filter_tenant := bson.D{{Key: "tenant_id", Value: tenant_id}}
		var result bson.M
		err = collection.FindOne(r.Context(), filter_tenant).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				_, err = collection.InsertOne(r.Context(), bson.D{{Key: "tenant_id", Value: tenant_id}, {Key: "sales", Value: []bson.D{}}})
				if err != nil {
					http.Error(w, "Error creating new tenant", http.StatusBadRequest)
					logger.Error(fmt.Sprintf("ERROR: %v", err))
					return
				}
			} else {
				http.Error(w, "Error searching for tenant in db", http.StatusBadRequest)
				logger.Error(fmt.Sprintf("ERROR: %v", err))
				return
			}
		}

		request_body := struct {
			Data []InventoryItemsDTO `json:"data"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&request_body); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}

		if len(request_body.Data) == 0 {
			w.WriteHeader(http.StatusOK)
			return
		}

		items := make([]models.InventoryItem, 0, len(request_body.Data))

		label = fmt.Sprintf("branch:%s", label)

		filter := bson.D{{Key: "tenant_id", Value: tenant_id}}

		// Remove all inventory_items with the same label
		removeFilter := bson.D{
			{Key: "tenant_id", Value: tenant_id},
			{Key: "inventory_items", Value: bson.D{
				{Key: "$elemMatch", Value: bson.D{
					{Key: "labels", Value: bson.D{{Key: "$in", Value: []string{label}}}},
				}},
			}},
		}

		// Select existing inventory_items object
		var existingInventoryItems models.Tenant
		err = collection.FindOne(r.Context(), filter).Decode(&existingInventoryItems)
		if err != nil {
			http.Error(w, "Failed to retrieve existing inventory items", http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}

		for _, item := range request_body.Data {

			existing_item_index, err := getInventoryItemById(item.ID, existingInventoryItems.InventoryItems)
			settings := models.InventoryItemSettings{
				AlertThreshold: 0,
				AlertEnabled:   false,
			}

			if err == nil {
				settings = existingInventoryItems.InventoryItems[existing_item_index].Settings
			}

			items = append(items, models.InventoryItem{
				ID:       item.ID,
				TenantID: tenant_id,
				Name:     item.Name,
				Quantity: item.Quantity,
				Unit:     item.Unit,
				Labels:   []string{label},
				Settings: settings,
			})
		}

		removeUpdate := bson.D{
			{Key: "$pull", Value: bson.D{
				{Key: "inventory_items", Value: bson.D{
					{Key: "labels", Value: label},
				}},
			}},
		}

		_, err = collection.UpdateOne(r.Context(), removeFilter, removeUpdate)
		if err != nil {
			http.Error(w, "Failed to remove items with the same label", http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}

		update := bson.D{{Key: "$addToSet", Value: bson.D{{Key: "inventory_items", Value: bson.D{{Key: "$each", Value: items}}}}}}
		_, err = collection.UpdateOne(r.Context(), filter, update)
		if err != nil {
			http.Error(w, "Failed to append items to the inventory_items property", http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}

	}
}

func getInventoryItemById(id string, items []models.InventoryItem) (index int, err error) {
	for index, item := range items {
		if item.ID == id {
			return index, nil
		}
	}
	return 0, fmt.Errorf("inventory item with ID %s not found", id)
}
