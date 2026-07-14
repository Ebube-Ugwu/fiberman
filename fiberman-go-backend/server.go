package fiberman

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fiberman/fiber-go-sdk/client"
	"github.com/fiberman/fiber-go-sdk/fibererr"
	"github.com/fiberman/fiber-go-sdk/model"
)

type Server struct {
	config        Config
	settings      *RuntimeSettings
	sessions      *SessionStore
	codeArtifacts *CodeArtifactService
	qrCodes       *QRCodeService
}

func NewServer(config Config) *Server {
	settings := NewRuntimeSettings(config)
	return &Server{
		config:        config,
		settings:      settings,
		sessions:      NewSessionStore(),
		codeArtifacts: NewCodeArtifactService(settings),
		qrCodes:       NewQRCodeService(),
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/fiber/node-info", s.withSession(s.handleNodeInfo))
	mux.HandleFunc("GET /api/fiber/channels", s.withSession(s.handleListChannels))
	mux.HandleFunc("POST /api/fiber/channels/details", s.withSession(s.handleGetChannel))
	mux.HandleFunc("GET /api/fiber/peers", s.withSession(s.handleListPeers))
	mux.HandleFunc("POST /api/fiber/invoices", s.withSession(s.handleCreateInvoice))
	mux.HandleFunc("POST /api/fiber/invoices/qr", s.withSession(s.handleInvoiceQRCode))
	mux.HandleFunc("POST /api/fiber/payments", s.withSession(s.handleSendPayment))
	mux.HandleFunc("POST /api/fiber/payments/status", s.withSession(s.handleGetPayment))
	mux.HandleFunc("GET /api/fiber/history", s.withSession(s.handleHistory))
	mux.HandleFunc("DELETE /api/fiber/history", s.withSession(s.handleClearHistory))
	mux.HandleFunc("GET /api/settings", s.handleGetSettings)
	mux.HandleFunc("PUT /api/settings", s.handleUpdateSettings)
	return withCORS(mux)
}

func (s *Server) withSession(next func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := s.sessions.SessionID(w, r)
		next(w, r, sessionID)
	}
}

func (s *Server) handleNodeInfo(w http.ResponseWriter, r *http.Request, sessionID string) {
	s.execute(w, sessionID, "/api/fiber/node-info", "node_info", map[string]any{}, "sdk.NodeInfo()", func(c *client.Client) (any, error) {
		return c.NodeInfo()
	})
}

func (s *Server) handleListChannels(w http.ResponseWriter, r *http.Request, sessionID string) {
	s.execute(w, sessionID, "/api/fiber/channels", "list_channels", map[string]any{}, "sdk.ListChannels()", func(c *client.Client) (any, error) {
		return c.ListChannels()
	})
}

func (s *Server) handleGetChannel(w http.ResponseWriter, r *http.Request, sessionID string) {
	parsed, ok := decodeJSON[GetChannelAPIRequest](w, r)
	if !ok {
		return
	}
	if strings.TrimSpace(parsed.ChannelID) == "" {
		writeValidationError(w, map[string]string{"channelId": "must not be blank"})
		return
	}

	s.execute(w, sessionID, "/api/fiber/channels/details", "get_channel", parsed, fmt.Sprintf("sdk.GetChannel(%q)", parsed.ChannelID), func(c *client.Client) (any, error) {
		return c.GetChannel(parsed.ChannelID)
	})
}

func (s *Server) handleListPeers(w http.ResponseWriter, r *http.Request, sessionID string) {
	s.execute(w, sessionID, "/api/fiber/peers", "list_peers", map[string]any{}, "sdk.ListPeers()", func(c *client.Client) (any, error) {
		return c.ListPeers()
	})
}

