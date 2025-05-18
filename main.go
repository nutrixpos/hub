package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nutrixpos/hub/handlers"
	"github.com/nutrixpos/pos/common/config"
	"github.com/nutrixpos/pos/common/logger"
	pos_middlewares "github.com/nutrixpos/pos/modules/core/middlewares"
)

func main() {

	logger := logger.NewZeroLog()

	// Create the configuration using the Viper config backend
	conf := config.ConfigFactory("viper", "config.yaml", &logger)

	r := mux.NewRouter()
	r.Handle("/v1/api/logs", handlers.LogsPost(conf, &logger)).Methods("POST", "OPTIONS")
	r.Handle("/v1/api/sales", pos_middlewares.AllowCors(handlers.GetSalesPerDay(conf, &logger))).Methods("GET", "OPTIONS")

	logger.Info("Serving on port 8001")
	err := http.ListenAndServe(":8001", r)
	if err != nil {
		panic(err)
	}
}
