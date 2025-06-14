package models

import core_models "github.com/nutrixpos/pos/modules/core/models"

type LogSalesPerDayOrder struct {
	core_models.Log  `json:",inline" bson:",inline" mapstructure:",squash"`
	SalesPerDayOrder SalesPerDayOrder `json:"sales_per_day_order" bson:"sales_per_day_order" mapstructure:"sales_per_day_order"`
	Labels           []string         `json:"labels" bson:"labels" mapstructure:"labels"`
}

type LogOrderItemRefund struct {
	core_models.LogOrderItemRefund `json:",inline" bson:",inline" mapstructure:",squash"`
	Labels                         []string `json:"labels" bson:"labels" mapstructure:"labels"`
}
