package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"github.com/nutrixpos/hub/common"
	"github.com/nutrixpos/hub/services"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	core_handlers "github.com/nutrixpos/pos/modules/core/handlers"
	core_models "github.com/nutrixpos/pos/modules/core/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetSalesPerDay(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		page_number, err := strconv.Atoi(r.URL.Query().Get("page[number]"))
		if err != nil || page_number == 0 {
			page_number = 1
		}

		page_size, err := strconv.Atoi(r.URL.Query().Get("page[size]"))
		if err != nil {
			page_size = 50
		}

		tenant_id := r.URL.Query().Get("tenant_id")
		if tenant_id == "" {
			http.Error(w, "tenant_id query string is required", http.StatusBadRequest)
			return
		}

		salesService := services.SalesService{
			Logger: logger,
			Config: config,
		}

		sales, totalRecords, err := salesService.GetSalesPerday(page_number, page_size, tenant_id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := core_handlers.JSONApiOkResponse{
			Meta: core_handlers.JSONAPIMeta{
				TotalRecords: totalRecords,
			},
			Data: sales,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func LogsPost(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenant_id := r.URL.Query().Get("tenant_id")
		if tenant_id == "" {
			http.Error(w, "tenant_id query string is required", http.StatusBadRequest)
			return
		}

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

			client_sales_orders := make([]core_models.SalesPerDayOrder, 0)
			client_sales_refunds := make([]core_models.LogOrderItemRefund, 0)

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

					var db_log core_models.LogSalesPerDayOrder
					decoder, _ = mapstructure.NewDecoder(&mapstructure.DecoderConfig{
						DecodeHook: common.StringToTimeHook(),
						Result:     &db_log,
					})
					err = decoder.Decode(v)

					if err != nil {
						logger.Error(fmt.Sprintf("Failed to decode log: %v", err))
						continue
					}

					client_sales_orders = append(client_sales_orders, db_log.SalesPerDayOrder)

				case core_models.LogTypeOrderItemRefunded:

					var db_log core_models.LogOrderItemRefund
					decoder, _ = mapstructure.NewDecoder(&mapstructure.DecoderConfig{
						DecodeHook: common.StringToTimeHook(),
						Result:     &db_log,
					})
					err = decoder.Decode(v)

					if err != nil {
						logger.Error(fmt.Sprintf("Failed to decode log: %v", err))
						continue
					}

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
