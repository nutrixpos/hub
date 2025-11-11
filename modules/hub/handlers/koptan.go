package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/nutrixpos/hub/common/config"
	"github.com/nutrixpos/hub/modules/hub/services"
	"github.com/nutrixpos/pos/common/logger"
	core_handlers "github.com/nutrixpos/pos/modules/core/handlers"
)

func GetKoptanSuggestions(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		page_number, err := strconv.Atoi(r.URL.Query().Get("page[number]"))
		if err != nil || page_number == 0 {
			page_number = 1
		}

		page_size, err := strconv.Atoi(r.URL.Query().Get("page[size]"))
		if err != nil {
			page_size = 50
		}

		tenant_id := "1"

		if config.Env == "prod" {
			token := r.Header.Get("X-Userinfo")
			if token == "" {
				http.Error(w, "X-Userinfo header is required", http.StatusBadRequest)
				return
			}

			decodedData, err := base64.StdEncoding.DecodeString(token)
			if err != nil {
				http.Error(w, "Failed to decode token", http.StatusBadRequest)
				logger.Error(fmt.Sprintf("ERROR: %v", err))
				return
			}

			// Create a map to hold the decoded JSON data
			var claims map[string]interface{}

			// Unmarshal the decoded data into the map
			err = json.Unmarshal(decodedData, &claims)
			if err != nil {
				http.Error(w, "Failed to unmarshal token", http.StatusBadRequest)
				logger.Error(fmt.Sprintf("ERROR: %v", err))
				return
			}

			var ok bool

			tenant_id, ok = claims["tenant_id"].(string)
			if !ok {
				http.Error(w, "tenant_id claim is required and must be a string", http.StatusBadRequest)
				logger.Error("ERROR: tenant_id claim is required and must be a string")
				return
			}
		}

		tenant_svc := services.TenantService{
			Config: config,
			Logger: logger,
		}

		tenant, err := tenant_svc.GetTenantById(tenant_id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}

		if tenant.Subscription.SubscriptionPlan == "free" {
			http.Error(w, "Upgrade to GOLD to unlock this feature", http.StatusForbidden)
			return
		}

		koptan_svc := services.KoptanService{
			Logger: logger,
			Config: config,
		}

		suggestions, totalRecords, err := koptan_svc.GetSuggestions(tenant_id, page_number, page_size)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}

		response := core_handlers.JSONApiOkResponse{
			Meta: core_handlers.JSONAPIMeta{
				TotalRecords: totalRecords,
			},
			Data: suggestions,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
