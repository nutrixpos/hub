package models

import (
	"time"
)

type Tenant struct {
	ID       string        `json:"id" bson:"id"`
	TenantID string        `bson:"tenant_id"`
	Sales    []SalesPerDay `bson:"sales"`
}

type TenantAPIKey struct {
	ID             string    `json:"id" bson:"id"`
	APIKey         string    `bson:"api_key"`
	Title          string    `bson:"title"`
	ExpirationDate time.Time `bson:"expiration_date"`
}
