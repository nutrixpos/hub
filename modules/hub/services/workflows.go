package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/nutrixpos/hub/common"
	"github.com/nutrixpos/hub/common/config"
	"github.com/nutrixpos/hub/modules/hub/events"
	"github.com/nutrixpos/hub/modules/hub/models"
	"github.com/nutrixpos/pos/common/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WorkflowsService struct {
	Config   config.Config
	Logger   logger.ILogger
	Settings models.Settings
}

func (ws *WorkflowsService) AddLogsToWorkflowRun(source string, tenant_id string, workflow_id string, run_id string, log models.WorkflowRunLog) (err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", ws.Config.Databases[0].Host, ws.Config.Databases[0].Port))

	db_connection_deadline := 5 * time.Second
	if ws.Config.Env == "dev" {
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

	collection := client.Database(ws.Config.Databases[0].Database).Collection(ws.Config.Databases[0].Tables["sales"])

	filter := bson.M{
		"tenant_id":         tenant_id,
		"workflows.id":      workflow_id,
		"workflows.runs.id": run_id,
	}

	// Create the update to push log to the specific run
	// We need to use array filters to target the correct nested elements
	update := bson.M{
		"$push": bson.M{
			"workflows.$[workflow].runs.$[run].logs": log,
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"workflow.id": workflow_id},
			bson.M{"run.id": run_id},
		},
	}

	opts := options.Update().SetArrayFilters(arrayFilters)

	// Execute the update
	result, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		ws.Logger.Error(err.Error())
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (ws *WorkflowsService) RunLowStockTriggeredWorkflows(events []events.EventLowStockData) (err error) {

	run_id := primitive.NewObjectID().Hex()

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", ws.Config.Databases[0].Host, ws.Config.Databases[0].Port))

	db_connection_deadline := 5 * time.Second
	if ws.Config.Env == "dev" {
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

	collection := client.Database(ws.Config.Databases[0].Database).Collection(ws.Config.Databases[0].Tables["sales"])

	// Filter for specific tenant
	filter := bson.M{
		"tenant_id":              events[0].TenantId,
		"workflows.trigger.type": models.WorkflowTriggerTypeLowStockLabel,
	}

	// Projection to select only workflows field
	opts := options.FindOne().SetProjection(bson.M{
		"workflows.$": 1,
	})

	// Using bson.M for dynamic structure
	var result bson.M
	err = collection.FindOne(context.Background(), filter, opts).Decode(&result)

	if err == mongo.ErrNoDocuments {
		return mongo.ErrNoDocuments
	}

	if err != nil {
		ws.Logger.Error(err.Error())
		return err
	}

	// 1. Assert to primitive.A (which is []interface{})
	rawWorkflows, ok := result["workflows"].(primitive.A)
	if !ok {
		return fmt.Errorf("result['workflows'] is not an array or is missing")
	}

	// 2. Check length safely
	if len(rawWorkflows) == 0 {
		return fmt.Errorf("no workflows found")
	}

	// 3. Convert to []bson.M manually
	workflows := make([]bson.M, len(rawWorkflows))
	for i, v := range rawWorkflows {
		workflows[i] = v.(bson.M) // Assert individual elements
	}

	output := models.WorkflowLowStockTriggerOutput{
		Items: make([]models.WorkflowLowStockTriggerOutputItem, 0),
	}

	for _, w := range workflows {
		var workflow models.Workflow
		if b, err := bson.Marshal(w); err == nil {
			_ = bson.Unmarshal(b, &workflow)
		}

		if workflow.Trigger.Type == models.WorkflowTriggerTypeLowStockLabel {
			var trigger models.WorkflowLowStockTrigger
			if b, err := bson.Marshal(w["trigger"]); err == nil {
				_ = bson.Unmarshal(b, &trigger)
			}

			if trigger.MonitorType == models.TriggerLowStockMonitorTypeAny {

				newRun := models.WorkflowRun{
					ID:        run_id,
					StartTime: time.Now(),
					Status:    "running",
					Logs: []models.WorkflowRunLog{
						{
							Level:     "INFO",
							TimeStamp: time.Now(),
							Message:   "Workflow execution started.",
						},
					},
				}

				_, err = collection.UpdateOne(ctx, bson.M{
					"tenant_id":    events[0].TenantId,
					"workflows.id": workflow.ID,
				}, bson.M{
					"$push": bson.M{
						"workflows.$.runs": newRun,
					},
				})
				if err != nil {
					ws.Logger.Error(err.Error())
				}
			} else {
				fmt.Printf("Specific items: %+v\n", trigger)
			}
		}

		ws.AddLogsToWorkflowRun(models.WorkflowTriggerTypeLowStockLabel, events[0].TenantId, workflow.ID, run_id, models.WorkflowRunLog{
			Level:     "INFO",
			Message:   fmt.Sprintf("Finished evaluating the %s trigger, running actions...", models.WorkflowTriggerTypeLowStockLabel),
			TimeStamp: time.Now(),
		})

		rawActions, ok := w["actions"].(primitive.A)
		if !ok {
			return fmt.Errorf("result['actions'] is not an array or is missing")
		}

		// 2. Check length safely
		if len(rawActions) == 0 {
			return fmt.Errorf("no actions found")
		}

		// 3. Convert to []bson.M manually
		actions_bson := make([]bson.M, len(rawActions))
		for i, v := range rawActions {
			actions_bson[i] = v.(bson.M) // Assert individual elements
		}

		actions_base := make([]models.WorkflowActionBase, 0)

		for _, a := range actions_bson {
			var action models.WorkflowActionBase
			if b, err := bson.Marshal(a); err == nil {
				_ = bson.Unmarshal(b, &action)
			}

			actions_base = append(actions_base, action)
		}

		if actions_base[0].Type == models.WorkflowActionTypeN8nWebhookLabel {

			var next_action models.WorkflowN8nWebhookAction
			if b, err := bson.Marshal(actions_bson[0]); err == nil {
				_ = bson.Unmarshal(b, &next_action)
			}

			err = ws.RunN8nAction(output, next_action, actions_bson[1:], events[0].TenantId, workflow.ID, run_id)
			if err != nil {
				ws.Logger.Error(err.Error())

				ws.AddLogsToWorkflowRun(
					models.WorkflowActionTypeN8nWebhookLabel,
					events[0].TenantId,
					workflow.ID,
					run_id,
					models.WorkflowRunLog{
						Level:     "ERROR",
						Message:   err.Error(),
						TimeStamp: time.Now(),
					},
				)

				filter := bson.M{
					"tenant_id":         events[0].TenantId,
					"workflows.id":      workflow.ID,
					"workflows.runs.id": run_id,
				}

				// Create the update to push log to the specific run
				// We need to use array filters to target the correct nested elements
				update := bson.M{
					"$set": bson.M{
						"workflows.$[workflow].runs.$[run].end_time": time.Now(),
						"workflows.$[workflow].runs.$[run].status":   "failed",
					},
				}

				arrayFilters := options.ArrayFilters{
					Filters: []interface{}{
						bson.M{"workflow.id": workflow.ID},
						bson.M{"run.id": run_id},
					},
				}

				opts := options.Update().SetArrayFilters(arrayFilters)

				// Execute the update
				result, err := collection.UpdateOne(context.Background(), filter, update, opts)
				if err != nil {
					ws.Logger.Error(err.Error())
					return err
				}

				if result.MatchedCount == 0 {
					return mongo.ErrNoDocuments
				}
			}
		}

	}

	return err
}

func (ws *WorkflowsService) RunN8nAction(input interface{}, action models.WorkflowN8nWebhookAction, next_actions []bson.M, tenant_id string, workflow_id string, run_id string) error {

	ws.AddLogsToWorkflowRun(models.WorkflowActionTypeN8nWebhookLabel, tenant_id, workflow_id, run_id, models.WorkflowRunLog{
		Level:     "INFO",
		Message:   "Running N8n action...",
		TimeStamp: time.Now(),
	})

	if len(next_actions) == 0 {

		clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", ws.Config.Databases[0].Host, ws.Config.Databases[0].Port))

		db_connection_deadline := 5 * time.Second
		if ws.Config.Env == "dev" {
			db_connection_deadline = 1000 * time.Second
		}

		// Create a context with a timeout (optional)
		ctx, cancel := context.WithTimeout(context.Background(), db_connection_deadline)
		defer cancel()

		// Connect to MongoDB
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			return err
		}

		// Ping the database to check connectivity
		err = client.Ping(ctx, nil)
		if err != nil {
			return err
		}

		// Connected successfully
		collection := client.Database(ws.Config.Databases[0].Database).Collection(ws.Config.Databases[0].Tables["sales"])

		filter := bson.M{
			"tenant_id": tenant_id,
		}

		// Get the tenant from the database
		var tenant models.Tenant
		err = collection.FindOne(ctx, filter).Decode(&tenant)
		if err != nil {
			return err
		}

		webhook_url_preprocessed := action.WebhookURL

		secrets_map := make(map[string]interface{})
		for _, env_var := range tenant.EnvVars {
			if env_var.IsSecret {
				secrets_map[env_var.Name] = env_var.Value
			}
		}

		webhook_url_processed, err := common.InterpretVars(webhook_url_preprocessed, secrets_map)
		if err != nil {
			return fmt.Errorf("Couldn't interpret env vars")
		}

		filter = bson.M{
			"tenant_id":         tenant_id,
			"workflows.id":      workflow_id,
			"workflows.runs.id": run_id,
		}

		// Create the update to push log to the specific run
		// We need to use array filters to target the correct nested elements
		update := bson.M{
			"$set": bson.M{
				"workflows.$[workflow].runs.$[run].end_time": time.Now(),
				"workflows.$[workflow].runs.$[run].status":   "completed",
			},
		}

		arrayFilters := options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"workflow.id": workflow_id},
				bson.M{"run.id": run_id},
			},
		}

		opts := options.Update().SetArrayFilters(arrayFilters)

		// Execute the update
		result, err := collection.UpdateOne(context.Background(), filter, update, opts)
		if err != nil {
			ws.Logger.Error(err.Error())
			return err
		}

		if result.MatchedCount == 0 {
			return mongo.ErrNoDocuments
		}

		// Send a POST request to the webhook URL with the input
		jsonData, err := json.Marshal(input)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return err
		}

		// Create HTTP request
		url := webhook_url_processed
		http_req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {

			secrets := make([]string, 0)
			for _, env_var := range secrets_map {
				secrets = append(secrets, fmt.Sprintf("%s", env_var))
			}

			fmt.Println("Error creating request: %s", common.MaskString(err.Error(), secrets))
			return err
		}

		// Set headers
		http_req.Header.Set("Content-Type", "application/json")
		http_req.Header.Set("User-Agent", "Go-HTTP-Client")

		// Create HTTP client with timeout
		http_client := &http.Client{
			Timeout: 10 * time.Second,
		}

		// Send request
		http_resp, err := http_client.Do(http_req)
		if err != nil {

			secrets := make([]string, 0)
			for _, env_var := range secrets_map {
				secrets = append(secrets, fmt.Sprintf("%s", env_var))
			}

			fmt.Println("Error sending request: %s", common.MaskString(err.Error(), secrets))
			return fmt.Errorf("Error sending request: %s", common.MaskString(err.Error(), secrets))
		}
		defer http_resp.Body.Close()

		// Read response
		var http_result map[string]interface{}
		json.NewDecoder(http_resp.Body).Decode(&http_result)

		if http_resp.StatusCode == 200 {
			ws.AddLogsToWorkflowRun(
				models.WorkflowActionTypeN8nWebhookLabel,
				tenant_id,
				workflow_id,
				run_id,
				models.WorkflowRunLog{
					Level:     "INFO",
					Message:   "Successfully called N8n webhook, response: " + fmt.Sprintf("%v\n", http_result),
					TimeStamp: time.Now(),
				},
			)
		} else {
			var responseBody map[string]interface{}
			err = json.NewDecoder(http_resp.Body).Decode(&responseBody)
			if err != nil && err.Error() != "EOF" {
				return fmt.Errorf(fmt.Sprintf("Failed to call N8n webhook, status code: %s", http_resp.Status))
			}
			ws.FailWorkflow(tenant_id, workflow_id, run_id, fmt.Sprintf("Failed to call N8n webhook, status code: %s", http_resp.Status))
			return fmt.Errorf(fmt.Sprintf("Failed to call N8n webhook, status code: %s, response: %v", http_resp.Status, responseBody))
		}

		// POST request finished

		ws.AddLogsToWorkflowRun(
			models.WorkflowActionTypeN8nWebhookLabel,
			tenant_id,
			workflow_id,
			run_id,
			models.WorkflowRunLog{
				Level:     "INFO",
				Message:   "Successfully finished running the workflow",
				TimeStamp: time.Now(),
			},
		)
	}

	return nil
}

