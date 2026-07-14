package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/fiberman/fiber-go-sdk/fibererr"
	"github.com/fiberman/fiber-go-sdk/model"
)

func TestNodeInfoUsesJSONRPCAndAuthHeader(t *testing.T) {
	var requestBody string
	var authHeader string

	sdk := mustClientWithRoundTripper(t, roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		authHeader = request.Header.Get("Authorization")
		body, err := io.ReadAll(request.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		requestBody = string(body)
		return jsonResponse(http.StatusOK, `{"jsonrpc":"2.0","id":"1","result":{"version":"0.4.0","node_name":"local-node"}}`), nil
	}))

	result, err := sdk.NodeInfo()
	if err != nil {
		t.Fatalf("node info failed: %v", err)
	}

	requestJSON := decodeMap(t, requestBody)
	if got := requestJSON["jsonrpc"]; got != "2.0" {
		t.Fatalf("expected jsonrpc 2.0, got %#v", got)
	}
	if got := requestJSON["method"]; got != "node_info" {
		t.Fatalf("expected method node_info, got %#v", got)
	}
	if id, _ := requestJSON["id"].(string); len(id) <= 10 {
		t.Fatalf("expected random id, got %#v", requestJSON["id"])
	}
	if authHeader != "Bearer secret-token" {
		t.Fatalf("expected auth header, got %q", authHeader)
	}

	resultMap := asMap(t, result)
	if got := resultMap["version"]; got != "0.4.0" {
		t.Fatalf("expected version 0.4.0, got %#v", got)
	}
	if got := resultMap["node_name"]; got != "local-node" {
		t.Fatalf("expected node_name local-node, got %#v", got)
	}
}

func TestNodeInfoReturnsActualJSONPayloadWhenUsingDefaultDecoding(t *testing.T) {
	sdk := mustClientWithRoundTripper(t, roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, `{"jsonrpc":"2.0","id":"1","result":{"version":"0.4.0","features":["rpc","payments"]}}`), nil
	}))

	result, err := sdk.NodeInfo()
	if err != nil {
		t.Fatalf("node info failed: %v", err)
	}

	resultMap := asMap(t, result)
	if got := resultMap["version"]; got != "0.4.0" {
		t.Fatalf("expected version 0.4.0, got %#v", got)
	}
	features := asSlice(t, resultMap["features"])
	if got := features[0]; got != "rpc" {
		t.Fatalf("expected first feature rpc, got %#v", got)
	}
}

func TestCreateInvoiceMapsToNewInvoiceAndSerializesFields(t *testing.T) {
	var requestBody string

	sdk := mustClientWithRoundTripper(t, roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		requestBody = string(body)
		return jsonResponse(http.StatusOK, `{"jsonrpc":"2.0","id":"1","result":{"invoice_address":"fiber1testinvoice"}}`), nil
	}))

	amount := int64(4200)
	expiry := int64(3600)
	description := "sdk test invoice"

	result, err := sdk.CreateInvoice(model.CreateInvoiceRequest{
		Amount:        &amount,
		Currency:      "FIBD",
		Description:   &description,
		ExpirySeconds: &expiry,
	})
	if err != nil {
		t.Fatalf("create invoice failed: %v", err)
	}

	requestJSON := decodeMap(t, requestBody)
	if got := requestJSON["method"]; got != "new_invoice" {
		t.Fatalf("expected method new_invoice, got %#v", got)
	}

	params := asMap(t, requestJSON["params"])
	nested := asMap(t, params["params"])
	if got := nested["amount"]; got != "4200" {
		t.Fatalf("expected amount \"4200\", got %#v", got)
	}
	if got := nested["currency"]; got != "FIBD" {
		t.Fatalf("expected currency FIBD, got %#v", got)
	}
	if got := nested["description"]; got != description {
		t.Fatalf("expected description %q, got %#v", description, got)
	}
	if got := int64(nested["expiry"].(float64)); got != expiry {
		t.Fatalf("expected expiry %d, got %d", expiry, got)
	}

	resultMap := asMap(t, result)
	if got := resultMap["invoice_address"]; got != "fiber1testinvoice" {
		t.Fatalf("expected invoice_address fiber1testinvoice, got %#v", got)
	}
}

