package events

const (
	EventLowStockId = "event_low_stock"
)

type EventLowStockData struct {
	TenantId  string
	ItemID    string
	ItemName  string
	Threshold float64
	Current   float64
}
