package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/nutrixpos/hub/common"
	"github.com/nutrixpos/hub/common/config"
	"github.com/nutrixpos/hub/modules/hub/models"
	"github.com/nutrixpos/pos/common/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateSettings is a post request handler that updates the settings in the database
// send a models.Settings directory to body to use it.
func WorkflowPUT(config config.Config, logger logger.ILogger) http.HandlerFunc {
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
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if request.Data == nil {
			http.Error(w, "data is required", http.StatusBadRequest)
			return
		}

		var workflow models.Workflow
		decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			DecodeHook: common.StringToTimeHook(),
			Result:     &workflow,
		})

		err = decoder.Decode(request.Data)

		if err != nil {
			logger.Error(fmt.Sprintf("Failed to decode log: %v", err))
			http.Error(w, "Failed to decode log", http.StatusBadRequest)
			return
		}

		db_workflow := map[string]interface{}{}
		db_workflow["actions"] = make([]interface{}, 0)

		// add new workflow
		if workflow.ID == "" {

			switch workflow.Trigger.Type {
			case models.WorkflowTriggerTypeLowStockLabel:
				var trigger models.WorkflowLowStockTrigger
				err = mapstructure.Decode(request.Data["trigger"], &trigger)
				if err != nil {
					logger.Error(fmt.Sprintf("Failed to decode low stock trigger: %v", err))
					http.Error(w, "Failed to decode low stock trigger", http.StatusBadRequest)
					return
				}
				db_workflow["trigger"] = trigger
				break
			}

			for index, action := range workflow.Actions {
				switch action.Type {
				case models.WorkflowActionTypeN8nWebhookLabel:
					var webhookAction models.WorkflowN8nWebhookAction
					err = mapstructure.Decode(request.Data["actions"].([]interface{})[index], &webhookAction)
					if err != nil {
						logger.Error(fmt.Sprintf("Failed to decode n8n webhook action: %v", err))
						http.Error(w, "Failed to decode n8n webhook action", http.StatusBadRequest)
						return
					}
					db_workflow["actions"] = append(db_workflow["actions"].([]interface{}), webhookAction)
				}
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

			// Access the collection
			collection := client.Database(config.Databases[0].Database).Collection(config.Databases[0].Tables["sales"])
			// Define the filter
			filter := bson.M{
				"tenant_id": tenant_id,
			}

			_, err = collection.UpdateOne(ctx, filter, bson.M{"$push": bson.M{"workflows": db_workflow}})
			if err != nil {
				http.Error(w, "Failed to insert workflow", http.StatusInternalServerError)
				logger.Error(fmt.Sprintf("ERROR: %v", err))
				return
			}

			w.WriteHeader(http.StatusOK)
			return
		}

	}
}
