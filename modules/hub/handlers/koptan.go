package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/nutrixpos/hub/common/config"
	"github.com/nutrixpos/hub/modules/hub/services"
	"github.com/nutrixpos/pos/common/logger"
	core_handlers "github.com/nutrixpos/pos/modules/core/handlers"
	"github.com/teilomillet/gollm"
)

func KoptanChat(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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
			Data string `json:"data"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if request.Data == "" {
			http.Error(w, "data is required", http.StatusBadRequest)
			return
		}

		if err != nil {
			logger.Error(fmt.Sprintf("Failed to decode data: %v", err))
			http.Error(w, "Failed to decode data", http.StatusBadRequest)
			return
		}

		sales_svc := services.SalesService{
			Logger: logger,
			Config: config,
		}

		sales, _, err := sales_svc.GetSalesPerday(1, 1, tenant_id, -1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error(err.Error())
			return
		}

		sales_json, err := json.Marshal(sales)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error(err.Error())
			return
		}

		// Using OpenAI
		llm, err := gollm.NewLLM(
			gollm.SetProvider("ollama"),
			gollm.SetModel("llama3.2:1b"),
		)

		if err != nil {
			logger.Error(fmt.Sprintf("Failed to create LLM: %v", err))
			http.Error(w, fmt.Sprintf("Failed to create LLM: %v", err), http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Minute)
		defer cancel()

		prompt := gollm.NewPrompt(
			fmt.Sprintf("Context Data:\n%s\n\nUser Question: %s", sales_json, request.Data),
			gollm.WithSystemPrompt("You are a senior business analyst. Your job is to look at sales and inventory data and provide actionable insights in natural language.", gollm.CacheTypeEphemeral),
			gollm.WithDirectives(
				"Identify the most urgent problem (e.g., low stock of high-sellers).",
				"Highlight underperforming items.",
				"Do not use tables or code; speak in a helpful, conversational tone.",
				"Keep the response under 150 words.",
			),
		)

		// Generate a llm_response
		llm_response, err := llm.Generate(ctx, prompt)
		if err != nil {
			log.Fatalf("Failed to generate text: %v", err)
		}

		response := core_handlers.JSONApiOkResponse{
			Data: llm_response,
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