func (ws *WorkflowsService) ReplaceEnvVars(plain string, tenant_id string) (interpreted string, err error) {

	re := regexp.MustCompile(`\{\{\s*(.*?)\s*\}\}`)
	matches := re.FindAllStringSubmatch(plain, -1)

	var env_var_names []string
	for _, match := range matches {
		if len(match) > 1 {
			env_var_names = append(env_var_names, match[1])
		}
	}

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", ws.Config.Databases[0].Host, ws.Config.Databases[0].Port))

	db_connection_deadline := 5 * time.Second
	if ws.Config.Env == "dev" {
		db_connection_deadline = 1000 * time.Second
	}

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), db_connection_deadline)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return interpreted, err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return interpreted, err
	}

	// Connected successfully

	collection := client.Database(ws.Config.Databases[0].Database).Collection(ws.Config.Databases[0].Tables["sales"])

	filter := bson.M{
		"tenant_id": tenant_id,
	}

	env_var_values := make([]string, 0)

	// Get the tenant from the database
	var tenant models.Tenant
	err = collection.FindOne(ctx, filter).Decode(&tenant)
	if err != nil {
		return interpreted, err
	}

	for _, name := range env_var_names {

		found := false

		for _, env_var := range tenant.EnvVars {

			if env_var.Name == name {
				env_var_values = append(env_var_values, env_var.Value)
				found = true
				break
			}
		}

		if !found {
			return interpreted, fmt.Errorf("environment variable %s not found for tenant %s", name, tenant_id)
		}
	}

	interpreted = plain
	for _, value := range env_var_values {

		re := regexp.MustCompile(`\{\{\s*(.*?)\s*\}\}`)

		// Find the first match and its submatch
		loc := re.FindStringSubmatchIndex(interpreted)

		if loc == nil {
			break
		}

		// loc[0], loc[1] = start/end of full match
		// loc[2], loc[3] = start/end of first submatch (content inside {{}})

		start, end := loc[0], loc[1]
		interpreted = interpreted[:start] + value + interpreted[end:]
	}

	return interpreted, nil
}

