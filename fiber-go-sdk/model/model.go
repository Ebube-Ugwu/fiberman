package model

type CreateInvoiceRequest struct {
	Amount        *int64  `json:"amount,omitempty"`
	Currency      string  `json:"currency"`
	Description   *string `json:"description,omitempty"`
	ExpirySeconds *int64  `json:"expiry,omitempty"`
}

type SendPaymentRequest struct {
	Invoice        string `json:"invoice"`
	Amount         *int64 `json:"amount,omitempty"`
	TimeoutSeconds *int64 `json:"timeout,omitempty"`
}

type JSONRPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      string `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
}

type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      any           `json:"id"`
	Result  any           `json:"result"`
	Error   *JSONRPCError `json:"error"`
}

type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
