package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nutrixpos/hub/common/config"
	"github.com/nutrixpos/hub/modules/hub/models"
	"github.com/nutrixpos/pos/common/logger"
	core_handlers "github.com/nutrixpos/pos/modules/core/handlers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SubscribeResponse struct {
	PaymentKeys      []PaymentKey    `json:"payment_keys"`
	IntentionOrderID int             `json:"intention_order_id"`
	ID               string          `json:"id"`
	IntentionDetail  IntentionDetail `json:"intention_detail"`
	ClientSecret     string          `json:"client_secret"`
	PaymentMethods   []PaymentMethod `json:"payment_methods"`
	SpecialReference string          `json:"special_reference"`
	Extras           Extras          `json:"extras"`
	Confirmed        bool            `json:"confirmed"`
	Status           string          `json:"status"`
	Created          time.Time       `json:"created"`
	CardDetail       CardDetail      `json:"card_detail"`
	CardTokens       []CardToken     `json:"card_tokens"`
	Object           string          `json:"object"`
}

type PaymentKey struct {
	Integration    int    `json:"integration"`
	Key            string `json:"key"`
	GatewayType    string `json:"gateway_type"`
	IframeID       string `json:"iframe_id"`
	OrderID        int    `json:"order_id"`
	RedirectionURL string `json:"redirection_url"`
	SaveCard       bool   `json:"save_card"`
}

type IntentionDetail struct {
	Amount      int         `json:"amount"`
	Items       []Item      `json:"items"`
	Currency    string      `json:"currency"`
	BillingData BillingData `json:"billing_data"`
}

type Item struct {
	Name        string `json:"name"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	Image       string `json:"image"`
}

type BillingData struct {
	Apartment      string `json:"apartment"`
	Floor          string `json:"floor"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Street         string `json:"street"`
	Building       string `json:"building"`
	PhoneNumber    string `json:"phone_number"`
	ShippingMethod string `json:"shipping_method"`
	City           string `json:"city"`
	Country        string `json:"country"`
	State          string `json:"state"`
	Email          string `json:"email"`
	PostalCode     string `json:"postal_code"`
}

type PaymentMethod struct {
	IntegrationID  int    `json:"integration_id"`
	Alias          string `json:"alias"`
	Name           string `json:"name"`
	MethodType     string `json:"method_type"`
	Currency       string `json:"currency"`
	Live           bool   `json:"live"`
	UseCVCWithMoto bool   `json:"use_cvc_with_moto"`
}

type Extras struct {
	CreationExtras     map[string]interface{} `json:"creation_extras"`
	ConfirmationExtras map[string]interface{} `json:"confirmation_extras"`
}

type CardDetail struct {
	CardNumber string `json:"card_number"`
	Expiry     string `json:"expiry"`
	CVV        string `json:"cvv"`
}

type CardToken struct {
	Integration int    `json:"integration"`
	Token       string `json:"token"`
	CardNumber  string `json:"card_number"`
	Expiry      string `json:"expiry"`
	CVV         string `json:"cvv"`
}

