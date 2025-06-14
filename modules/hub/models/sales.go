package models

import core_models "github.com/nutrixpos/pos/modules/core/models"

type SalesPerDay struct {
	Id           string                   `json:"id" bson:"id,omitempty" mapstructure:"id"`
	Date         string                   `json:"date" bson:"date" mapstructure:"date"`
	Orders       []SalesPerDayOrder       `json:"orders" bson:"orders" mapstructure:"orders"`
	Refunds      []core_models.ItemRefund `json:"refunds" bson:"refunds" mapstructure:"refunds"`
	Costs        float64                  `json:"costs" bson:"costs" mapstructure:"costs"`
	TotalSales   float64                  `json:"total_sales" bson:"total_sales" mapstructure:"total_sales"`
	RefundsValue float64                  `json:"refunds_value" bson:"refunds_value" mapstructure:"refunds_value"`
}

type SalesPerDayOrder struct {
	core_models.SalesPerDayOrder `json:",inline" bson:",inline" mapstructure:",squash"`
	Labels                       []string `json:"labels" bson:"labels" mapstructure:"labels"`
}