func (s *Server) handleCreateInvoice(w http.ResponseWriter, r *http.Request, sessionID string) {
	parsed, ok := decodeJSON[CreateInvoiceAPIRequest](w, r)
	if !ok {
		return
	}

	validation := map[string]string{}
	if parsed.Amount == nil || *parsed.Amount < 1 {
		validation["amount"] = "must be greater than or equal to 1"
	}
	if strings.TrimSpace(parsed.Currency) == "" {
		validation["currency"] = "must not be blank"
	}
	if parsed.ExpirySeconds != nil && *parsed.ExpirySeconds < 1 {
		validation["expirySeconds"] = "must be greater than or equal to 1"
	}
	if len(validation) > 0 {
		writeValidationError(w, validation)
		return
	}

	request := model.CreateInvoiceRequest{
		Amount:        parsed.Amount,
		Currency:      parsed.Currency,
		Description:   parsed.Description,
		ExpirySeconds: parsed.ExpirySeconds,
	}
	goCall := fmt.Sprintf("sdk.CreateInvoice(model.CreateInvoiceRequest{Amount: int64Ptr(%d), Currency: %q, Description: stringPtr(%s), ExpirySeconds: int64PtrOrNil(%s)})", derefInt64(parsed.Amount), parsed.Currency, quoteNullableString(parsed.Description), quoteNullableInt64(parsed.ExpirySeconds))
	s.execute(w, sessionID, "/api/fiber/invoices", "new_invoice", parsed, goCall, func(c *client.Client) (any, error) {
		return c.CreateInvoice(request)
	})
}

func (s *Server) handleInvoiceQRCode(w http.ResponseWriter, r *http.Request, sessionID string) {
	parsed, ok := decodeJSON[InvoiceQRCodeRequest](w, r)
	if !ok {
		return
	}
	if strings.TrimSpace(parsed.Invoice) == "" {
		writeValidationError(w, map[string]string{"invoice": "must not be blank"})
		return
	}
	if parsed.Size != nil && *parsed.Size < 64 {
		writeValidationError(w, map[string]string{"size": "must be greater than or equal to 64"})
		return
	}

	qrCode, err := s.qrCodes.Generate(parsed.Invoice, parsed.Size)
	if err != nil {
		writeAppError(w, err)
		return
	}

	params := map[string]any{
		"invoice": parsed.Invoice,
		"size":    parsed.Size,
	}
	response := FiberCallResponse{
		HistoryID:     randomID(),
		Method:        "invoice_qr",
		BackendPath:   "/api/fiber/invoices/qr",
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		Success:       true,
		Params:        params,
		CodeArtifacts: nil,
		InvoiceQRCode: qrCode,
	}
	s.sessions.Append(sessionID, response)
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) handleSendPayment(w http.ResponseWriter, r *http.Request, sessionID string) {
	parsed, ok := decodeJSON[SendPaymentAPIRequest](w, r)
	if !ok {
		return
	}

	validation := map[string]string{}
	if strings.TrimSpace(parsed.Invoice) == "" {
		validation["invoice"] = "must not be blank"
	}
	if parsed.Amount != nil && *parsed.Amount < 1 {
		validation["amount"] = "must be greater than or equal to 1"
	}
	if parsed.TimeoutSeconds != nil && *parsed.TimeoutSeconds < 1 {
		validation["timeoutSeconds"] = "must be greater than or equal to 1"
	}
	if len(validation) > 0 {
		writeValidationError(w, validation)
		return
	}

	request := model.SendPaymentRequest{
		Invoice:        parsed.Invoice,
		Amount:         parsed.Amount,
		TimeoutSeconds: parsed.TimeoutSeconds,
	}
	goCall := fmt.Sprintf("sdk.SendPayment(model.SendPaymentRequest{Invoice: %q, Amount: int64PtrOrNil(%s), TimeoutSeconds: int64PtrOrNil(%s)})", parsed.Invoice, quoteNullableInt64(parsed.Amount), quoteNullableInt64(parsed.TimeoutSeconds))
	s.execute(w, sessionID, "/api/fiber/payments", "send_payment", parsed, goCall, func(c *client.Client) (any, error) {
		return c.SendPayment(request)
	})
}

