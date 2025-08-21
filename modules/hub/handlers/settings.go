package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nutrixpos/hub/common/config"
	"github.com/nutrixpos/hub/modules/hub/models"
	"github.com/nutrixpos/hub/modules/hub/services"
	"github.com/nutrixpos/pos/common/logger"
	core_handlers "github.com/nutrixpos/pos/modules/core/handlers"
)

// UpdateSettings is a post request handler that updates the settings in the database
// send a models.Settings directory to body to use it.
func UpdateSettings(conf config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		request := struct {
			Data models.Settings `json:"data"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		settings_svc := services.SettingsService{
			Config: conf,
		}

		err = settings_svc.UpdateSettings(request.Data)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	}
}

// GetSettings is an http get handlers that just returns the settings object from the db
func GetSettings(conf config.Config, logger logger.ILogger) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		settings_svc := services.SettingsService{
			Config: conf,
		}

		settings, err := settings_svc.GetSettings()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := core_handlers.JSONApiOkResponse{
			Data: settings,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "Failed to marshal order settings response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}

}
