package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"github.com/nutrixpos/hub/common"
	"github.com/nutrixpos/hub/common/config"
	"github.com/nutrixpos/hub/modules/hub/models"
	"github.com/nutrixpos/hub/modules/hub/services"
	"github.com/nutrixpos/pos/common/logger"
	core_models "github.com/nutrixpos/pos/modules/core/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func LogsPost(config config.Config, logger logger.ILogger) http.HandlerFunc {
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

			if tenant_id == "" {
				http.Error(w, "tenant_id claim is required", http.StatusBadRequest)
				logger.Error("ERROR: tenant_id claim is required")
				return
			}

			label = claims["name"].(string)
			if label == "" {
				http.Error(w, "name claim is required", http.StatusBadRequest)
				logger.Error("ERROR: name claim is required")
				return
			}
		}

		label = fmt.Sprintf("branch:%s", label)

		request_body := struct {
			Data []interface{} `json:"data"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&request_body); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}

		if len(request_body.Data) == 0 {
			http.Error(w, "Logs array is empty", http.StatusBadRequest)
			logger.Error("ERROR: Logs array is empty")
			return
		}

		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
		if err != nil {
			http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}

		defer client.Disconnect(context.TODO())

		if len(request_body.Data) > 0 {

			client_sales_orders := make([]models.SalesPerDayOrder, 0)
			client_sales_refunds := make([]models.LogOrderItemRefund, 0)

			for _, v := range request_body.Data {

				var log core_models.Log
				decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
					DecodeHook: common.StringToTimeHook(),
					Result:     &log,
				})

				err = decoder.Decode(v)

				if err != nil {
					logger.Error(fmt.Sprintf("Failed to decode log: %v", err))
					continue
				}

				switch log.Type {

				case core_models.LogTypeSalesPerDayOrder:

					var db_log models.LogSalesPerDayOrder
					decoder, _ = mapstructure.NewDecoder(&mapstructure.DecoderConfig{
						DecodeHook: common.StringToTimeHook(),
						Result:     &db_log,
					})
					err = decoder.Decode(v)

					if err != nil {
						logger.Error(fmt.Sprintf("Failed to decode log: %v", err))
						continue
					}

					db_log.Labels = []string{}
					db_log.SalesPerDayOrder.Labels = []string{}

					db_log.Labels = append(db_log.Labels, label)
					db_log.SalesPerDayOrder.Labels = append(db_log.SalesPerDayOrder.Labels, label)

					client_sales_orders = append(client_sales_orders, db_log.SalesPerDayOrder)

				case core_models.LogTypeOrderItemRefunded:

					var db_log models.LogOrderItemRefund
					decoder, _ = mapstructure.NewDecoder(&mapstructure.DecoderConfig{
						DecodeHook: common.StringToTimeHook(),
						Result:     &db_log,
					})
					err = decoder.Decode(v)

					if err != nil {
						logger.Error(fmt.Sprintf("Failed to decode log: %v", err))
						continue
					}

					db_log.Labels = []string{}
					db_log.Labels = append(db_log.Labels, label)

					client_sales_refunds = append(client_sales_refunds, db_log)

				default:
					logger.Warning(fmt.Sprintf("Unknown log type: %s", log.Type))
				}
			}

			sales_svc := services.SalesService{
				Logger: logger,
				Config: config,
			}

			if len(client_sales_orders) > 0 {
				err = sales_svc.InsertClientSalesOrders(tenant_id, client_sales_orders)
				if err != nil {
					http.Error(w, "Failed to insert logs", http.StatusInternalServerError)
					logger.Error(fmt.Sprintf("ERROR: %v", err))
					return
				}
			}

			if len(client_sales_refunds) > 0 {
				err = sales_svc.InsertClientSalesRefunds(tenant_id, client_sales_refunds)
				if err != nil {
					http.Error(w, "Failed to insert logs", http.StatusInternalServerError)
					logger.Error(fmt.Sprintf("ERROR: %v", err))
					return
				}
			}

		}

		w.WriteHeader(http.StatusCreated)
	}
}
