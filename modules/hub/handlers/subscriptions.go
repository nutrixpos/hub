package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

type PaymobTime time.Time

func (t PaymobTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2025-08-22T14:08:47.177486"))
	return []byte(stamp), nil
}

func (t *PaymobTime) UnmarshalJSON(b []byte) error {
	s := string(b)

	// Handle null values
	if s == "null" {
		return nil
	}

	// Remove quotes from JSON string
	s = strings.Trim(s, `"`)

	// Parse the timestamp using the correct format
	parsedTime, err := time.Parse("2006-01-02T15:04:05.999999", s)
	if err != nil {
		// Try alternative formats if needed
		return fmt.Errorf("failed to parse time %s: %v", s, err)
	}

	*t = PaymobTime(parsedTime)
	return nil
}

type PaymobSubscribeResponse struct {
	PaymentKeys      []PaymobPaymentKey          `json:"payment_keys"`
	IntentionOrderID int                         `json:"intention_order_id"`
	ID               string                      `json:"id"`
	IntentionDetail  PaymobIntentionDetail       `json:"intention_detail"`
	ClientSecret     string                      `json:"client_secret"`
	PaymentMethods   []PaymobPaymentMethod       `json:"payment_methods"`
	SpecialReference string                      `json:"special_reference"`
	Extras           PaymobSubscrbeResponseExtra `json:"extras"`
	Confirmed        bool                        `json:"confirmed"`
	Status           string                      `json:"status"`
	Created          PaymobTime                  `json:"created"`
	CardDetail       PaymobCardDetail            `json:"card_detail"`
	CardTokens       []PaymobCardToken           `json:"card_tokens"`
	Object           string                      `json:"object"`
}

type PaymobPaymentKey struct {
	Integration    int    `json:"integration"`
	Key            string `json:"key"`
	GatewayType    string `json:"gateway_type"`
	IframeID       string `json:"iframe_id"`
	OrderID        int    `json:"order_id"`
	RedirectionURL string `json:"redirection_url"`
	SaveCard       bool   `json:"save_card"`
}

type PaymobIntentionDetail struct {
	Amount      int               `json:"amount"`
	Items       []PaymobItem      `json:"items"`
	Currency    string            `json:"currency"`
	BillingData PaymobBillingData `json:"billing_data"`
}

