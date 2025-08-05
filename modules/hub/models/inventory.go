package models

type InventoryItem struct {
	ID       string                `json:"id" bson:"id"`
	TenantID string                `bson:"tenant_id"`
	Name     string                `bson:"name"`
	Quantity float64               `bson:"quantity"`
	Labels   []string              `bson:"labels"`   // Labels for categorization
	Settings InventoryItemSettings `bson:"settings"` // Additional settings for the inventory item
}

type InventoryItemSettings struct {
	AlertThreshold float64 `bson:"alert_threshold"` // Threshold for alerting when stock is low
	AlertEnabled   bool    `bson:"alert_enabled"`   // Whether alerts are enabled for this item
}
