package models

import (
	pos_core_models "github.com/nutrixpos/pos/modules/core/models"
)

type Tenant struct {
	ID       string                        `json:"id" bson:"id"`
	TenantID string                        `bson:"tenant_id"`
	Sales    []pos_core_models.SalesPerDay `bson:"sales"`
}