type PaymobBillingData struct {
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

type PaymobPaymentMethod struct {
	IntegrationID  int    `json:"integration_id"`
	Alias          string `json:"alias"`
	Name           string `json:"name"`
	MethodType     string `json:"method_type"`
	Currency       string `json:"currency"`
	Live           bool   `json:"live"`
	UseCVCWithMoto bool   `json:"use_cvc_with_moto"`
}

type PaymobCardDetail struct {
	CardNumber string `json:"card_number"`
	Expiry     string `json:"expiry"`
	CVV        string `json:"cvv"`
}

type PaymobCardToken struct {
	Integration int    `json:"integration"`
	Token       string `json:"token"`
	CardNumber  string `json:"card_number"`
	Expiry      string `json:"expiry"`
	CVV         string `json:"cvv"`
}

type PaymobSubscribeRequest struct {
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

type PaymobSubscribePaymentCallback struct {
	Type                                  string            `json:"type"`
	Obj                                   PaymobTransaction `json:"obj"`
	AcceptFees                            int               `json:"accept_fees"`
	IssuerBank                            interface{}       `json:"issuer_bank"`
	TransactionProcessedCallbackResponses string            `json:"transaction_processed_callback_responses"`
}

type PaymobTransaction struct {
	ID                                    int64                  `json:"id"`
	Pending                               bool                   `json:"pending"`
	AmountCents                           int64                  `json:"amount_cents"`
	Success                               bool                   `json:"success"`
	IsAuth                                bool                   `json:"is_auth"`
	IsCapture                             bool                   `json:"is_capture"`
	IsStandalonePayment                   bool                   `json:"is_standalone_payment"`
	IsVoided                              bool                   `json:"is_voided"`
	IsRefunded                            bool                   `json:"is_refunded"`
	Is3DSecure                            bool                   `json:"is_3d_secure"`
	IntegrationID                         int64                  `json:"integration_id"`
	ProfileID                             int64                  `json:"profile_id"`
	HasParentTransaction                  bool                   `json:"has_parent_transaction"`
	Order                                 PaymobOrder            `json:"order"`
	CreatedAt                             string                 `json:"created_at"`
	TransactionProcessedCallbackResponses []interface{}          `json:"transaction_processed_callback_responses"`
	Currency                              string                 `json:"currency"`
	SourceData                            PaymobSourceData       `json:"source_data"`
	APISource                             string                 `json:"api_source"`
	TerminalID                            interface{}            `json:"terminal_id"`
	MerchantCommission                    int64                  `json:"merchant_commission"`
	AcceptFees                            int64                  `json:"accept_fees"`
	Installment                           interface{}            `json:"installment"`
	DiscountDetails                       []interface{}          `json:"discount_details"`
	IsVoid                                bool                   `json:"is_void"`
	IsRefund                              bool                   `json:"is_refund"`
	Data                                  PaymobTransactionData  `json:"data"`
	IsHidden                              bool                   `json:"is_hidden"`
	PaymentKeyClaims                      PaymobPaymentKeyClaims `json:"payment_key_claims"`
	ErrorOccured                          bool                   `json:"error_occured"`
	IsLive                                bool                   `json:"is_live"`
	OtherEndpointReference                interface{}            `json:"other_endpoint_reference"`
	RefundedAmountCents                   int64                  `json:"refunded_amount_cents"`
	SourceID                              int64                  `json:"source_id"`
	IsCaptured                            bool                   `json:"is_captured"`
	CapturedAmount                        int64                  `json:"captured_amount"`
	MerchantStaffTag                      interface{}            `json:"merchant_staff_tag"`
	UpdatedAt                             string                 `json:"updated_at"`
	IsSettled                             bool                   `json:"is_settled"`
	BillBalanced                          bool                   `json:"bill_balanced"`
	IsBill                                bool                   `json:"is_bill"`
	Owner                                 int64                  `json:"owner"`
	ParentTransaction                     interface{}            `json:"parent_transaction"`
}

type PaymobOrder struct {
	ID                  int64              `json:"id"`
	CreatedAt           string             `json:"created_at"`
	DeliveryNeeded      bool               `json:"delivery_needed"`
	Merchant            PaymobMerchant     `json:"merchant"`
	Collector           interface{}        `json:"collector"`
	AmountCents         int64              `json:"amount_cents"`
	ShippingData        PaymobShippingData `json:"shipping_data"`
	Currency            string             `json:"currency"`
	IsPaymentLocked     bool               `json:"is_payment_locked"`
	IsReturn            bool               `json:"is_return"`
	IsCancel            bool               `json:"is_cancel"`
	IsReturned          bool               `json:"is_returned"`
	IsCanceled          bool               `json:"is_canceled"`
	MerchantOrderID     interface{}        `json:"merchant_order_id"`
	WalletNotification  interface{}        `json:"wallet_notification"`
	PaidAmountCents     int64              `json:"paid_amount_cents"`
	NotifyUserWithEmail bool               `json:"notify_user_with_email"`
	Items               []PaymobItem       `json:"items"`
	OrderURL            string             `json:"order_url"`
	CommissionFees      int64              `json:"commission_fees"`
	DeliveryFeesCents   int64              `json:"delivery_fees_cents"`
	DeliveryVatCents    int64              `json:"delivery_vat_cents"`
	PaymentMethod       string             `json:"payment_method"`
	MerchantStaffTag    interface{}        `json:"merchant_staff_tag"`
	APISource           string             `json:"api_source"`
	Data                interface{}        `json:"data"`
	PaymentStatus       string             `json:"payment_status"`
}

type PaymobMerchant struct {
	ID            int64    `json:"id"`
	CreatedAt     string   `json:"created_at"`
	Phones        []string `json:"phones"`
	CompanyEmails []string `json:"company_emails"`
	CompanyName   string   `json:"company_name"`
	State         string   `json:"state"`
	Country       string   `json:"country"`
	City          string   `json:"city"`
	PostalCode    string   `json:"postal_code"`
	Street        string   `json:"street"`
}

type PaymobShippingData struct {
	ID               int64  `json:"id"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Street           string `json:"street"`
	Building         string `json:"building"`
	Floor            string `json:"floor"`
	Apartment        string `json:"apartment"`
	City             string `json:"city"`
	State            string `json:"state"`
	Country          string `json:"country"`
	Email            string `json:"email"`
	PhoneNumber      string `json:"phone_number"`
	PostalCode       string `json:"postal_code"`
	ExtraDescription string `json:"extra_description"`
	ShippingMethod   string `json:"shipping_method"`
	OrderID          int64  `json:"order_id"`
	Order            int64  `json:"order"`
}

type PaymobItem struct {
	Name        string `json:"name"`
	AmountCents int64  `json:"amount_cents"`
	Quantity    int64  `json:"quantity"`
}

type PaymobSourceData struct {
	Pan     string      `json:"pan"`
	Type    string      `json:"type"`
	Tenure  interface{} `json:"tenure"`
	SubType string      `json:"sub_type"`
}

type PaymobTransactionData struct {
	GatewayIntegrationPk int64                 `json:"gateway_integration_pk"`
	Klass                string                `json:"klass"`
	CreatedAt            string                `json:"created_at"`
	Amount               int64                 `json:"amount"`
	Currency             string                `json:"currency"`
	MigsOrder            PaymobMigsOrder       `json:"migs_order"`
	Merchant             string                `json:"merchant"`
	MigsResult           string                `json:"migs_result"`
	MigsTransaction      PaymobMigsTransaction `json:"migs_transaction"`
	TxnResponseCode      string                `json:"txn_response_code"`
	AcqResponseCode      string                `json:"acq_response_code"`
	Message              string                `json:"message"`
	MerchantTxnRef       string                `json:"merchant_txn_ref"`
	OrderInfo            string                `json:"order_info"`
	ReceiptNo            string                `json:"receipt_no"`
	TransactionNo        string                `json:"transaction_no"`
	BatchNo              int64                 `json:"batch_no"`
	AuthorizeID          string                `json:"authorize_id"`
	CardType             string                `json:"card_type"`
	CardNum              string                `json:"card_num"`
	SecureHash           string                `json:"secure_hash"`
	AvsResultCode        string                `json:"avs_result_code"`
	AvsAcqResponseCode   string                `json:"avs_acq_response_code"`
	CapturedAmount       int64                 `json:"captured_amount"`
	AuthorisedAmount     int64                 `json:"authorised_amount"`
	RefundedAmount       int64                 `json:"refunded_amount"`
	AcsEci               string                `json:"acs_eci"`
}

type PaymobMigsOrder struct {
	AcceptPartialAmount   bool             `json:"acceptPartialAmount"`
	Amount                int64            `json:"amount"`
	AuthenticationStatus  string           `json:"authenticationStatus"`
	Chargeback            PaymobChargeback `json:"chargeback"`
	CreationTime          string           `json:"creationTime"`
	Currency              string           `json:"currency"`
	Description           string           `json:"description"`
	ID                    string           `json:"id"`
	LastUpdatedTime       string           `json:"lastUpdatedTime"`
	MerchantAmount        int64            `json:"merchantAmount"`
	MerchantCategoryCode  string           `json:"merchantCategoryCode"`
	MerchantCurrency      string           `json:"merchantCurrency"`
	Status                string           `json:"status"`
	TotalAuthorizedAmount int64            `json:"totalAuthorizedAmount"`
	TotalCapturedAmount   int64            `json:"totalCapturedAmount"`
	TotalRefundedAmount   int64            `json:"totalRefundedAmount"`
}

type PaymobChargeback struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type PaymobMigsTransaction struct {
	Acquirer             PaymobAcquirer `json:"acquirer"`
	Amount               int64          `json:"amount"`
	AuthenticationStatus string         `json:"authenticationStatus"`
	AuthorizationCode    string         `json:"authorizationCode"`
	Currency             string         `json:"currency"`
	ID                   string         `json:"id"`
	Receipt              string         `json:"receipt"`
	Source               string         `json:"source"`
	Stan                 string         `json:"stan"`
	Terminal             string         `json:"terminal"`
	Type                 string         `json:"type"`
}

type PaymobAcquirer struct {
	Batch          int64  `json:"batch"`
	Date           string `json:"date"`
	ID             string `json:"id"`
	MerchantID     string `json:"merchantId"`
	SettlementDate string `json:"settlementDate"`
	TimeZone       string `json:"timeZone"`
	TransactionID  string `json:"transactionId"`
}

type PaymobPaymentKeyClaims struct {
	Extra                PaymobSubscrbeResponseExtra `json:"extra"`
	UserID               int64                       `json:"user_id"`
	Currency             string                      `json:"currency"`
	OrderID              int64                       `json:"order_id"`
	CreatedBy            int64                       `json:"created_by"`
	IsPartner            bool                        `json:"is_partner"`
	AmountCents          int64                       `json:"amount_cents"`
	BillingData          PaymobBillingData           `json:"billing_data"`
	RedirectURL          string                      `json:"redirect_url"`
	IntegrationID        int64                       `json:"integration_id"`
	LockOrderWhenPaid    bool                        `json:"lock_order_when_paid"`
	SubscriptionPlanID   string                      `json:"subscription_plan_id"`
	NextPaymentIntention string                      `json:"next_payment_intention"`
	SinglePaymentAttempt bool                        `json:"single_payment_attempt"`
}

type PaymobSubscrbeResponseExtra struct {
	CreationExtras struct {
		Ee              int     `json:"ee"`
		MerchantOrderID *string `json:"merchant_order_id"`
	} `json:"creation_extras"`
	ConfirmationExtras *map[string]interface{} `json:"confirmation_extras"`
}

// ///////////// Paymob Token generation https://accept.paymob.com/api/auth/tokens
type PaymobTokensPostResponseUser struct {
	ID              int     `json:"id"`
	Username        string  `json:"username"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	DateJoined      string  `json:"date_joined"`
	Email           string  `json:"email"`
	IsActive        bool    `json:"is_active"`
	IsStaff         bool    `json:"is_staff"`
	IsSuperuser     bool    `json:"is_superuser"`
	LastLogin       *string `json:"last_login"`
	Groups          []any   `json:"groups"`
	UserPermissions []int   `json:"user_permissions"`
}

type PaymobTokensPostResponseProfile struct {
	ID                                      int                          `json:"id"`
	User                                    PaymobTokensPostResponseUser `json:"user"`
	CreatedAt                               string                       `json:"created_at"`
	Active                                  bool                         `json:"active"`
	ProfileType                             string                       `json:"profile_type"`
	Phones                                  []string                     `json:"phones"`
	CompanyEmails                           []string                     `json:"company_emails"`
	CompanyName                             string                       `json:"company_name"`
	State                                   string                       `json:"state"`
	Country                                 string                       `json:"country"`
	City                                    string                       `json:"city"`
	PostalCode                              string                       `json:"postal_code"`
	Street                                  string                       `json:"street"`
	EmailNotification                       bool                         `json:"email_notification"`
	OrderRetrievalEndpoint                  *string                      `json:"order_retrieval_endpoint"`
	DeliveryUpdateEndpoint                  *string                      `json:"delivery_update_endpoint"`
	LogoURL                                 *string                      `json:"logo_url"`
	IsMobadra                               bool                         `json:"is_mobadra"`
	Sector                                  string                       `json:"sector"`
	Is2FAEnabled                            bool                         `json:"is_2fa_enabled"`
	OTPSentTo                               string                       `json:"otp_sent_to"`
	ActivationMethod                        int                          `json:"activation_method"`
	SignedUpThrough                         int                          `json:"signed_up_through"`
	FailedAttempts                          int                          `json:"failed_attempts"`
	CustomExportColumns                     []any                        `json:"custom_export_columns"`
	ServerIP                                []any                        `json:"server_IP"`
	Username                                *string                      `json:"username"`
	PrimaryPhoneNumber                      string                       `json:"primary_phone_number"`
	PrimaryPhoneVerified                    bool                         `json:"primary_phone_verified"`
	IsTempPassword                          bool                         `json:"is_temp_password"`
	OTP2FASentAt                            *string                      `json:"otp_2fa_sent_at"`
	OTP2FAAttempt                           any                          `json:"otp_2fa_attempt"`
	OTPSentAt                               string                       `json:"otp_sent_at"`
	OTPValidatedAt                          *string                      `json:"otp_validated_at"`
	AWBBanner                               *string                      `json:"awb_banner"`
	EmailBanner                             *string                      `json:"email_banner"`
	IdentificationNumber                    *string                      `json:"identification_number"`
	DeliveryStatusCallback                  string                       `json:"delivery_status_callback"`
	MerchantExternalLink                    *string                      `json:"merchant_external_link"`
	MerchantStatus                          int                          `json:"merchant_status"`
	DeactivatedByBank                       bool                         `json:"deactivated_by_bank"`
	BankDeactivationReason                  *string                      `json:"bank_deactivation_reason"`
	BankMerchantStatus                      int                          `json:"bank_merchant_status"`
	NationalID                              *string                      `json:"national_id"`
	SuperAgent                              any                          `json:"super_agent"`
	WalletLimitProfile                      any                          `json:"wallet_limit_profile"`
	Address                                 *string                      `json:"address"`
	CommercialRegistration                  *string                      `json:"commercial_registration"`
	CommercialRegistrationArea              *string                      `json:"commercial_registration_area"`
	DistributorCode                         *string                      `json:"distributor_code"`
	DistributorBranchCode                   *string                      `json:"distributor_branch_code"`
	AllowTerminalOrderID                    bool                         `json:"allow_terminal_order_id"`
	AllowEncryptionBypass                   bool                         `json:"allow_encryption_bypass"`
	WalletPhoneNumber                       *string                      `json:"wallet_phone_number"`
	Suspicious                              int                          `json:"suspicious"`
	Latitude                                any                          `json:"latitude"`
	Longitude                               any                          `json:"longitude"`
	BankStaffs                              map[string]any               `json:"bank_staffs"`
	BankRejectionReason                     *string                      `json:"bank_rejection_reason"`
	BankReceivedDocuments                   bool                         `json:"bank_received_documents"`
	BankMerchantDigitalStatus               int                          `json:"bank_merchant_digital_status"`
	BankDigitalRejectionReason              *string                      `json:"bank_digital_rejection_reason"`
	FilledBusinessData                      bool                         `json:"filled_business_data"`
	DayStartTime                            string                       `json:"day_start_time"`
	DayEndTime                              *string                      `json:"day_end_time"`
	WithholdTransfers                       bool                         `json:"withhold_transfers"`
	ManualSettlement                        bool                         `json:"manual_settlement"`
	SMSSenderName                           string                       `json:"sms_sender_name"`
	WithholdTransfersReason                 *string                      `json:"withhold_transfers_reason"`
	WithholdTransfersNotes                  *string                      `json:"withhold_transfers_notes"`
	CanBillDepositWithCard                  bool                         `json:"can_bill_deposit_with_card"`
	CanTopupMerchants                       bool                         `json:"can_topup_merchants"`
	TopupTransferID                         any                          `json:"topup_transfer_id"`
	ReferralEligible                        bool                         `json:"referral_eligible"`
	IsEligibleToBeRanger                    bool                         `json:"is_eligible_to_be_ranger"`
	EligibleForManualRefunds                bool                         `json:"eligible_for_manual_refunds"`
	IsRanger                                bool                         `json:"is_ranger"`
	IsPoaching                              bool                         `json:"is_poaching"`
	PaymobAppMerchant                       bool                         `json:"paymob_app_merchant"`
	SettlementFrequency                     *string                      `json:"settlement_frequency"`
	DayOfTheWeek                            *string                      `json:"day_of_the_week"`
	DayOfTheMonth                           *string                      `json:"day_of_the_month"`
	AllowTransactionNotifications           bool                         `json:"allow_transaction_notifications"`
	AllowTransferNotifications              bool                         `json:"allow_transfer_notifications"`
	SallefnyAmountWhole                     float64                      `json:"sallefny_amount_whole"`
	SallefnyFeesWhole                       float64                      `json:"sallefny_fees_whole"`
	PaymobAppFirstLogin                     *string                      `json:"paymob_app_first_login"`
	PaymobAppLastActivity                   *string                      `json:"paymob_app_last_activity"`
	PayoutEnabled                           bool                         `json:"payout_enabled"`
	PayoutTerms                             bool                         `json:"payout_terms"`
	IsBillsNew                              bool                         `json:"is_bills_new"`
	CanProcessMultipleRefunds               bool                         `json:"can_process_multiple_refunds"`
	SettlementClassification                int                          `json:"settlement_classification"`
	VATClassification                       int                          `json:"vat_classification"`
	InstantSettlementEnabled                bool                         `json:"instant_settlement_enabled"`
	InstantSettlementTransactionOTPVerified bool                         `json:"instant_settlement_transaction_otp_verified"`
	PreferredLanguage                       *string                      `json:"preferred_language"`
	IgnoreFlashCallbacks                    bool                         `json:"ignore_flash_callbacks"`
	ReceiveCallbackCardInfo                 bool                         `json:"receive_callback_card_info"`
	AcqPartner                              any                          `json:"acq_partner"`
	DOM                                     any                          `json:"dom"`
	BankRelated                             any                          `json:"bank_related"`
	Permissions                             []any                        `json:"permissions"`
}

type PaymobTokensPostResponse struct {
	Profile PaymobTokensPostResponseProfile `json:"profile"`
	Token   string                          `json:"token"`
}

///////////////////!SECTION

// ///////// !Paymob Transaction inquiry
type PaymobTransactionInquiryMerchant struct {
	ID            int      `json:"id"`
	CreatedAt     string   `json:"created_at"`
	Phones        []string `json:"phones"`
	CompanyEmails []string `json:"company_emails"`
	CompanyName   string   `json:"company_name"`
	State         string   `json:"state"`
	Country       string   `json:"country"`
	City          string   `json:"city"`
	PostalCode    string   `json:"postal_code"`
	Street        string   `json:"street"`
}

type PaymobTransactionInquiryShippingData struct {
	ID               int    `json:"id"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Street           string `json:"street"`
	Building         string `json:"building"`
	Floor            string `json:"floor"`
	Apartment        string `json:"apartment"`
	City             string `json:"city"`
	State            string `json:"state"`
	Country          string `json:"country"`
	Email            string `json:"email"`
	PhoneNumber      string `json:"phone_number"`
	PostalCode       string `json:"postal_code"`
	ExtraDescription string `json:"extra_description"`
	ShippingMethod   string `json:"shipping_method"`
	OrderID          int    `json:"order_id"`
	Order            int    `json:"order"`
}

type PaymobTransactionInquiryItem struct {
	Name        string `json:"name"`
	AmountCents int    `json:"amount_cents"`
	Quantity    int    `json:"quantity"`
}

type PaymobTransactionInquiryOrder struct {
	ID                  int                                  `json:"id"`
	CreatedAt           string                               `json:"created_at"`
	DeliveryNeeded      bool                                 `json:"delivery_needed"`
	Merchant            PaymobTransactionInquiryMerchant     `json:"merchant"`
	Collector           interface{}                          `json:"collector"`
	AmountCents         int                                  `json:"amount_cents"`
	ShippingData        PaymobTransactionInquiryShippingData `json:"shipping_data"`
	Currency            string                               `json:"currency"`
	IsPaymentLocked     bool                                 `json:"is_payment_locked"`
	IsReturn            bool                                 `json:"is_return"`
	IsCancel            bool                                 `json:"is_cancel"`
	IsReturned          bool                                 `json:"is_returned"`
	IsCanceled          bool                                 `json:"is_canceled"`
	MerchantOrderID     interface{}                          `json:"merchant_order_id"`
	WalletNotification  interface{}                          `json:"wallet_notification"`
	PaidAmountCents     int                                  `json:"paid_amount_cents"`
	NotifyUserWithEmail bool                                 `json:"notify_user_with_email"`
	Items               []PaymobTransactionInquiryItem       `json:"items"`
	OrderURL            string                               `json:"order_url"`
	CommissionFees      int                                  `json:"commission_fees"`
	DeliveryFeesCents   int                                  `json:"delivery_fees_cents"`
	DeliveryVatCents    int                                  `json:"delivery_vat_cents"`
	PaymentMethod       string                               `json:"payment_method"`
	MerchantStaffTag    interface{}                          `json:"merchant_staff_tag"`
	APISource           string                               `json:"api_source"`
	Data                map[string]interface{}               `json:"data"`
	PaymentStatus       string                               `json:"payment_status"`
}

type PaymobTransactionInquiryCallbackResponse struct {
	Response struct {
		Status   string `json:"status"`
		Content  string `json:"content"`
		Headers  string `json:"headers"`
		Encoding string `json:"encoding"`
	} `json:"response"`
	CallbackURL        string `json:"callback_url"`
	ResponseReceivedAt string `json:"response_received_at"`
}

type PaymobTransactionInquirySourceData struct {
	Pan     string      `json:"pan"`
	Type    string      `json:"type"`
	Tenure  interface{} `json:"tenure"`
	SubType string      `json:"sub_type"`
}

type PaymobTransactionInquiryMigsOrder struct {
	ID         string  `json:"id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	Currency   string  `json:"currency"`
	Chargeback struct {
		Amount   int    `json:"amount"`
		Currency string `json:"currency"`
	} `json:"chargeback"`
	Description           string  `json:"description"`
	CreationTime          string  `json:"creationTime"`
	MerchantAmount        float64 `json:"merchantAmount"`
	LastUpdatedTime       string  `json:"lastUpdatedTime"`
	MerchantCurrency      string  `json:"merchantCurrency"`
	AcceptPartialAmount   bool    `json:"acceptPartialAmount"`
	TotalCapturedAmount   float64 `json:"totalCapturedAmount"`
	TotalRefundedAmount   float64 `json:"totalRefundedAmount"`
	AuthenticationStatus  string  `json:"authenticationStatus"`
	MerchantCategoryCode  string  `json:"merchantCategoryCode"`
	TotalAuthorizedAmount float64 `json:"totalAuthorizedAmount"`
}

type PaymobTransactionInquiryAcquirer struct {
	ID             string `json:"id"`
	Date           string `json:"date"`
	Batch          int    `json:"batch"`
	TimeZone       string `json:"timeZone"`
	MerchantID     string `json:"merchantId"`
	TransactionID  string `json:"transactionId"`
	SettlementDate string `json:"settlementDate"`
}

type PaymobTransactionInquiryMigsTransaction struct {
	ID                   string                           `json:"id"`
	Stan                 string                           `json:"stan"`
	Type                 string                           `json:"type"`
	Amount               float64                          `json:"amount"`
	Source               string                           `json:"source"`
	Receipt              string                           `json:"receipt"`
	Acquirer             PaymobTransactionInquiryAcquirer `json:"acquirer"`
	Currency             string                           `json:"currency"`
	Terminal             string                           `json:"terminal"`
	AuthorizationCode    string                           `json:"authorizationCode"`
	AuthenticationStatus string                           `json:"authenticationStatus"`
}

type PaymobTransactionInquiryTransactionData struct {
	Klass                string                                  `json:"klass"`
	Amount               float64                                 `json:"amount"`
	AcsEci               string                                  `json:"acs_eci"`
	Message              string                                  `json:"message"`
	BatchNo              int                                     `json:"batch_no"`
	CardNum              string                                  `json:"card_num"`
	Currency             string                                  `json:"currency"`
	Merchant             string                                  `json:"merchant"`
	CardType             string                                  `json:"card_type"`
	CreatedAt            string                                  `json:"created_at"`
	MigsOrder            PaymobTransactionInquiryMigsOrder       `json:"migs_order"`
	OrderInfo            string                                  `json:"order_info"`
	ReceiptNo            string                                  `json:"receipt_no"`
	MigsResult           string                                  `json:"migs_result"`
	SecureHash           string                                  `json:"secure_hash"`
	AuthorizeID          string                                  `json:"authorize_id"`
	TransactionNo        string                                  `json:"transaction_no"`
	AvsResultCode        string                                  `json:"avs_result_code"`
	CapturedAmount       float64                                 `json:"captured_amount"`
	RefundedAmount       float64                                 `json:"refunded_amount"`
	MerchantTxnRef       string                                  `json:"merchant_txn_ref"`
	MigsTransaction      PaymobTransactionInquiryMigsTransaction `json:"migs_transaction"`
	AcqResponseCode      string                                  `json:"acq_response_code"`
	AuthorisedAmount     float64                                 `json:"authorised_amount"`
	TxnResponseCode      string                                  `json:"txn_response_code"`
	AvsAcqResponseCode   string                                  `json:"avs_acq_response_code"`
	GatewayIntegrationPk int                                     `json:"gateway_integration_pk"`
}

type PaymobTransactionInquiryBillingData struct {
	City             string `json:"city"`
	Email            string `json:"email"`
	Floor            string `json:"floor"`
	State            string `json:"state"`
	Street           string `json:"street"`
	Country          string `json:"country"`
	Building         string `json:"building"`
	Apartment        string `json:"apartment"`
	LastName         string `json:"last_name"`
	FirstName        string `json:"first_name"`
	PostalCode       string `json:"postal_code"`
	PhoneNumber      string `json:"phone_number"`
	ExtraDescription string `json:"extra_description"`
}

type PaymobTransactionInquiryPaymentKeyClaims struct {
	Extra struct {
		TenantID        string      `json:"tenant_id"`
		StartsAt        time.Time   `json:"starts_at"`
		EndsAt          time.Time   `json:"ends_at"`
		MerchantOrderID interface{} `json:"merchant_order_id"`
	} `json:"extra"`
	UserID               int                                 `json:"user_id"`
	Currency             string                              `json:"currency"`
	OrderID              int                                 `json:"order_id"`
	CreatedBy            int                                 `json:"created_by"`
	IsPartner            bool                                `json:"is_partner"`
	AmountCents          int                                 `json:"amount_cents"`
	BillingData          PaymobTransactionInquiryBillingData `json:"billing_data"`
	RedirectURL          string                              `json:"redirect_url"`
	IntegrationID        int                                 `json:"integration_id"`
	LockOrderWhenPaid    bool                                `json:"lock_order_when_paid"`
	SubscriptionPlanID   string                              `json:"subscription_plan_id"`
	NextPaymentIntention string                              `json:"next_payment_intention"`
	SinglePaymentAttempt bool                                `json:"single_payment_attempt"`
}

type PaymobTransactionInquiryResponse struct {
	ID                                    int                                        `json:"id"`
	Pending                               bool                                       `json:"pending"`
	AmountCents                           int                                        `json:"amount_cents"`
	Success                               bool                                       `json:"success"`
	IsAuth                                bool                                       `json:"is_auth"`
	IsCapture                             bool                                       `json:"is_capture"`
	IsStandalonePayment                   bool                                       `json:"is_standalone_payment"`
	IsVoided                              bool                                       `json:"is_voided"`
	IsRefunded                            bool                                       `json:"is_refunded"`
	Is3DSecure                            bool                                       `json:"is_3d_secure"`
	IntegrationID                         int                                        `json:"integration_id"`
	ProfileID                             int                                        `json:"profile_id"`
	HasParentTransaction                  bool                                       `json:"has_parent_transaction"`
	Order                                 PaymobTransactionInquiryOrder              `json:"order"`
	CreatedAt                             string                                     `json:"created_at"`
	TransactionProcessedCallbackResponses []PaymobTransactionInquiryCallbackResponse `json:"transaction_processed_callback_responses"`
	Currency                              string                                     `json:"currency"`
	SourceData                            PaymobTransactionInquirySourceData         `json:"source_data"`
	APISource                             string                                     `json:"api_source"`
	TerminalID                            interface{}                                `json:"terminal_id"`
	MerchantCommission                    int                                        `json:"merchant_commission"`
	AcceptFees                            int                                        `json:"accept_fees"`
	Installment                           interface{}                                `json:"installment"`
	DiscountDetails                       []interface{}                              `json:"discount_details"`
	IsVoid                                bool                                       `json:"is_void"`
	IsRefund                              bool                                       `json:"is_refund"`
	Data                                  PaymobTransactionInquiryTransactionData    `json:"data"`
	IsHidden                              bool                                       `json:"is_hidden"`
	PaymentKeyClaims                      PaymobTransactionInquiryPaymentKeyClaims   `json:"payment_key_claims"`
	ErrorOccured                          bool                                       `json:"error_occured"`
	IsLive                                bool                                       `json:"is_live"`
	OtherEndpointReference                interface{}                                `json:"other_endpoint_reference"`
	RefundedAmountCents                   int                                        `json:"refunded_amount_cents"`
	SourceID                              int                                        `json:"source_id"`
	IsCaptured                            bool                                       `json:"is_captured"`
	CapturedAmount                        int                                        `json:"captured_amount"`
	MerchantStaffTag                      interface{}                                `json:"merchant_staff_tag"`
	UpdatedAt                             string                                     `json:"updated_at"`
	IsSettled                             bool                                       `json:"is_settled"`
	BillBalanced                          bool                                       `json:"bill_balanced"`
	IsBill                                bool                                       `json:"is_bill"`
	Owner                                 int                                        `json:"owner"`
	ParentTransaction                     interface{}                                `json:"parent_transaction"`
	UniqueRef                             string                                     `json:"unique_ref"`
}

//////////// !Paymob Transaction inquiry

func PaymobSubscribeCallbackPOST(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := PaymobSubscribePaymentCallback{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tokenUrl := "https://accept.paymob.com/api/auth/tokens"
		tokenRequest := struct {
			ApiKey string `json:"api_key"`
		}{
			ApiKey: config.Payment.ApiKey,
		}
		tokenReqBody, err := json.Marshal(tokenRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tokenReq, err := http.NewRequest("POST", tokenUrl, bytes.NewBuffer(tokenReqBody))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tokenReq.Header.Set("Content-Type", "application/json")
		tokenResp, err := http.DefaultClient.Do(tokenReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer tokenResp.Body.Close()

		var tokenRespBody PaymobTokensPostResponse
		err = json.NewDecoder(tokenResp.Body).Decode(&tokenRespBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// !Get order details
		transactionInquiryUrl := "https://accept.paymob.com/api/ecommerce/orders/transaction_inquiry"
		transactionInquiryReqBody, err := json.Marshal(struct {
			OrderId int64 `json:"order_id"`
		}{
			OrderId: request.Obj.Order.ID,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		transactionInquiryReq, err := http.NewRequest("POST", transactionInquiryUrl, bytes.NewBuffer(transactionInquiryReqBody))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		transactionInquiryReq.Header.Set("Content-Type", "application/json")
		transactionInquiryReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenRespBody.Token))

		resp, err := http.DefaultClient.Do(transactionInquiryReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, fmt.Sprintf("Paymob transaction inquiry request responded with status code %d", resp.StatusCode), http.StatusInternalServerError)
			return
		}

		var transactionInquiryResp PaymobTransactionInquiryResponse
		err = json.NewDecoder(resp.Body).Decode(&transactionInquiryResp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//// !!Get order details

		if transactionInquiryResp.Success && !transactionInquiryResp.IsRefunded && !transactionInquiryResp.IsRefund && !transactionInquiryResp.Pending {
			tenant_id := transactionInquiryResp.PaymentKeyClaims.Extra.TenantID
			// created_at := transactionInquiryResp.CreatedAt
			// response_received_at := transactionInquiryResp.TransactionProcessedCallbackResponses[0].ResponseReceivedAt
			starts_at := transactionInquiryResp.PaymentKeyClaims.Extra.StartsAt
			ends_at := transactionInquiryResp.PaymentKeyClaims.Extra.EndsAt

			if transactionInquiryResp.Order.Items[0].Name == "Standard subscription" {
				today := time.Now()

				clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", config.Databases[0].Host, config.Databases[0].Port))
				deadline := 5 * time.Second
				if config.Env == "dev" {
					deadline = 1000 * time.Second
				}

				ctx, cancel := context.WithTimeout(context.Background(), deadline)
				defer cancel()

				client, err := mongo.Connect(ctx, clientOptions)

				// connected to db

				collection := client.Database(config.Databases[0].Database).Collection(config.Databases[0].Tables["sales"])
				// check if document with tenant_id exists, otherwise create it with empty object value
				filter := bson.D{{Key: "tenant_id", Value: tenant_id}}
				var tenant models.Tenant
				err = collection.FindOne(ctx, filter).Decode(&tenant)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					logger.Error(err.Error())
					return
				}

				if today.Before(ends_at) && today.After(starts_at) {
					tenant.Subscription.SubscriptionPlan = "standard"
					tenant.Subscription.StartDate = starts_at
					tenant.Subscription.EndDate = ends_at
				}

				_, err = collection.UpdateOne(ctx, filter, bson.M{
					"$set": bson.M{
						"subscription": tenant.Subscription,
					},
				})
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					logger.Error(err.Error())
					return
				}
			}

			if transactionInquiryResp.Order.Items[0].Name == "Gold subscription" {
				today := time.Now()

				clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", config.Databases[0].Host, config.Databases[0].Port))
				deadline := 5 * time.Second
				if config.Env == "dev" {
					deadline = 1000 * time.Second
				}

				ctx, cancel := context.WithTimeout(context.Background(), deadline)
				defer cancel()

				client, err := mongo.Connect(ctx, clientOptions)

				// connected to db

				collection := client.Database(config.Databases[0].Database).Collection(config.Databases[0].Tables["sales"])
				// check if document with tenant_id exists, otherwise create it with empty object value
				filter := bson.D{{Key: "tenant_id", Value: tenant_id}}
				var tenant models.Tenant
				err = collection.FindOne(ctx, filter).Decode(&tenant)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					logger.Error(err.Error())
					return
				}

				if today.Before(ends_at) && today.After(starts_at) {
					tenant.Subscription.SubscriptionPlan = "gold"
					tenant.Subscription.StartDate = starts_at
					tenant.Subscription.EndDate = ends_at
				}

				_, err = collection.UpdateOne(ctx, filter, bson.M{
					"$set": bson.M{
						"subscription": tenant.Subscription,
					},
				})
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					logger.Error(err.Error())
					return
				}
			}
		}

	}
}

func SubcriptionRequest(config config.Config, logger logger.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenant_id := "1"
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

			subscription_request := PaymobSubscribeRequest{
				Amount:             150000, // Example amount in cents
				Currency:           "EGP",
				PaymentMethods:     []int{5229518},
				SubscriptionPlanID: 4088, // Example plan ID
				Items: []struct {
					Name     string `json:"name"`
					Amount   int    `json:"amount"`
					Quantity int    `json:"quantity"`
				}{{
					Name:     "Standard subscription",
					Amount:   150000,
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
					FirstName:   given_name,
					LastName:    family_name,
					Street:      "",
					Building:    "",
					PhoneNumber: phone,
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
					FirstName: given_name,
					LastName:  family_name,
					Email:     client_email,
					Extras: map[string]interface{}{
						"tenant_id": tenant_id,
						"starts_at": time.Now().Format(time.RFC3339),
						"ends_at":   time.Now().AddDate(0, 1, 0).Format(time.RFC3339),
					},
				},
				Extras: map[string]interface{}{
					"tenant_id": tenant_id,
					"starts_at": time.Now().Format(time.RFC3339),
					"ends_at":   time.Now().AddDate(0, 1, 0).Format(time.RFC3339),
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

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			req = req.WithContext(ctx)
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

			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
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

			var subscribeResponse PaymobSubscribeResponse
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
					PaymentURL       string `json:"payment_url"`
				}{
					SubscriptionPlan: "standard",
					ClientSecret:     subscribeResponse.ClientSecret,
					PublicKey:        config.Payment.PublicKey,
					PaymentURL:       fmt.Sprintf("https://accept.paymob.com/unifiedcheckout/?publicKey=%s&clientSecret=%s", config.Payment.PublicKey, subscribeResponse.ClientSecret),
				},
			}

			jsonResponse, err := json.Marshal(response)
			if err != nil {
				logger.Error(err.Error())
				http.Error(w, "Failed to marshal order settings response", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonResponse)

		case "gold":

			subscription_request := PaymobSubscribeRequest{
				Amount:             250000, // Example amount in cents
				Currency:           "EGP",
				PaymentMethods:     []int{5229518},
				SubscriptionPlanID: 4203, // Example plan ID
				Items: []struct {
					Name     string `json:"name"`
					Amount   int    `json:"amount"`
					Quantity int    `json:"quantity"`
				}{{
					Name:     "Gold subscription",
					Amount:   250000,
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
					FirstName:   given_name,
					LastName:    family_name,
					Street:      "",
					Building:    "",
					PhoneNumber: phone,
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
					FirstName: given_name,
					LastName:  family_name,
					Email:     client_email,
					Extras: map[string]interface{}{
						"tenant_id": tenant_id,
						"starts_at": time.Now().Format(time.RFC3339),
						"ends_at":   time.Now().AddDate(0, 1, 0).Format(time.RFC3339),
					},
				},
				Extras: map[string]interface{}{
					"tenant_id": tenant_id,
					"starts_at": time.Now().Format(time.RFC3339),
					"ends_at":   time.Now().AddDate(0, 1, 0).Format(time.RFC3339),
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

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			req = req.WithContext(ctx)
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

			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
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

			var subscribeResponse PaymobSubscribeResponse
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
					PaymentURL       string `json:"payment_url"`
				}{
					SubscriptionPlan: "gold",
					ClientSecret:     subscribeResponse.ClientSecret,
					PublicKey:        config.Payment.PublicKey,
					PaymentURL:       fmt.Sprintf("https://accept.paymob.com/unifiedcheckout/?publicKey=%s&clientSecret=%s", config.Payment.PublicKey, subscribeResponse.ClientSecret),
				},
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

		if tenant.Subscription.SubscriptionPlan == "" {
			tenant.Subscription = models.TenantSubscription{
				ID:               primitive.NewObjectID().Hex(),
				SubscriptionPlan: "free",
				StartDate:        time.Now(),
				EndDate:          time.Now().AddDate(99, 0, 0), // 99 years from now
				Status:           "active",
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
