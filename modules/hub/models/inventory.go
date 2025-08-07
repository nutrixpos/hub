package models

type InventoryItem struct {
	ID       string                `json:"id" bson:"id"`
	TenantID string                `json:"tenant_id" bson:"tenant_id"`
	Name     string                `json:"name" bson:"name"`
	Quantity float64               `json:"quantity" bson:"quantity"`
	Unit     string                `json:"unit" bson:"unit"`
	Labels   []string              `json:"labels" bson:"labels"`     // Labels for categorization
	Settings InventoryItemSettings `json:"settings" bson:"settings"` // Additional settings for the inventory item
}

type InventoryItemSettings struct {
	AlertThreshold float64 `json:"alert_threshold" bson:"alert_threshold"` // Threshold for alerting when stock is low
	AlertEnabled   bool    `json:"alert_enabled" bson:"alert_enabled"`     // Whether alerts are enabled for this item
}