func (s *Server) handleGetPayment(w http.ResponseWriter, r *http.Request, sessionID string) {
	parsed, ok := decodeJSON[GetPaymentStatusAPIRequest](w, r)
	if !ok {
		return
	}
	if strings.TrimSpace(parsed.PaymentID) == "" {
		writeValidationError(w, map[string]string{"paymentId": "must not be blank"})
		return
	}

	s.execute(w, sessionID, "/api/fiber/payments/status", "get_payment", parsed, fmt.Sprintf("sdk.GetPayment(%q)", parsed.PaymentID), func(c *client.Client) (any, error) {
		return c.GetPayment(parsed.PaymentID)
	})
}

func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request, sessionID string) {
	writeJSON(w, http.StatusOK, s.sessions.History(sessionID))
}

func (s *Server) handleClearHistory(w http.ResponseWriter, r *http.Request, sessionID string) {
	s.sessions.Clear(sessionID)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleGetSettings(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.settings.Get())
}

func (s *Server) handleUpdateSettings(w http.ResponseWriter, r *http.Request) {
	parsed, ok := decodeJSON[UpdatePlaygroundSettingsRequest](w, r)
	if !ok {
		return
	}

	validation := map[string]string{}
	if strings.TrimSpace(parsed.NodeURL) == "" {
		validation["nodeUrl"] = "must not be blank"
	}
	if parsed.TimeoutSeconds < 1 {
		validation["timeoutSeconds"] = "must be greater than or equal to 1"
	}
	if len(validation) > 0 {
		writeValidationError(w, validation)
		return
	}

	writeJSON(w, http.StatusOK, s.settings.Update(parsed))
}

func (s *Server) execute(
	w http.ResponseWriter,
	sessionID string,
	backendPath string,
	method string,
	requestBody any,
	goMethodCall string,
	call func(*client.Client) (any, error),
) {
	fiberClient, err := s.settings.Client()
	if err != nil {
		writeAppError(w, err)
		return
	}

	var artifactRequest any = requestBody
	if emptyMap, ok := requestBody.(map[string]any); ok && len(emptyMap) == 0 {
		artifactRequest = nil
	}
	codeArtifacts := s.codeArtifacts.Generate(backendPath, goMethodCall, artifactRequest)
	params := toMap(requestBody)

	result, err := call(fiberClient)
	if err != nil {
		response := FiberCallResponse{
			HistoryID:     randomID(),
			Method:        method,
			BackendPath:   backendPath,
			Timestamp:     time.Now().UTC().Format(time.RFC3339),
			Success:       false,
			Params:        params,
			Error:         toErrorPayload(err),
			CodeArtifacts: codeArtifacts,
		}
		s.sessions.Append(sessionID, response)
		writeMappedFiberError(w, err)
		return
	}

	response := FiberCallResponse{
		HistoryID:     randomID(),
		Method:        method,
		BackendPath:   backendPath,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		Success:       true,
		Params:        params,
		Result:        result,
		CodeArtifacts: codeArtifacts,
		InvoiceQRCode: s.extractInvoiceQRCode(result),
	}
	s.sessions.Append(sessionID, response)
	writeJSON(w, http.StatusOK, response)
}

func (s *Server) extractInvoiceQRCode(result any) *InvoiceQRCodeResponse {
	resultMap, ok := result.(map[string]any)
	if !ok {
		return nil
	}

	invoiceValue, ok := resultMap["invoice_address"]
	if !ok || invoiceValue == nil {
		invoiceValue = resultMap["invoice"]
	}

	invoice, ok := invoiceValue.(string)
	if !ok || strings.TrimSpace(invoice) == "" {
		return nil
	}

	qrCode, err := s.qrCodes.Generate(invoice, nil)
	if err != nil {
		return nil
	}
	return qrCode
}

func toMap(value any) map[string]any {
	if value == nil {
		return map[string]any{}
	}
	if typed, ok := value.(map[string]any); ok {
		return typed
	}
	bytes, err := json.Marshal(value)
	if err != nil {
		return map[string]any{}
	}
	var result map[string]any
	if err := json.Unmarshal(bytes, &result); err != nil {
		return map[string]any{}
	}
	return result
}

