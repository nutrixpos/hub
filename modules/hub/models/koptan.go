package models

import "time"

type KoptanSuggestion struct {
	ID        int       `bson:"id" json:"id"`
	TenantID  string    `bson:"tenant_id" json:"tenant_id"`
	Type      string    `bson:"type" json:"type"`
	Content   string    `bson:"content" json:"content"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
