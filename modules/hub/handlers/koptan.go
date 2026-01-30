package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/nutrixpos/hub/common/config"
	"github.com/nutrixpos/hub/modules/hub/services"
	"github.com/nutrixpos/pos/common/logger"
	core_handlers "github.com/nutrixpos/pos/modules/core/handlers"
)

type LLMRequest struct {
	MetaData map[string]string `json:"metadata"`
	Messages []LLMMessage      `json:"messages"`
}

type LLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func KoptanChat(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
		defer cancel()

		tenant_id := "1"
		label := "dev"

		if config.Env != "dev" {
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

			var claims map[string]interface{}
			err = json.Unmarshal(decodedData, &claims)
			if err != nil {
				http.Error(w, "Failed to unmarshal token", http.StatusBadRequest)
				logger.Error(fmt.Sprintf("ERROR: %v", err))
				return
			}

			var ok bool
			tenant_id, ok = claims["tenant_id"].(string)
			if !ok || tenant_id == "" {
				http.Error(w, "tenant_id claim is required and must be a string", http.StatusBadRequest)
				logger.Error("ERROR: tenant_id claim is required and must be a string")
				return
			}

			label = claims["name"].(string)
			if label == "" {
				http.Error(w, "name claim is required", http.StatusBadRequest)
				logger.Error("ERROR: name claim is required")
				return
			}
		}

		request := struct {
			Data struct {
				Messages []LLMMessage `json:"messages"`
			} `json:"data"`
		}{}

		requestBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Error(err.Error())
			return
		}
		r.Body.Close()

		if err := json.Unmarshal(requestBytes, &request); err != nil {
			logger.Error(fmt.Sprintf("Failed to parse JSON: %v", err))
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		if len(request.Data.Messages) == 0 {
			http.Error(w, "messages are required", http.StatusBadRequest)
			return
		}

		if err != nil {
			logger.Error(fmt.Sprintf("Failed to decode data: %v", err))
			http.Error(w, "Failed to decode data", http.StatusBadRequest)
			return
		}

		messages := []LLMMessage{
			{
				Role:    "system",
				Content: "You are a helpful restaurant manager, and a sales expert, and you are here to help the owner of the restaurant to boost his sales and manage his inventory.  You have set of tools that can help you get data from the restaurant database, you can run them then, analyse their output without exposing their output to the user, and return back the final answer to the user, make the output very concise and to the point.",
			},
		}

		for _, msg := range request.Data.Messages {
			messages = append(messages, LLMMessage{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}

		llmReq := LLMRequest{
			MetaData: map[string]string{
				"tenant_id": tenant_id,
				"label":     label,
			},
			Messages: messages,
		}

		LLMRequestJSON, err := json.Marshal(struct {
			Data LLMRequest `json:"data"`
		}{
			Data: llmReq,
		})

		if err != nil {
			logger.Error(fmt.Sprintf("Failed to encode LLMRequest: %v", err))
			http.Error(w, "Failed to encode LLMRequest", http.StatusInternalServerError)
			return
		}

		url := "http://localhost:8000"
		req, err := http.NewRequest("POST", url, bytes.NewBufferString(string(LLMRequestJSON)))

		if err != nil {
			logger.Error(fmt.Sprintf("Failed to create POST request: %v to koptan service", err))
			http.Error(w, fmt.Sprintf("Failed to create POST request to koptan service"), http.StatusInternalServerError)
			return
		}

		resp, err := http.DefaultClient.Do(req.WithContext(ctx))
		defer resp.Body.Close()
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to send POST request: %v to koptan service", err))
			http.Error(w, fmt.Sprintf("Failed to send POST request to koptan service"), http.StatusInternalServerError)
			return
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to read upstream response: %v", err))
			http.Error(w, "Failed to read response from LLM service", http.StatusInternalServerError)
			return
		}

		response := core_handlers.JSONApiOkResponse{
			Data: string(bodyBytes),
		}

		data, err := json.Marshal(response)
		if err != nil {
			// No headers have been sent yet, so we can return a 500
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		n, err := w.Write(data)
		if err != nil {
			// Note: If headers were already sent, this might only be useful for logging.
			logger.Error(fmt.Sprintf("Write failed: %v (bytes written: %d)", err, n))
			return
		}

	}
}

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
