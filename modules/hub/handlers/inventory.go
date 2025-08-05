package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nutrixpos/hub/modules/hub/models"
	"github.com/nutrixpos/pos/common/config"
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
		if len(result) > 0 {
			inventoryItems = result[0].InventoryItems
		} else {
			http.Error(w, "No inventory items found", http.StatusNotFound)
			logger.Error("ERROR: No inventory items found")
			return
		}

		json.NewEncoder(w).Encode(inventoryItems)
		response := core_handlers.JSONApiOkResponse{
			Data: inventoryItems,
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

		for _, item := range request_body.Data {
			items = append(items, models.InventoryItem{
				ID:       item.ID,
				TenantID: tenant_id,
				Name:     item.Name,
				Quantity: item.Quantity,
				Labels:   []string{label},
			})
		}

		client, err := mongo.Connect(r.Context(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", config.Databases[0].Host, config.Databases[0].Port)))
		if err != nil {
			http.Error(w, "Failed to connect to db", http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}
		defer client.Disconnect(r.Context())

		collection := client.Database(config.Databases[0].Database).Collection(config.Databases[0].Tables["sales"])

		filter := bson.D{{Key: "tenant_id", Value: tenant_id}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "inventory_items", Value: items}}}}
		_, err = collection.UpdateOne(r.Context(), filter, update)
		if err != nil {
			http.Error(w, "Failed to insert items into the inventory_items property", http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}

	}
}
