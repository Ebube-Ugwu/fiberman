package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fiberman/fiber-go-sdk/client"
	"github.com/fiberman/fiber-go-sdk/model"
)

func main() {
	baseURL := requiredSetting("FIBER_NODE_URL")
	authToken := setting("FIBER_NODE_AUTH_TOKEN")
	timeout := time.Duration(longSetting("FIBER_NODE_TIMEOUT_SECONDS", 30)) * time.Second

	sdk, err := client.New(client.Config{
		BaseURL:        baseURL,
		AuthToken:      authToken,
		ConnectTimeout: timeout,
		RequestTimeout: timeout,
	})
	if err != nil {
		log.Fatal(err)
	}

	result, err := sdk.NodeInfo()
	runStep("nodeInfo", result, err)

	result, err = sdk.ListChannels()
	runStep("listChannels", result, err)

	invoiceAmount := setting("FIBER_TEST_INVOICE_AMOUNT")
	invoiceCurrency := setting("FIBER_TEST_INVOICE_CURRENCY")
	invoiceDescription := setting("FIBER_TEST_INVOICE_DESCRIPTION")
	if invoiceAmount != "" {
		if invoiceCurrency == "" {
			log.Fatal("missing required environment variable: FIBER_TEST_INVOICE_CURRENCY")
		}

		amount := mustParseInt64(invoiceAmount)
		result, err = sdk.CreateInvoice(model.CreateInvoiceRequest{
			Amount:        &amount,
			Currency:      invoiceCurrency,
			Description:   stringPtrOrNil(invoiceDescription),
			ExpirySeconds: int64SettingNullable("FIBER_TEST_INVOICE_EXPIRY_SECONDS"),
		})
		runStep("createInvoice", result, err)
	} else {
		fmt.Println("Skipping createInvoice: set FIBER_TEST_INVOICE_AMOUNT and FIBER_TEST_INVOICE_CURRENCY to enable it.")
	}

	paymentInvoice := setting("FIBER_TEST_PAYMENT_INVOICE")
	if paymentInvoice != "" {
		result, err = sdk.SendPayment(model.SendPaymentRequest{
			Invoice:        paymentInvoice,
			Amount:         int64SettingNullable("FIBER_TEST_PAYMENT_AMOUNT"),
			TimeoutSeconds: int64SettingNullable("FIBER_TEST_PAYMENT_TIMEOUT_SECONDS"),
		})
		runStep("sendPayment", result, err)
	} else {
		fmt.Println("Skipping sendPayment: set FIBER_TEST_PAYMENT_INVOICE to enable it.")
	}
}

func runStep(name string, result any, err error) {
	if err != nil {
		log.Fatalf("=== %s FAILURE ===\n%v", name, err)
	}

	output, marshalErr := json.MarshalIndent(result, "", "  ")
	if marshalErr != nil {
		log.Fatalf("failed to format %s result: %v", name, marshalErr)
	}

	fmt.Printf("=== %s SUCCESS ===\n%s\n", name, output)
}

func requiredSetting(name string) string {
	value := setting(name)
	if value == "" {
		log.Fatalf("missing required setting: %s", name)
	}
	return value
}

func setting(name string) string {
	return strings.TrimSpace(os.Getenv(name))
}

func longSetting(name string, defaultValue int64) int64 {
	value := setting(name)
	if value == "" {
		return defaultValue
	}
	return mustParseInt64(value)
}

func int64SettingNullable(name string) *int64 {
	value := setting(name)
	if value == "" {
		return nil
	}
	parsed := mustParseInt64(value)
	return &parsed
}

func mustParseInt64(value string) int64 {
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Fatalf("invalid integer value %q: %v", value, err)
	}
	return parsed
}

func stringPtrOrNil(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
