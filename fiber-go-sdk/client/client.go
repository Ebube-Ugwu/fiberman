package client

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/fiberman/fiber-go-sdk/model"
	"github.com/fiberman/fiber-go-sdk/transport"
)

type Client struct {
	transport *transport.Transport
}

type Config struct {
	BaseURL        string
	AuthToken      string
	Headers        map[string]string
	ConnectTimeout time.Duration
	RequestTimeout time.Duration
	HTTPClient     *http.Client
}

func New(config Config) (*Client, error) {
	baseURL := strings.TrimSpace(config.BaseURL)
	if baseURL == "" {
		return nil, errors.New("baseURL must be provided")
	}

	if _, err := url.Parse(baseURL); err != nil {
		return nil, fmt.Errorf("baseURL must be a valid URL: %w", err)
	}

	connectTimeout := config.ConnectTimeout
	if connectTimeout <= 0 {
		connectTimeout = 10 * time.Second
	}

	requestTimeout := config.RequestTimeout
	if requestTimeout <= 0 {
		requestTimeout = 30 * time.Second
	}

	headers := cloneHeaders(config.Headers)
	if strings.TrimSpace(config.AuthToken) != "" {
		headers["Authorization"] = "Bearer " + strings.TrimSpace(config.AuthToken)
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout: connectTimeout,
				}).DialContext,
				ForceAttemptHTTP2: true,
			},
		}
	}

	return &Client{
		transport: transport.New(httpClient, baseURL, requestTimeout, headers),
	}, nil
}

func (c *Client) NodeInfo() (any, error) {
	return c.Call("node_info", map[string]any{})
}

func (c *Client) CreateInvoice(request model.CreateInvoiceRequest) (any, error) {
	return c.Call("new_invoice", wrapParams(createInvoiceParams(request)))
}

func (c *Client) SendPayment(request model.SendPaymentRequest) (any, error) {
	return c.Call("send_payment", wrapParams(sendPaymentParams(request)))
}

func (c *Client) ListChannels() (any, error) {
	return c.Call("list_channels", wrapParams(map[string]any{}))
}

func (c *Client) ListPeers() (any, error) {
	return c.Call("list_peers", wrapParams(map[string]any{}))
}

func (c *Client) GetChannel(channelID string) (any, error) {
	return c.Call("get_channel", wrapParams(map[string]any{"channel_id": channelID}))
}

func (c *Client) GetPayment(paymentID string) (any, error) {
	return c.Call("get_payment", wrapParams(map[string]any{"payment_id": paymentID}))
}

func (c *Client) Invoke(method string, params any) (any, error) {
	return c.Call(method, params)
}

func (c *Client) Call(method string, params any) (any, error) {
	return c.transport.Call(method, params)
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

func wrapParams(params any) map[string]any {
	return map[string]any{"params": params}
}

func createInvoiceParams(request model.CreateInvoiceRequest) map[string]any {
	params := map[string]any{
		"currency": request.Currency,
	}
	if request.Amount != nil {
		params["amount"] = strconv.FormatInt(*request.Amount, 10)
	}
	if request.Description != nil {
		params["description"] = *request.Description
	}
	if request.ExpirySeconds != nil {
		params["expiry"] = *request.ExpirySeconds
	}
	return params
}

func sendPaymentParams(request model.SendPaymentRequest) map[string]any {
	params := map[string]any{
		"invoice": request.Invoice,
	}
	if request.Amount != nil {
		params["amount"] = strconv.FormatInt(*request.Amount, 10)
	}
	if request.TimeoutSeconds != nil {
		params["timeout"] = *request.TimeoutSeconds
	}
	return params
}
