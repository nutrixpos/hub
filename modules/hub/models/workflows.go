package models

import (
	"time"
)

const (
	WorkflowTriggerTypeLowStockLabel   = "trigger_low_stock"
	WorkflowActionTypeN8nWebhookLabel  = "action_n8n_webhook"
	TriggerLowStockMonitorTypeAny      = "any_item"
	TriggerLowStockMonitorTypeSpecific = "specific_items"
)

type WorkflowEnvVar struct {
	Name     string `json:"name" bson:"name" mapstructure:"name"`
	Value    string `json:"value" bson:"value" mapstructure:"value"`
	IsSecret bool   `json:"is_secret" bson:"is_secret" mapstructure:"is_secret"`
}

type Workflow struct {
	ID          string               `json:"id" bson:"id" mapstructure:"id"`
	Name        string               `json:"name" bson:"name" mapstructure:"name"`
	Status      string               `json:"status" bson:"status" mapstructure:"status"`
	Description string               `json:"description" bson:"description" mapstructure:"description"`
	Enabled     bool                 `json:"enabled" bson:"enabled" mapstructure:"enabled"`
	Trigger     WorkflowTriggerBase  `json:"trigger" bson:"trigger" mapstructure:"trigger"`
	Actions     []WorkflowActionBase `json:"actions" bson:"actions" mapstructure:"actions"`
	Runs        []WorkflowRun        `json:"runs" bson:"runs" mapstructure:"runs"`
}

type WorkflowRun struct {
	ID        string           `json:"id" bson:"id" mapstructure:"id"`
	Logs      []WorkflowRunLog `json:"logs" bson:"logs" mapstructure:"logs"`
	StartTime time.Time        `json:"start_time" bson:"start_time" mapstructure:"start_time"`
	EndTime   time.Time        `json:"end_time" bson:"end_time" mapstructure:"end_time"`
	Status    string           `json:"status" bson:"status" mapstructure:"status"`
}

type WorkflowRunLog struct {
	Level     string    `json:"level" bson:"level" mapstructure:"level"`
	TimeStamp time.Time `json:"timestamp" bson:"timestamp" mapstructure:"timestamp"`
	Message   string    `json:"message" bson:"message" mapstructure:"message"`
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
	Output              string   `json:"output" bson:"output" mapstructure:"output"`
}

type WorkflowLowStockTriggerOutput struct {
	Items []WorkflowLowStockTriggerOutputItem `json:"items" bson:"items" mapstructure:"items"`
}

type WorkflowLowStockTriggerOutputItem struct {
	TenantId string   `json:"tenant_id" bson:"tenant_id" mapstructure:"tenant_id"`
	Labels   []string `json:"labels" bson:"labels" mapstructure:"labels"`
	ItemID   string   `json:"item_id" bson:"item_id" mapstructure:"item_id"`
	ItemName string   `json:"item_name" bson:"item_name" mapstructure:"item_name"`
	Quantity float64  `json:"quantity" bson:"quantity" mapstructure:"quantity"`
	Unit     string   `json:"unit" bson:"unit" mapstructure:"unit"`
}

type WorkflowN8nWebhookAction struct {
	WorkflowActionBase `json:",inline" bson:",inline" mapstructure:",squash"`
	Input              string            `json:"input" bson:"input" mapstructure:"input"`
	WebhookURL         string            `json:"webhook_url" bson:"webhook_url" mapstructure:"webhook_url"`
	Method             string            `json:"method" bson:"method" mapstructure:"method"`
	Headers            map[string]string `json:"headers" bson:"headers" mapstructure:"headers"`
	Timeout            int               `json:"timeout" bson:"timeout" mapstructure:"timeout"`
	Output             string            `json:"output" bson:"output" mapstructure:"output"`
}