func (ws *WorkflowsService) FailWorkflow(tenant_id string, workflow_id string, run_id string, message string) (err error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", ws.Config.Databases[0].Host, ws.Config.Databases[0].Port))

	db_connection_deadline := 5 * time.Second
	if ws.Config.Env == "dev" {
		db_connection_deadline = 1000 * time.Second
	}

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), db_connection_deadline)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	// Connected successfully

	collection := client.Database(ws.Config.Databases[0].Database).Collection(ws.Config.Databases[0].Tables["sales"])

	ws.AddLogsToWorkflowRun(
		models.WorkflowActionTypeN8nWebhookLabel,
		tenant_id,
		workflow_id,
		run_id,
		models.WorkflowRunLog{
			Level:     "ERROR",
			Message:   message,
			TimeStamp: time.Now(),
		},
	)

	filter := bson.M{
		"tenant_id":         tenant_id,
		"workflows.id":      workflow_id,
		"workflows.runs.id": run_id,
	}

	// Create the update to push log to the specific run
	// We need to use array filters to target the correct nested elements
	update := bson.M{
		"$set": bson.M{
			"workflows.$[workflow].runs.$[run].end_time": time.Now(),
			"workflows.$[workflow].runs.$[run].status":   "failed",
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"workflow.id": workflow_id},
			bson.M{"run.id": run_id},
		},
	}

	opts := options.Update().SetArrayFilters(arrayFilters)

	// Execute the update
	result, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		ws.Logger.Error(err.Error())
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
