package events

const (
	EventLowStockId = "event_low_stock"
)

type EventLowStockData struct {
	ItemID    string
	ItemName  string
	Threshold int
	Current   int
}