func toErrorPayload(err error) map[string]any {
	var rpcErr *fibererr.RPCError
	if errors.As(err, &rpcErr) {
		return map[string]any{
			"type":    "fiber_rpc_error",
			"message": rpcErr.Error(),
			"code":    rpcErr.Code,
			"data":    rpcErr.Data,
		}
	}

	var httpErr *fibererr.HTTPError
	if errors.As(err, &httpErr) {
		return map[string]any{
			"type":         "fiber_http_error",
			"message":      httpErr.Error(),
			"status":       httpErr.StatusCode,
			"responseBody": httpErr.ResponseBody,
		}
	}

	var timeoutErr *fibererr.TimeoutError
	if errors.As(err, &timeoutErr) {
		return map[string]any{
			"type":    "fiber_timeout",
			"message": timeoutErr.Error(),
		}
	}

	var transportErr *fibererr.TransportError
	if errors.As(err, &transportErr) {
		return map[string]any{
			"type":    "fiber_backend_error",
			"message": transportErr.Error(),
		}
	}

	var serializationErr *fibererr.SerializationError
	if errors.As(err, &serializationErr) {
		return map[string]any{
			"type":    "fiber_backend_error",
			"message": serializationErr.Error(),
		}
	}

	return map[string]any{
		"type":    "internal_error",
		"message": err.Error(),
	}
}

func writeMappedFiberError(w http.ResponseWriter, err error) {
	var rpcErr *fibererr.RPCError
	if errors.As(err, &rpcErr) {
		writeJSONWithStatus(w, http.StatusBadRequest, map[string]any{
			"error":   "fiber_rpc_error",
			"message": rpcErr.Error(),
			"code":    rpcErr.Code,
			"data":    rpcErr.Data,
		})
		return
	}

	var httpErr *fibererr.HTTPError
	if errors.As(err, &httpErr) {
		writeJSONWithStatus(w, http.StatusBadGateway, map[string]any{
			"error":        "fiber_http_error",
			"message":      httpErr.Error(),
			"status":       httpErr.StatusCode,
			"responseBody": httpErr.ResponseBody,
		})
		return
	}

	var timeoutErr *fibererr.TimeoutError
	if errors.As(err, &timeoutErr) {
		writeJSONWithStatus(w, http.StatusGatewayTimeout, map[string]any{
			"error":   "fiber_timeout",
			"message": timeoutErr.Error(),
		})
		return
	}

	var transportErr *fibererr.TransportError
	if errors.As(err, &transportErr) {
		writeJSONWithStatus(w, http.StatusBadGateway, map[string]any{
			"error":   "fiber_backend_error",
			"message": transportErr.Error(),
		})
		return
	}

	var serializationErr *fibererr.SerializationError
	if errors.As(err, &serializationErr) {
		writeJSONWithStatus(w, http.StatusBadGateway, map[string]any{
			"error":   "fiber_backend_error",
			"message": serializationErr.Error(),
		})
		return
	}

	writeAppError(w, err)
}

func writeAppError(w http.ResponseWriter, err error) {
	writeJSONWithStatus(w, http.StatusInternalServerError, map[string]any{
		"error":   "internal_error",
		"message": err.Error(),
	})
}

func writeValidationError(w http.ResponseWriter, details map[string]string) {
	writeJSONWithStatus(w, http.StatusBadRequest, map[string]any{
		"error":   "validation_failed",
		"details": details,
	})
}

func decodeJSON[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
	var value T
	if err := json.NewDecoder(r.Body).Decode(&value); err != nil {
		writeJSONWithStatus(w, http.StatusBadRequest, map[string]any{
			"error":   "validation_failed",
			"details": map[string]string{"body": "invalid JSON body"},
		})
		return value, false
	}
	return value, true
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	writeJSONWithStatus(w, status, value)
}

func writeJSONWithStatus(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func quoteNullableString(value *string) string {
	if value == nil {
		return "nil"
	}
	return fmt.Sprintf("%q", *value)
}

func quoteNullableInt64(value *int64) string {
	if value == nil {
		return "nil"
	}
	return fmt.Sprintf("%d", *value)
}

func derefInt64(value *int64) int64 {
	if value == nil {
		return 0
	}
	return *value
}
