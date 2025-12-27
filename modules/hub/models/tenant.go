package models

import (
	"time"
)

type TenantSubscription struct {
	ID               string    `json:"id" bson:"id"`
	SubscriptionPlan string    `json:"subscription_plan" bson:"subscription_plan"`
	StartDate        time.Time `json:"start_date" bson:"start_date"`
	EndDate          time.Time `json:"end_date" bson:"end_date"`
	Status           string    `json:"status" bson:"status"`
}

type Tenant struct {
	ID             string             `json:"id" bson:"id" mapstructure:"id"`
	TenantID       string             `bson:"tenant_id" json:"tenant_id" mapstructure:"tenant_id"`
	Sales          []SalesPerDay      `bson:"sales" json:"sales" mapstructure:"sales"`
	InventoryItems []InventoryItem    `bson:"inventory_items" json:"inventory_items" mapstructure:"inventory_items"`
	Subscription   TenantSubscription `bson:"subscription" json:"subscription" mapstructure:"subscription"`
	Workflows      []interface{}      `json:"workflows" bson:"workflows" mapstructure:"workflows"`
}

type TenantAPIKey struct {
	ID             string    `json:"id" bson:"id"`
	APIKey         string    `bson:"api_key"`
	Title          string    `bson:"title"`
	ExpirationDate time.Time `bson:"expiration_date"`
}