func TestGetChannelMapsChannelIDCorrectly(t *testing.T) {
	var requestBody string

	sdk := mustClientWithRoundTripper(t, roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		requestBody = string(body)
		return jsonResponse(http.StatusOK, `{"jsonrpc":"2.0","id":"1","result":{"channel_id":"abc123","state":"OPEN"}}`), nil
	}))

	result, err := sdk.GetChannel("abc123")
	if err != nil {
		t.Fatalf("get channel failed: %v", err)
	}

	requestJSON := decodeMap(t, requestBody)
	if got := requestJSON["method"]; got != "get_channel" {
		t.Fatalf("expected method get_channel, got %#v", got)
	}
	params := asMap(t, requestJSON["params"])
	nested := asMap(t, params["params"])
	if got := nested["channel_id"]; got != "abc123" {
		t.Fatalf("expected channel_id abc123, got %#v", got)
	}

	resultMap := asMap(t, result)
	if got := resultMap["state"]; got != "OPEN" {
		t.Fatalf("expected state OPEN, got %#v", got)
	}
}

func TestGetPaymentMapsPaymentIDCorrectly(t *testing.T) {
	var requestBody string

	sdk := mustClientWithRoundTripper(t, roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		requestBody = string(body)
		return jsonResponse(http.StatusOK, `{"jsonrpc":"2.0","id":"1","result":{"payment_id":"pay_123","status":"SETTLED"}}`), nil
	}))

	result, err := sdk.GetPayment("pay_123")
	if err != nil {
		t.Fatalf("get payment failed: %v", err)
	}

	requestJSON := decodeMap(t, requestBody)
	if got := requestJSON["method"]; got != "get_payment" {
		t.Fatalf("expected method get_payment, got %#v", got)
	}
	params := asMap(t, requestJSON["params"])
	nested := asMap(t, params["params"])
	if got := nested["payment_id"]; got != "pay_123" {
		t.Fatalf("expected payment_id pay_123, got %#v", got)
	}

	resultMap := asMap(t, result)
	if got := resultMap["status"]; got != "SETTLED" {
		t.Fatalf("expected status SETTLED, got %#v", got)
	}
}

func TestListPeersUsesExpectedMethodName(t *testing.T) {
	var requestBody string

	sdk := mustClientWithRoundTripper(t, roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		requestBody = string(body)
		return jsonResponse(http.StatusOK, `{"jsonrpc":"2.0","id":"1","result":[{"peer_id":"peer-1"}]}`), nil
	}))

	result, err := sdk.ListPeers()
	if err != nil {
		t.Fatalf("list peers failed: %v", err)
	}

	requestJSON := decodeMap(t, requestBody)
	if got := requestJSON["method"]; got != "list_peers" {
		t.Fatalf("expected method list_peers, got %#v", got)
	}
	params := asMap(t, requestJSON["params"])
	if _, ok := params["params"]; !ok {
		t.Fatalf("expected nested params field, got %#v", params)
	}

	resultList := asSlice(t, result)
	first := asMap(t, resultList[0])
	if got := first["peer_id"]; got != "peer-1" {
		t.Fatalf("expected peer_id peer-1, got %#v", got)
	}
}

func TestListChannelsUsesNestedParamsShape(t *testing.T) {
	var requestBody string

	sdk := mustClientWithRoundTripper(t, roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		requestBody = string(body)
		return jsonResponse(http.StatusOK, `{"jsonrpc":"2.0","id":"1","result":[]}`), nil
	}))

	if _, err := sdk.ListChannels(); err != nil {
		t.Fatalf("list channels failed: %v", err)
	}

	requestJSON := decodeMap(t, requestBody)
	if got := requestJSON["method"]; got != "list_channels" {
		t.Fatalf("expected method list_channels, got %#v", got)
	}
	params := asMap(t, requestJSON["params"])
	if _, ok := params["params"]; !ok {
		t.Fatalf("expected nested params field, got %#v", params)
	}
}

func TestInvokeAllowsAdHocMethods(t *testing.T) {
	var requestBody string

	sdk := mustClientWithRoundTripper(t, roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		requestBody = string(body)
		return jsonResponse(http.StatusOK, `{"jsonrpc":"2.0","id":"1","result":{"ok":true}}`), nil
	}))

	result, err := sdk.Invoke("custom_method", map[string]any{"flag": true})
	if err != nil {
		t.Fatalf("invoke failed: %v", err)
	}

	requestJSON := decodeMap(t, requestBody)
	if got := requestJSON["method"]; got != "custom_method" {
		t.Fatalf("expected method custom_method, got %#v", got)
	}
	params := asMap(t, requestJSON["params"])
	if got := params["flag"]; got != true {
		t.Fatalf("expected flag true, got %#v", got)
	}

	resultMap := asMap(t, result)
	if got := resultMap["ok"]; got != true {
		t.Fatalf("expected ok true, got %#v", got)
	}
}

