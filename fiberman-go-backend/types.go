package fiberman

type CodeArtifacts struct {
	Curl        string `json:"curl"`
	JavaSnippet string `json:"javaSnippet"`
	GoSnippet   string `json:"goSnippet"`
}

type InvoiceQRCodeResponse struct {
	Value     string `json:"value"`
	Size      int    `json:"size"`
	PNGBase64 string `json:"pngBase64"`
	DataURL   string `json:"dataUrl"`
}

type FiberCallResponse struct {
	HistoryID     string                 `json:"historyId"`
	Method        string                 `json:"method"`
	BackendPath   string                 `json:"backendPath"`
	Timestamp     string                 `json:"timestamp"`
	Success       bool                   `json:"success"`
	Params        map[string]any         `json:"params"`
	Result        any                    `json:"result"`
	Error         map[string]any         `json:"error"`
	CodeArtifacts *CodeArtifacts         `json:"codeArtifacts"`
	InvoiceQRCode *InvoiceQRCodeResponse `json:"invoiceQrCode"`
}

type CreateInvoiceAPIRequest struct {
	Amount        *int64  `json:"amount"`
	Currency      string  `json:"currency"`
	Description   *string `json:"description"`
	ExpirySeconds *int64  `json:"expirySeconds"`
}

type SendPaymentAPIRequest struct {
	Invoice        string `json:"invoice"`
	Amount         *int64 `json:"amount"`
	TimeoutSeconds *int64 `json:"timeoutSeconds"`
}

type GetChannelAPIRequest struct {
	ChannelID string `json:"channelId"`
}

type GetPaymentStatusAPIRequest struct {
	PaymentID string `json:"paymentId"`
}

type InvoiceQRCodeRequest struct {
	Invoice string `json:"invoice"`
	Size    *int   `json:"size"`
}

type PlaygroundSettingsResponse struct {
	NodeURL                string `json:"nodeUrl"`
	AuthToken              string `json:"authToken"`
	TimeoutSeconds         int64  `json:"timeoutSeconds"`
	DefaultInvoiceCurrency string `json:"defaultInvoiceCurrency"`
	PlaygroundBaseURL      string `json:"playgroundBaseUrl"`
}

type UpdatePlaygroundSettingsRequest struct {
	NodeURL                string `json:"nodeUrl"`
	AuthToken              string `json:"authToken"`
	TimeoutSeconds         int64  `json:"timeoutSeconds"`
	DefaultInvoiceCurrency string `json:"defaultInvoiceCurrency"`
}
