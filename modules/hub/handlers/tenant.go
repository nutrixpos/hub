package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nutrixpos/hub/modules/hub/models"
	"github.com/nutrixpos/hub/modules/hub/services"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
)

// GetSettings is an http get handlers that just returns the settings object from the db
func PostAPIKey(conf config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		tenant_id := r.URL.Query().Get("tenant_id")
		if tenant_id == "" {
			http.Error(w, "tenant_id query string is required", http.StatusBadRequest)
			return
		}

		tenant_svc := services.TenantService{
			Config: conf,
		}

		request := struct {
			Data models.TenantAPIKey `json:"data"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = tenant_svc.AddAPIKey(tenant_id, request.Data.Title, request.Data.ExpirationDate)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