type SubscribeRequest struct {
	Amount             int    `json:"amount"`
	Currency           string `json:"currency"`
	PaymentMethods     []int  `json:"payment_methods"`
	SubscriptionPlanID int    `json:"subscription_plan_id"`
	Items              []struct {
		Name     string `json:"name"`
		Amount   int    `json:"amount"`
		Quantity int    `json:"quantity"`
	} `json:"items"`
	BillingData struct {
		Apartment   string `json:"apartment"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Street      string `json:"street"`
		Building    string `json:"building"`
		PhoneNumber string `json:"phone_number"`
		Country     string `json:"country"`
		Email       string `json:"email"`
		Floor       string `json:"floor"`
		State       string `json:"state"`
	} `json:"billing_data"`
	Customer struct {
		FirstName string                 `json:"first_name"`
		LastName  string                 `json:"last_name"`
		Email     string                 `json:"email"`
		Extras    map[string]interface{} `json:"extras"`
	} `json:"customer"`
	Extras map[string]interface{} `json:"extras"`
}

func SubcriptionSubscribe(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func SubcriptionRequest(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenant_id := "1"
		client_name := "Dev"
		client_email := "dev@dev.dev"
		given_name := "Dev"
		family_name := "Dev"
		phone := "1234567890"

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
		}

		request := struct {
			Data struct {
				Plan string `json:"plan"`
			} `json:"data"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Set up MongoDB connection
		clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", config.Databases[0].Host, config.Databases[0].Port))
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}
		defer client.Disconnect(ctx)

		// Access the collection
		collection := client.Database(config.Databases[0].Database).Collection(config.Databases[0].Tables["sales"])

		// Define the filter
		filter := bson.M{
			"tenant_id": tenant_id,
		}

		// Find the document
		var existing_tenant models.Tenant
		err = collection.FindOne(ctx, filter).Decode(&existing_tenant)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "No document found with the specified tenant_id", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to retrieve document", http.StatusInternalServerError)
				logger.Error(fmt.Sprintf("ERROR: %v", err))
			}
			return
		}

		new_subscription := existing_tenant.Subscription

		switch request.Data.Plan {
		case "standard":
			if existing_tenant.Subscription.SubcriptionPlan == "free" {
				subscription_request := SubscribeRequest{
					Amount:             1500, // Example amount in cents
					Currency:           "EGP",
					PaymentMethods:     []int{5229518},
					SubscriptionPlanID: 4088, // Example plan ID
					Items: []struct {
						Name     string `json:"name"`
						Amount   int    `json:"amount"`
						Quantity int    `json:"quantity"`
					}{{
						Name:     "Standard subscription",
						Amount:   1500,
						Quantity: 1,
					}},
					BillingData: struct {
						Apartment   string `json:"apartment"`
						FirstName   string `json:"first_name"`
						LastName    string `json:"last_name"`
						Street      string `json:"street"`
						Building    string `json:"building"`
						PhoneNumber string `json:"phone_number"`
						Country     string `json:"country"`
						Email       string `json:"email"`
						Floor       string `json:"floor"`
						State       string `json:"state"`
					}{
						Apartment:   "",
						FirstName:   client_name,
						LastName:    "",
						Street:      "",
						Building:    "",
						PhoneNumber: "",
						Country:     "",
						Email:       client_email,
						Floor:       "",
						State:       "",
					},
					Customer: struct {
						FirstName string                 `json:"first_name"`
						LastName  string                 `json:"last_name"`
						Email     string                 `json:"email"`
						Extras    map[string]interface{} `json:"extras"`
					}{
						FirstName: client_name,
						LastName:  "",
						Email:     client_email,
						Extras: map[string]interface{}{
							"tenant_id": tenant_id,
						},
					},
					Extras: map[string]interface{}{
						"tenant_id": tenant_id,
					},
				}

				json_data, err := json.Marshal(subscription_request)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				req, err := http.NewRequest("POST", config.Payment.SubscribingURL, bytes.NewBuffer(json_data))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", fmt.Sprintf("Token %s", config.Payment.SecretKey))

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					var errorResponse map[string]interface{}
					err = json.NewDecoder(resp.Body).Decode(&errorResponse)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						logger.Error(fmt.Sprintf("Failed to create subscription, unable to decode body:: %v", err))
						return
					}

					http.Error(w, fmt.Sprintf("Failed to create subscription, error: %v", errorResponse), http.StatusInternalServerError)
					logger.Error(fmt.Sprintf("Failed to create subscription, error:: %v", errorResponse))
					return
				}

				var subscribeResponse SubscribeResponse
				err = json.NewDecoder(resp.Body).Decode(&subscribeResponse)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				response := core_handlers.JSONApiOkResponse{
					Data: struct {
						SubscriptionPlan string `json:"subscription_id"`
						ClientSecret     string `json:"client_secret"`
						PublicKey        string `json:"public_key"`
					}{
						SubscriptionPlan: "standard",
						ClientSecret:     subscribeResponse.ClientSecret,
						PublicKey:        config.Payment.PublicKey,
					},
				}

				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(response); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		existing_tenant.Subscription = new_subscription

		// Update the document
		_, err = collection.ReplaceOne(ctx, filter, existing_tenant)
		if err != nil {
			http.Error(w, "Failed to update document", http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func SubcriptionGET(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tenant_id := "1"

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
		}

		deadline := 5 * time.Second
		if config.Env == "dev" {
			deadline = 1000 * time.Second
		}

		// Set up MongoDB connection
		clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", config.Databases[0].Host, config.Databases[0].Port))
		ctx, cancel := context.WithTimeout(context.Background(), deadline)
		defer cancel()

		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}
		defer client.Disconnect(ctx)

		tenant_col := client.Database("nutrixhub").Collection(config.Databases[0].Tables["sales"])
		filter := bson.M{"tenant_id": tenant_id}
		var tenant models.Tenant
		err = tenant_col.FindOne(ctx, filter).Decode(&tenant)
		if err != nil {
			http.Error(w, "Failed to get tenant", http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}

		if tenant.Subscription.SubcriptionPlan == "" {
			tenant.Subscription = models.TenantSubscription{
				ID:              primitive.NewObjectID().Hex(),
				SubcriptionPlan: "free",
				StartDate:       time.Now(),
				EndDate:         time.Now().AddDate(99, 0, 0), // 99 years from now
				Status:          "active",
			}
		}

		_, err = tenant_col.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"subscription": tenant.Subscription}})
		if err != nil {
			http.Error(w, "Failed to update tenant subscription", http.StatusInternalServerError)
			logger.Error(fmt.Sprintf("ERROR: %v", err))
			return
		}

		response := core_handlers.JSONApiOkResponse{
			Data: tenant.Subscription,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
