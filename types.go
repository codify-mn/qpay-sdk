package qpay

// Config is a convenience struct matching what many consumers load from YAML.
// Equivalent to passing WithBaseURL / WithCredentials / WithTerminalID separately.
type Config struct {
	BaseURL    string `yaml:"base_url"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	TerminalID string `yaml:"terminal_id"`
}

// TokenResponse is the raw /v2/auth/token response.
type TokenResponse struct {
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
}

// BankAccount represents a merchant bank account.
type BankAccount struct {
	AccountBankCode string `json:"account_bank_code"`
	AccountNumber   string `json:"account_number"`
	AccountName     string `json:"account_name"`
	IsDefault       bool   `json:"is_default"`
}

// CreateCompanyMerchantRequest registers a company as a QPay merchant.
type CreateCompanyMerchantRequest struct {
	OwnerRegisterNo string      `json:"owner_register_no"`
	OwnerFirstName  string      `json:"owner_first_name"`
	OwnerLastName   string      `json:"owner_last_name"`
	RegisterNumber  string      `json:"register_number"`
	Name            string      `json:"name"`
	MCCCode         string      `json:"mcc_code"`
	City            string      `json:"city"`
	District        string      `json:"district"`
	Address         string      `json:"address"`
	Phone           string      `json:"phone"`
	Email           string      `json:"email"`
	BankAccount     BankAccount `json:"bank_account"`
}

// CreatePersonMerchantRequest registers a person as a QPay merchant.
type CreatePersonMerchantRequest struct {
	RegisterNumber string      `json:"register_number"`
	FirstName      string      `json:"first_name"`
	LastName       string      `json:"last_name"`
	MCCCode        string      `json:"mcc_code"`
	City           string      `json:"city"`
	District       string      `json:"district"`
	Address        string      `json:"address"`
	BusinessName   string      `json:"business_name"`
	Phone          string      `json:"phone"`
	Email          string      `json:"email"`
	BankAccount    BankAccount `json:"bank_account"`
}

// UpdateMerchantRequest updates merchant details.
type UpdateMerchantRequest struct {
	Name        string      `json:"name,omitempty"`
	MCCCode     string      `json:"mcc_code,omitempty"`
	City        string      `json:"city,omitempty"`
	District    string      `json:"district,omitempty"`
	Address     string      `json:"address,omitempty"`
	Phone       string      `json:"phone,omitempty"`
	Email       string      `json:"email,omitempty"`
	BankAccount BankAccount `json:"bank_account,omitempty"`
}

// Merchant is a QPay merchant record.
type Merchant struct {
	ID             string `json:"id"`
	VendorID       string `json:"vendor_id"`
	Type           string `json:"type"`
	RegisterNumber string `json:"register_number"`
	Name           string `json:"name"`
	P2PTerminalID  string `json:"p2p_terminal_id"`
	CardTerminalID string `json:"card_terminal_id"`
	MCCCode        string `json:"mcc_code"`
	City           string `json:"city"`
	District       string `json:"district"`
	Address        string `json:"address"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
}

// MerchantList is the paginated response from ListMerchants.
type MerchantList struct {
	Count   int        `json:"count"`
	Rows    []Merchant `json:"rows"`
	Offset  int        `json:"offset"`
	Limit   int        `json:"limit"`
	Message string     `json:"message,omitempty"`
}

// ListOptions is pagination for list endpoints.
type ListOptions struct {
	Offset int
	Limit  int
}

// InvoiceLine is a line item on an invoice.
type InvoiceLine struct {
	Name     string  `json:"line_description"`
	Quantity int     `json:"line_quantity"`
	Price    float64 `json:"line_unit_price"`
}

// CreateInvoiceRequest creates a payment invoice via the QPay Quick Pay API.
//
// Required: MerchantID, Amount, Currency, CustomerName, CallbackURL, Description, BankAccounts.
type CreateInvoiceRequest struct {
	MerchantID   string        `json:"merchant_id"`
	BranchCode   string        `json:"branch_code,omitempty"`
	Amount       float64       `json:"amount"`
	Currency     string        `json:"currency"`      // "MNT"
	CustomerName string        `json:"customer_name"` // merchant-chosen display name
	CustomerLogo string        `json:"customer_logo,omitempty"`
	CallbackURL  string        `json:"callback_url"`
	Description  string        `json:"description"`
	MCCCode      string        `json:"mcc_code,omitempty"`
	BankAccounts []BankAccount `json:"bank_accounts,omitempty"`
}