func TestNon2xxResponsesBecomeHTTPError(t *testing.T) {
	sdk := mustClientWithRoundTripper(t, roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusUnauthorized, `{"error":"unauthorized"}`), nil
	}))

	_, err := sdk.NodeInfo()
	if err == nil {
		t.Fatal("expected error")
	}

	var httpErr *fibererr.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected HTTPError, got %T", err)
	}
	if httpErr.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", httpErr.StatusCode)
	}
	if !strings.Contains(httpErr.Error(), "node_info") {
		t.Fatalf("expected method context, got %q", httpErr.Error())
	}
	if !strings.Contains(httpErr.ResponseBody, "unauthorized") {
		t.Fatalf("expected response body, got %q", httpErr.ResponseBody)
	}
}

func TestRPCErrorsBecomeRPCError(t *testing.T) {
	sdk := mustClientWithRoundTripper(t, roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, `{"jsonrpc":"2.0","id":"1","error":{"code":-32001,"message":"bad invoice","data":{"field":"invoice"}}}`), nil
	}))

	_, err := sdk.NodeInfo()
	if err == nil {
		t.Fatal("expected error")
	}

	var rpcErr *fibererr.RPCError
	if !errors.As(err, &rpcErr) {
		t.Fatalf("expected RPCError, got %T", err)
	}
	if rpcErr.Code != -32001 {
		t.Fatalf("expected code -32001, got %d", rpcErr.Code)
	}
	if !strings.Contains(rpcErr.Error(), "node_info") {
		t.Fatalf("expected method context, got %q", rpcErr.Error())
	}
	data := asMap(t, rpcErr.Data)
	if got := data["field"]; got != "invoice" {
		t.Fatalf("expected field invoice, got %#v", got)
	}
}

func TestMalformedJSONBecomesSerializationError(t *testing.T) {
	sdk := mustClientWithRoundTripper(t, roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		return jsonResponse(http.StatusOK, `not-json`), nil
	}))

	_, err := sdk.NodeInfo()
	if err == nil {
		t.Fatal("expected error")
	}

	var serializationErr *fibererr.SerializationError
	if !errors.As(err, &serializationErr) {
		t.Fatalf("expected SerializationError, got %T", err)
	}
	if !strings.Contains(serializationErr.Error(), "node_info") {
		t.Fatalf("expected method context, got %q", serializationErr.Error())
	}
}

func TestRequestTimeoutBecomesTimeoutError(t *testing.T) {
	sdk := mustClientWithRoundTripper(t, roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		select {
		case <-time.After(200 * time.Millisecond):
			return jsonResponse(http.StatusOK, `{"jsonrpc":"2.0","id":"1","result":{"status":"late"}}`), nil
		case <-request.Context().Done():
			return nil, request.Context().Err()
		}
	}), 50*time.Millisecond)

	_, err := sdk.NodeInfo()
	if err == nil {
		t.Fatal("expected error")
	}

	var timeoutErr *fibererr.TimeoutError
	if !errors.As(err, &timeoutErr) {
		t.Fatalf("expected TimeoutError, got %T", err)
	}
	if !strings.Contains(timeoutErr.Error(), "node_info") {
		t.Fatalf("expected method context, got %q", timeoutErr.Error())
	}
	if !strings.Contains(timeoutErr.Error(), "50ms") {
		t.Fatalf("expected timeout duration, got %q", timeoutErr.Error())
	}
}

func TestMissingBaseURLFailsFast(t *testing.T) {
	_, err := New(Config{})
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "baseURL must be provided" {
		t.Fatalf("expected baseURL error, got %q", err.Error())
	}
}

func TestNewRejectsInvalidBaseURL(t *testing.T) {
	_, err := New(Config{BaseURL: "://bad url"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "baseURL must be a valid URL") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTimeoutErrorWrapsOriginalCause(t *testing.T) {
	sdk := mustClientWithRoundTripper(t, roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		<-request.Context().Done()
		return nil, request.Context().Err()
	}), 10*time.Millisecond)

	_, err := sdk.NodeInfo()
	if err == nil {
		t.Fatal("expected error")
	}
	if errors.Unwrap(err) == nil {
		t.Fatalf("expected wrapped cause, got nil")
	}
}

