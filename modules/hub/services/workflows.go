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

func (ws *WorkflowsService) AddLogsToWorkflowRun() (err error) {

	return nil
}

func (ws *WorkflowsService) RunLowStockTriggeredWorkflows(event events.EventLowStockData) (err error) {

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
		"tenant_id":              event.TenantId,
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
					ID:        primitive.NewObjectID().Hex(),
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
					"tenant_id":    event.TenantId,
					"workflows.id": workflow.ID,
				}, bson.M{
					"$push": bson.M{
						"workflows.$.runs": newRun,
					},
				})
				if err != nil {
					ws.Logger.Error(err.Error())
				}

				// TODO: Run workflow
				// Retrieve all inventory items from db with the mentioned tenant id
				filter := bson.M{
					"tenant_id": event.TenantId,
				}

				var result []models.Tenant
				var inventoryItems []models.InventoryItem
				collection := client.Database(ws.Config.Databases[0].Database).Collection(ws.Config.Databases[0].Tables["sales"])
				cursor, err := collection.Find(context.TODO(), filter, &options.FindOptions{})
				if err != nil {
					return err
				}

				err = cursor.All(context.TODO(), &result)
				if err != nil {
					return err
				}

				inventoryItems = make([]models.InventoryItem, 0)

				if len(result) > 0 {
					if result[0].InventoryItems == nil {
						result[0].InventoryItems = make([]models.InventoryItem, 0)
					}

					if len(result[0].InventoryItems) > 0 {
						inventoryItems = append(inventoryItems, result[0].InventoryItems[0])
					}

					for _, item := range inventoryItems {
						if item.Quantity <= item.Settings.AlertThreshold {
							output.Items = append(output.Items, models.WorkflowLowStockTriggerOutputItem{
								Labels:   item.Labels,
								TenantId: item.TenantID,
								ItemID:   item.ID,
								Quantity: item.Quantity,
								ItemName: item.Name,
								Unit:     item.Unit,
							})
						}
					}

				} else {
					return mongo.ErrNoDocuments
				}
			} else {
				fmt.Printf("Specific items: %+v\n", trigger)
			}
		}

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

		actions := make([]any, 0)

		for _, a := range actions_bson {
			var action interface{}
			if b, err := bson.Marshal(a); err == nil {
				_ = bson.Unmarshal(b, &action)
			}

			actions = append(actions, action)
		}

		if action, ok := actions[0].(models.WorkflowN8nWebhookAction); ok {
			ws.RunN8nAction(output, action, actions[1:])
		}
	}

	return err
}

func (ws *WorkflowsService) RunN8nAction(input interface{}, action models.WorkflowN8nWebhookAction, next_actions []interface{}) error {

	fmt.Println("Running n8n action")

	return nil
}