// Invoice is a created invoice containing QR + deep links.
type Invoice struct {
	InvoiceID string       `json:"id"`
	QRCode    string       `json:"qr_code"`
	QRImage   string       `json:"qr_image"`
	URLs      []InvoiceURL `json:"urls"`
	ShortURL  string       `json:"qPay_shortUrl,omitempty"`
}

// InvoiceURL is a bank app deep link rendered on the client.
type InvoiceURL struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Link        string `json:"link"`
}

// InvoiceDetail is the detailed invoice record from GetInvoice.
type InvoiceDetail struct {
	InvoiceID   string  `json:"invoice_id"`
	InvoiceCode string  `json:"invoice_code"`
	Amount      float64 `json:"amount"`
	Status      string  `json:"invoice_status"` // "OPEN" | "PAID" | "CANCELLED"
	Description string  `json:"invoice_description"`
	CreatedDate string  `json:"created_date"`
	PaidDate    string  `json:"paid_date,omitempty"`
}

// checkPaymentRequest is the internal body for POST /v2/payment/check.
type checkPaymentRequest struct {
	ObjectType string `json:"object_type"`
	ObjectID   string `json:"object_id"`
	Offset     int    `json:"offset,omitempty"`
	Limit      int    `json:"limit,omitempty"`
}

// PaymentCheck is the response from CheckPayment.
type PaymentCheck struct {
	Count       int       `json:"count"`
	PaidAmount  float64   `json:"paid_amount"`
	CheckedDate string    `json:"checked_date"`
	Rows        []Payment `json:"rows"`
}

// IsPaid reports whether any row in the check response has status "PAID".
func (p *PaymentCheck) IsPaid() (bool, *Payment) {
	if p == nil {
		return false, nil
	}
	for i := range p.Rows {
		if p.Rows[i].PaymentStatus == "PAID" {
			return true, &p.Rows[i]
		}
	}
	return false, nil
}

// Payment is one payment transaction.
type Payment struct {
	PaymentID     string  `json:"payment_id"`
	PaymentStatus string  `json:"payment_status"`
	PaymentDate   string  `json:"payment_date"`
	PaymentAmount float64 `json:"payment_amount"`
	PaymentFee    float64 `json:"payment_fee"`
	PaymentBank   string  `json:"payment_bank"`
	PaymentMethod string  `json:"payment_method"`
	TerminalID    string  `json:"terminal_id"`
}

// PaymentListRequest filters the POST /v2/payment/list endpoint.
type PaymentListRequest struct {
	MerchantID     string `json:"merchant_id,omitempty"`
	MerchantBranch string `json:"merchant_branch_code,omitempty"`
	PaymentStatus  string `json:"payment_status,omitempty"`
	StartDate      string `json:"start_date,omitempty"`
	EndDate        string `json:"end_date,omitempty"`
	Offset         int    `json:"offset,omitempty"`
	Limit          int    `json:"limit,omitempty"`
}

// PaymentList is the response from ListPayments.
type PaymentList struct {
	Count      int       `json:"count"`
	PaidAmount float64   `json:"paid_amount"`
	Rows       []Payment `json:"rows"`
}

// RefundRequest supplies a reason for refund/cancel payment endpoints.
type RefundRequest struct {
	CallbackURL string `json:"callback_url,omitempty"`
	Note        string `json:"note,omitempty"`
}

// CreateEbarimtRequest creates an eBarimt receipt.
type CreateEbarimtRequest struct {
	PaymentID       string `json:"payment_id"`
	EbarimtReceiver string `json:"ebarimt_receiver_type"` // "CITIZEN" | "ORGANIZATION"
}

// Ebarimt is the created eBarimt receipt.
type Ebarimt struct {
	EbarimtID      string `json:"id"`
	PaymentID      string `json:"payment_id"`
	EbarimtType    string `json:"ebarimt_type"`
	EbarimtReceipt string `json:"ebarimt_receipt"` // PDF / JSON payload
	EbarimtQR      string `json:"ebarimt_qr,omitempty"`
}

// Location is a city/district/bank record.
type Location struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// Bank is a supported bank for merchant onboarding.
type Bank struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Logo string `json:"logo,omitempty"`
}

// WebhookPayload is the body QPay POSTs to an invoice callback URL.
type WebhookPayload struct {
	Type      string `json:"type"`
	ObjectID  string `json:"object_id"`
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
}
