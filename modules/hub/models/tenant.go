package models

import (
	"time"
)

type TenantSubscription struct {
	ID              string    `json:"id" bson:"id"`
	SubcriptionPlan string    `bson:"subscription_plan"`
	StartDate       time.Time `bson:"start_date"`
	EndDate         time.Time `bson:"end_date"`
	Status          string    `bson:"status"`
}

type Tenant struct {
	ID             string             `json:"id" bson:"id"`
	TenantID       string             `bson:"tenant_id"`
	Sales          []SalesPerDay      `bson:"sales"`
	InventoryItems []InventoryItem    `bson:"inventory_items"`
	Subscription   TenantSubscription `bson:"subscription"`
}

type TenantAPIKey struct {
	ID             string    `json:"id" bson:"id"`
	APIKey         string    `bson:"api_key"`
	Title          string    `bson:"title"`
	ExpirationDate time.Time `bson:"expiration_date"`
}