func TestClientAllowsCustomHeaders(t *testing.T) {
	var gotHeader string

	sdk := mustClientWithRoundTripper(t, roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		gotHeader = request.Header.Get("X-Test-Header")
		return jsonResponse(http.StatusOK, `{"jsonrpc":"2.0","id":"1","result":{"ok":true}}`), nil
	}))

	sdkWithHeaders, err := New(Config{
		BaseURL: "http://fiber.test/",
		Headers: map[string]string{"X-Test-Header": "present"},
		HTTPClient: &http.Client{
			Transport: roundTripperFunc(func(request *http.Request) (*http.Response, error) {
				gotHeader = request.Header.Get("X-Test-Header")
				return jsonResponse(http.StatusOK, `{"jsonrpc":"2.0","id":"1","result":{"ok":true}}`), nil
			}),
		},
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	if _, err := sdkWithHeaders.NodeInfo(); err != nil {
		t.Fatalf("node info failed: %v", err)
	}
	if gotHeader != "present" {
		t.Fatalf("expected custom header, got %q", gotHeader)
	}

	if _, err := sdk.NodeInfo(); err != nil {
		t.Fatalf("node info failed: %v", err)
	}
}

func ExampleClient_NodeInfo() {
	sdk, _ := New(Config{
		BaseURL: "http://fiber.test/",
		HTTPClient: &http.Client{
			Transport: roundTripperFunc(func(request *http.Request) (*http.Response, error) {
				return jsonResponse(http.StatusOK, `{"jsonrpc":"2.0","id":"1","result":{"version":"0.4.0"}}`), nil
			}),
		},
	})

	result, _ := sdk.NodeInfo()
	payload, _ := json.Marshal(result)
	fmt.Println(string(payload))
	// Output: {"version":"0.4.0"}
}

func mustClientWithRoundTripper(t *testing.T, roundTripper http.RoundTripper, requestTimeout ...time.Duration) *Client {
	t.Helper()

	timeout := time.Second
	if len(requestTimeout) > 0 {
		timeout = requestTimeout[0]
	}

	sdk, err := New(Config{
		BaseURL:        "http://fiber.test/",
		RequestTimeout: timeout,
		HTTPClient: &http.Client{
			Transport: roundTripper,
		},
		AuthToken: "secret-token",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return sdk
}

func decodeMap(t *testing.T, raw string) map[string]any {
	t.Helper()
	var value map[string]any
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		t.Fatalf("failed to decode map: %v", err)
	}
	return value
}

func asMap(t *testing.T, value any) map[string]any {
	t.Helper()
	typed, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", value)
	}
	return typed
}

func asSlice(t *testing.T, value any) []any {
	t.Helper()
	typed, ok := value.([]any)
	if !ok {
		t.Fatalf("expected []any, got %T", value)
	}
	return typed
}

func jsonResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return fn(request)
}

type timeoutNetError struct{}

func (timeoutNetError) Error() string   { return "timeout" }
func (timeoutNetError) Timeout() bool   { return true }
func (timeoutNetError) Temporary() bool { return true }

func TestNetTimeoutAlsoBecomesTimeoutError(t *testing.T) {
	sdk := mustClientWithRoundTripper(t, roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		return nil, timeoutNetError{}
	}))

	_, err := sdk.NodeInfo()
	if err == nil {
		t.Fatal("expected error")
	}

	var timeoutErr *fibererr.TimeoutError
	if !errors.As(err, &timeoutErr) {
		t.Fatalf("expected TimeoutError, got %T", err)
	}
}

func TestTransportFailuresBecomeTransportError(t *testing.T) {
	sdk := mustClientWithRoundTripper(t, roundTripperFunc(func(request *http.Request) (*http.Response, error) {
		return nil, net.UnknownNetworkError("unknown")
	}))

	_, err := sdk.NodeInfo()
	if err == nil {
		t.Fatal("expected error")
	}

	var transportErr *fibererr.TransportError
	if !errors.As(err, &transportErr) {
		t.Fatalf("expected TransportError, got %T", err)
	}
	if !strings.Contains(transportErr.Error(), "node_info") {
		t.Fatalf("expected method context, got %q", transportErr.Error())
	}
}
