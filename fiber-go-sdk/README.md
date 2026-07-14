# fiber-go-sdk

`fiber-go-sdk` is a plain Go SDK for talking to a Fiber node over JSON-RPC.

## Features

- no framework dependency
- Fiber JSON-RPC over HTTP POST
- auth-token header support
- clean transport, timeout, HTTP, RPC, and serialization errors
- typed request models for invoice and payment calls
- raw JSON-compatible results for flexible inspection

## Supported Methods

- `NodeInfo`
- `CreateInvoice`
- `SendPayment`
- `ListChannels`
- `ListPeers`
- `GetChannel`
- `GetPayment`
- `Invoke`

## Example

```go
sdk, err := client.New(client.Config{
    BaseURL:   os.Getenv("FIBER_NODE_URL"),
    AuthToken: os.Getenv("FIBER_NODE_AUTH_TOKEN"),
})
if err != nil {
    log.Fatal(err)
}

result, err := sdk.NodeInfo()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("%#v\n", result)
```

## Smoke Example

Run:

```bash
cd fiber-go-sdk
go run ./example/smoke
```

Environment variables:

- `FIBER_NODE_URL`
- `FIBER_NODE_AUTH_TOKEN`
- `FIBER_NODE_TIMEOUT_SECONDS`
- `FIBER_TEST_INVOICE_AMOUNT`
- `FIBER_TEST_INVOICE_CURRENCY`
- `FIBER_TEST_INVOICE_DESCRIPTION`
- `FIBER_TEST_INVOICE_EXPIRY_SECONDS`
- `FIBER_TEST_PAYMENT_INVOICE`
- `FIBER_TEST_PAYMENT_AMOUNT`
- `FIBER_TEST_PAYMENT_TIMEOUT_SECONDS`
