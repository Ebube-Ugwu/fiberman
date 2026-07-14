package transport

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/fiberman/fiber-go-sdk/fibererr"
	"github.com/fiberman/fiber-go-sdk/model"
)

const contentType = "application/json"

type Transport struct {
	httpClient     *http.Client
	endpoint       string
	requestTimeout time.Duration
	headers        map[string]string
}

func New(httpClient *http.Client, endpoint string, requestTimeout time.Duration, headers map[string]string) *Transport {
	return &Transport{
		httpClient:     httpClient,
		endpoint:       endpoint,
		requestTimeout: requestTimeout,
		headers:        cloneHeaders(headers),
	}
}

func (t *Transport) Call(method string, params any) (any, error) {
	request := model.JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      requestID(),
		Method:  method,
		Params:  params,
	}

	body, err := t.serializeRequest(method, request)
	if err != nil {
		return nil, err
	}

	httpRequest, err := t.buildHTTPRequest(body)
	if err != nil {
		return nil, fibererr.NewTransportError(
			fmt.Sprintf("Failed to construct Fiber RPC request for method '%s'", method),
			err,
		)
	}

	httpResponse, err := t.send(method, httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fibererr.NewTransportError(
			fmt.Sprintf("Transport failure reading Fiber RPC response for method '%s'", method),
			err,
		)
	}

	if httpResponse.StatusCode < 200 || httpResponse.StatusCode >= 300 {
		return nil, fibererr.NewHTTPError(
			fibererr.FormatHTTPMessage(method, httpResponse.StatusCode),
			httpResponse.StatusCode,
			string(responseBody),
		)
	}

	rpcResponse, err := t.deserializeResponse(method, responseBody)
	if err != nil {
		return nil, err
	}

	if rpcResponse.Error != nil {
		return nil, fibererr.NewRPCError(
			fibererr.FormatRPCMessage(method, rpcResponse.Error.Code, rpcResponse.Error.Message),
			rpcResponse.Error.Code,
			rpcResponse.Error.Data,
		)
	}

	return rpcResponse.Result, nil
}

func (t *Transport) serializeRequest(method string, request model.JSONRPCRequest) ([]byte, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fibererr.NewSerializationError(
			fmt.Sprintf("Failed to serialize Fiber RPC request for method '%s'", method),
			err,
		)
	}
	return body, nil
}

func (t *Transport) deserializeResponse(method string, body []byte) (*model.JSONRPCResponse, error) {
	var response model.JSONRPCResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fibererr.NewSerializationError(
			fmt.Sprintf("Failed to deserialize Fiber RPC response for method '%s'", method),
			err,
		)
	}
	return &response, nil
}

func (t *Transport) buildHTTPRequest(body []byte) (*http.Request, error) {
	request, err := http.NewRequest(http.MethodPost, t.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Accept", contentType)
	for name, value := range t.headers {
		request.Header.Set(name, value)
	}

	return request, nil
}

func (t *Transport) send(method string, request *http.Request) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(request.Context(), t.requestTimeout)
	defer cancel()

	response, err := t.httpClient.Do(request.WithContext(ctx))
	if err == nil {
		return response, nil
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return nil, fibererr.NewTimeoutError(
			fmt.Sprintf("Fiber RPC call timed out for method '%s' after %s", method, t.requestTimeout),
			err,
		)
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return nil, fibererr.NewTimeoutError(
			fmt.Sprintf("Fiber RPC call timed out for method '%s' after %s", method, t.requestTimeout),
			err,
		)
	}

	return nil, fibererr.NewTransportError(
		fmt.Sprintf("Transport failure calling Fiber RPC method '%s'", method),
		err,
	)
}

func cloneHeaders(headers map[string]string) map[string]string {
	if len(headers) == 0 {
		return map[string]string{}
	}
	cloned := make(map[string]string, len(headers))
	for key, value := range headers {
		cloned[key] = value
	}
	return cloned
}

func requestID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}
