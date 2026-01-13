package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"github.com/nutrixpos/hub/common"
	"github.com/nutrixpos/hub/common/config"
	"github.com/nutrixpos/hub/modules/hub/models"
	"github.com/nutrixpos/pos/common/logger"
	core_handlers "github.com/nutrixpos/pos/modules/core/handlers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EnvVarPATCH(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenant_id := "1"
		label := "dev"

		params := mux.Vars(r)
		env_var := params["name"]

		if env_var == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

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

		request := struct {
			Data models.WorkflowEnvVar `json:"data"`
		}{}

		err = json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		filter := bson.M{
			"tenant_id":     tenant_id,
			"env_vars.name": request.Data.Name,
		}
		update := bson.M{}

		count, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			http.Error(w, "Failed to count documents", http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}

		if count == 0 {
			// Insert new env_var item if the item doesn't exist
			filter = bson.M{"tenant_id": tenant_id}
			update = bson.M{
				"$push": bson.M{
					"env_vars": bson.M{
						"name":      request.Data.Name,
						"value":     request.Data.Value,
						"is_secret": request.Data.IsSecret,
					},
				},
			}
		} else {
			// Update existing env_var item if the item exists
			update = bson.M{
				"$set": bson.M{
					"env_vars.$.value": request.Data.Value,
				},
			}
		}

		_, err = collection.UpdateOne(ctx, filter, update, options.Update())
		if err != nil {
			http.Error(w, "Failed to update environment variable", http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}
	}
}

func EnvVarsGet(config config.Config, logger logger.ILogger) http.HandlerFunc {
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

		filter := bson.M{"tenant_id": tenant_id}
		var result models.Tenant
		err = collection.FindOne(ctx, filter).Decode(&result)
		if err != nil {
			http.Error(w, "Failed to find tenant", http.StatusNotFound)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}

		env_vars := make([]models.WorkflowEnvVar, len(result.EnvVars))
		for i, env_var := range result.EnvVars {
			env_vars[i] = env_var

			if env_var.IsSecret {
				env_vars[i].Value = "********"
			}
		}

		response := core_handlers.JSONApiOkResponse{
			Meta: core_handlers.JSONAPIMeta{
				TotalRecords: len(env_vars),
			},
			Data: env_vars,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "Failed to marshal environment variables response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)

	}
}

func EnvVarDelete(config config.Config, logger logger.ILogger) http.HandlerFunc {
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

		filter := bson.M{"tenant_id": tenant_id}
		var result models.Tenant
		err = collection.FindOne(ctx, filter).Decode(&result)
		if err != nil {
			http.Error(w, "Failed to find tenant", http.StatusNotFound)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}

		params := mux.Vars(r)
		env_var := params["name"]

		if env_var == "" {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}
		update := bson.M{"$pull": bson.M{"env_vars": bson.M{"name": env_var}}}
		_, err = collection.UpdateOne(ctx, filter, update)
		if err != nil {
			http.Error(w, "Failed to delete environment variable", http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}

		w.WriteHeader(http.StatusOK)

	}
}

func WorkflowGET(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tenant_id := "1"
		label := "dev"

		params := mux.Vars(r)
		workflow_id := params["id"]

		if workflow_id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

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

		// Filter for specific tenant
		filter := bson.M{
			"tenant_id":    tenant_id,
			"workflows.id": workflow_id,
		}

		// Projection to select only workflows field
		opts := options.FindOne().SetProjection(bson.M{
			"workflows.$": 1,
		})

		// Using bson.M for dynamic structure
		var result bson.M
		err = collection.FindOne(context.Background(), filter, opts).Decode(&result)
		if err != nil {
			http.Error(w, "Failed to fetch workflows", http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}

		if len(result["workflows"].(primitive.A)) == 0 {
			http.Error(w, "No workflows found", http.StatusNotFound)
			return
		}

		response := core_handlers.JSONApiOkResponse{
			Meta: core_handlers.JSONAPIMeta{
				TotalRecords: len(result["workflows"].(primitive.A)),
			},
			Data: result["workflows"].(primitive.A)[0],
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "Failed to marshal order settings response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func WorkflowsGET(config config.Config, logger logger.ILogger) http.HandlerFunc {
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

		// Filter for specific tenant
		filter := bson.M{"tenant_id": tenant_id}

		// Projection to select only workflows field
		opts := options.FindOne().SetProjection(bson.M{
			"workflows": 1,
			"id":        1,
		})

		// Using bson.M for dynamic structure
		var result bson.M
		err = collection.FindOne(context.Background(), filter, opts).Decode(&result)
		if err != nil {
			http.Error(w, "Failed to fetch workflows", http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}

		response := core_handlers.JSONApiOkResponse{
			Meta: core_handlers.JSONAPIMeta{
				TotalRecords: len(result["workflows"].(primitive.A)),
			},
			Data: result["workflows"],
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "Failed to marshal order settings response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

func WorkflowPOST(config config.Config, logger logger.ILogger) http.HandlerFunc {
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

		deadline := 5 * time.Second
		if config.Env == "dev" {
			deadline = 1000 * time.Second
		}

		ctx, cancel := context.WithTimeout(context.Background(), deadline)
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
		db_workflow["id"] = primitive.NewObjectID().Hex()
		db_workflow["name"] = workflow.Name
		db_workflow["description"] = workflow.Description
		db_workflow["enabled"] = workflow.Enabled
		db_workflow["status"] = "idle"
		db_workflow["runs"] = make([]models.WorkflowRun, 0)

		// add new workflow
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
	}
}

func WorkflowPATCH(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenant_id := "1"
		label := "dev"

		// Get workflow ID from URL params
		params := mux.Vars(r)
		workflow_id := params["id"]

		if workflow_id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

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
			logger.Error(fmt.Sprintf("Failed to decode workflow: %v", err))
			http.Error(w, "Failed to decode workflow", http.StatusBadRequest)
			return
		}

		db_workflow := map[string]interface{}{}
		db_workflow["actions"] = make([]interface{}, 0)

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

		deadline := 5 * time.Second
		if config.Env == "dev" {
			deadline = 1000 * time.Second
		}

		ctx, cancel := context.WithTimeout(context.Background(), deadline)
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
		db_workflow["id"] = workflow_id
		db_workflow["name"] = workflow.Name
		db_workflow["description"] = workflow.Description
		db_workflow["enabled"] = workflow.Enabled
		db_workflow["runs"] = workflow.Runs

		// Update existing workflow
		filter := bson.M{
			"tenant_id":    tenant_id,
			"workflows.id": workflow_id,
		}

		update := bson.M{
			"$set": bson.M{
				"workflows.$": db_workflow,
			},
		}

		result, err := collection.UpdateOne(ctx, filter, update)
		if err != nil {
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			http.Error(w, "Failed to update workflow", http.StatusInternalServerError)
			return
		}

		if result.MatchedCount == 0 {
			http.Error(w, "Workflow not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
