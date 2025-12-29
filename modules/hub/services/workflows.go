package services

import (
	"context"
	"fmt"
	"time"

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

			ws.RunN8nAction(output, next_action, actions_bson[1:], events[0].TenantId, workflow.ID, run_id)
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
