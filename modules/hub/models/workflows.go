package models

const (
	WorkflowTriggerTypeLowStockLabel  = "trigger_low_stock"
	WorkflowActionTypeN8nWebhookLabel = "action_n8n_webhook"
)

type Workflow struct {
	ID          string               `json:"id" bson:"id" mapstructure:"id"`
	Name        string               `json:"name" bson:"name" mapstructure:"name"`
	Description string               `json:"description" bson:"description" mapstructure:"description"`
	Enabled     bool                 `json:"enabled" bson:"enabled" mapstructure:"enabled"`
	Trigger     WorkflowTriggerBase  `json:"trigger" bson:"trigger" mapstructure:"trigger"`
	Actions     []WorkflowActionBase `json:"actions" bson:"actions" mapstructure:"actions"`
}

type WorkflowTriggerBase struct {
	Type string `json:"type" bson:"type" mapstructure:"type"`
}

type WorkflowActionBase struct {
	Type string `json:"type" bson:"type" mapstructure:"type"`
}

type WorkflowLowStockTrigger struct {
	WorkflowTriggerBase `json:",inline" bson:",inline" mapstructure:",squash"`
	MonitorType         string   `json:"monitor_type" bson:"monitor_type" mapstructure:"monitor_type"`
	ProductIDs          []string `json:"product_ids" bson:"product_ids" mapstructure:"product_ids"`
	Threshold           float64  `json:"threshold" bson:"threshold" mapstructure:"threshold"`
	Output              string   `json:"output" bson:"output" mapstructure:"output"`
}

type WorkflowN8nWebhookAction struct {
	WorkflowActionBase `json:",inline" bson:",inline" mapstructure:",squash"`
	Input              string `json:"input" bson:"input" mapstructure:"input"`
	WebhookURL         string `json:"webhook_url" bson:"webhook_url" mapstructure:"webhook_url"`
	Method             string `json:"method" bson:"method" mapstructure:"method"`
	Output             string `json:"output" bson:"output" mapstructure:"output"`
}
