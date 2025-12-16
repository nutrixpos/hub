package hub

import (
	"github.com/gorilla/mux"
	"github.com/nutrixpos/hub/common"
	"github.com/nutrixpos/hub/common/config"
	"github.com/nutrixpos/hub/modules"
	"github.com/nutrixpos/hub/modules/hub/events"
	"github.com/nutrixpos/hub/modules/hub/handlers"
	"github.com/nutrixpos/hub/modules/hub/models"
	"github.com/nutrixpos/hub/modules/hub/services"
	"github.com/nutrixpos/pos/common/logger"
	pos_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
)

type HubModule struct {
	Name            string
	Config          config.Config
	Logger          logger.ILogger
	Settings        models.Settings
	EventBus        common.EventBus
	EventHandlersId []string
}

func (h *HubModule) SetName(name string) error {
	h.Name = name
	return nil
}

func (h *HubModule) GetName() string {
	return h.Name
}

// OnStart is called when the core module is started.
func (h *HubModule) OnStart() func() error {
	return func() error {
		err := h.EnsureSeeded()
		if err != nil {
			return err
		}

		return nil
	}
}

// OnEnd is called when the core module is ended.
func (h *HubModule) OnEnd() func() {
	return func() {
		for _, handler_id := range h.EventHandlersId {
			h.EventBus.UnregisterHandler(events.EventLowStockId, handler_id)
		}
	}
}

func (h *HubModule) RegisterHttpHandlers(router *mux.Router, prefix string) {
	router.Handle("/v1/api/logs", pos_middlewares.AllowCors(handlers.LogsPost(h.Config, h.Logger))).Methods("POST", "OPTIONS")
	router.Handle("/v1/api/inventories", pos_middlewares.AllowCors(handlers.InventoryItemsPut(h.Config, h.Logger))).Methods("PUT", "OPTIONS")
	router.Handle("/v1/api/inventories", pos_middlewares.AllowCors(handlers.InventoryItemsGet(h.Config, h.Logger))).Methods("GET", "OPTIONS")
	router.Handle("/v1/api/inventories", pos_middlewares.AllowCors(handlers.InventoryItemsPatch(h.Config, h.Logger))).Methods("PATCH", "OPTIONS")
	router.Handle("/v1/api/sales", pos_middlewares.AllowCors(handlers.GetSalesPerDay(h.Config, h.Logger))).Methods("GET", "OPTIONS")
	router.Handle("/v1/api/workflows", pos_middlewares.AllowCors(handlers.WorkflowsGET(h.Config, h.Logger))).Methods("GET", "OPTIONS")
	router.Handle("/v1/api/workflows/{id}", pos_middlewares.AllowCors(handlers.WorkflowGET(h.Config, h.Logger))).Methods("GET", "OPTIONS")
	router.Handle("/v1/api/workflows", pos_middlewares.AllowCors(handlers.WorkflowPOST(h.Config, h.Logger))).Methods("POST", "OPTIONS")
	router.Handle("/v1/api/workflows/{id}", pos_middlewares.AllowCors(handlers.WorkflowPATCH(h.Config, h.Logger))).Methods("PATCH", "OPTIONS")
	router.Handle("/v1/api/languages", pos_middlewares.AllowCors(handlers.GetAvailableLanguages(h.Config, h.Logger))).Methods("GET", "OPTIONS")
	router.Handle("/v1/api/languages/{code}", pos_middlewares.AllowCors(handlers.GetLanguage(h.Config, h.Logger))).Methods("GET", "OPTIONS")
	router.Handle("/v1/api/settings", pos_middlewares.AllowCors(handlers.GetSettings(h.Config, h.Logger))).Methods("GET", "OPTIONS")
	router.Handle("/v1/api/settings", pos_middlewares.AllowCors(handlers.UpdateSettings(h.Config, h.Logger))).Methods("PATCH", "OPTIONS")
	router.Handle("/v1/api/koptan/suggestions", pos_middlewares.AllowCors(handlers.GetKoptanSuggestions(h.Config, h.Logger))).Methods("GET", "OPTIONS")
	router.Handle("/v1/api/subscriptions", pos_middlewares.AllowCors(handlers.SubcriptionGET(h.Config, h.Logger))).Methods("GET", "OPTIONS")
	router.Handle("/v1/api/subscriptions/request", pos_middlewares.AllowCors(handlers.SubcriptionRequest(h.Config, h.Logger))).Methods("POST", "OPTIONS")
	router.Handle("/v1/api/subscriptions/payment_callback", pos_middlewares.AllowCors(handlers.PaymobSubscribeCallbackPOST(h.Config, h.Logger))).Methods("POST", "OPTIONS")
	router.Handle("/v1/api/subscriptions/request_cancellation", pos_middlewares.AllowCors(handlers.SubscriptionRequestCancellation(h.Config, h.Logger))).Methods("POST", "OPTIONS")
}

func (h *HubModule) RegisterBackgroundWorkers() []modules.Worker {
	return nil
}

func (h *HubModule) RegisterEventBus(eb common.EventBus) error {
	h.EventBus = eb
	return nil
}

func (h *HubModule) EnsureSeeded() error {

	seeder_svc := services.SeederService{
		Config: &h.Config,
		Logger: h.Logger,
	}

	err := seeder_svc.Seed()
	if err != nil {
		return err
	}

	return nil
}
