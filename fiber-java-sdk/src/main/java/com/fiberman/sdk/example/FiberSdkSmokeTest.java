package com.fiberman.sdk.example;

import com.fasterxml.jackson.databind.JsonNode;
import com.fiberman.sdk.client.FiberClient;
import com.fiberman.sdk.model.CreateInvoiceRequest;
import com.fiberman.sdk.model.SendPaymentRequest;
import java.time.Duration;

public final class FiberSdkSmokeTest {
    private FiberSdkSmokeTest() {
    }

    public static void main(String[] args) {
        String baseUrl = requiredSetting("FIBER_NODE_URL");
        String authToken = setting("FIBER_NODE_AUTH_TOKEN");
        Duration timeout = Duration.ofSeconds(longSetting("FIBER_NODE_TIMEOUT_SECONDS", 30L));

        FiberClient.Builder builder = FiberClient.builder()
                .baseUrl(baseUrl)
                .requestTimeout(timeout)
                .connectTimeout(timeout);

        if (authToken != null && !authToken.isBlank()) {
            builder.authToken(authToken);
        }

        FiberClient client = builder.build();

        runStep("nodeInfo", client.nodeInfo());
        runStep("listChannels", client.listChannels());

        String invoiceAmount = setting("FIBER_TEST_INVOICE_AMOUNT");
        String invoiceCurrency = setting("FIBER_TEST_INVOICE_CURRENCY");
        String invoiceDescription = setting("FIBER_TEST_INVOICE_DESCRIPTION");
        if (invoiceAmount != null && !invoiceAmount.isBlank()) {
            if (invoiceCurrency == null || invoiceCurrency.isBlank()) {
                throw new IllegalStateException("Missing required environment variable: FIBER_TEST_INVOICE_CURRENCY");
            }
            runStep("createInvoice", client.createInvoice(new CreateInvoiceRequest(
                    Long.parseLong(invoiceAmount),
                    invoiceCurrency,
                    blankToNull(invoiceDescription),
                    longSettingNullable("FIBER_TEST_INVOICE_EXPIRY_SECONDS"))));
        } else {
            System.out.println("Skipping createInvoice: set FIBER_TEST_INVOICE_AMOUNT and FIBER_TEST_INVOICE_CURRENCY to enable it.");
        }

        String paymentInvoice = setting("FIBER_TEST_PAYMENT_INVOICE");
        if (paymentInvoice != null && !paymentInvoice.isBlank()) {
            runStep("sendPayment", client.sendPayment(new SendPaymentRequest(
                    paymentInvoice,
                    longSettingNullable("FIBER_TEST_PAYMENT_AMOUNT"),
                    longSettingNullable("FIBER_TEST_PAYMENT_TIMEOUT_SECONDS"))));
        } else {
            System.out.println("Skipping sendPayment: set FIBER_TEST_PAYMENT_INVOICE to enable it.");
        }
    }

    private static void runStep(String methodName, JsonNode result) {
        System.out.println("=== " + methodName + " SUCCESS ===");
        System.out.println(result.toPrettyString());
    }

    private static String requiredSetting(String name) {
        String value = setting(name);
        if (value == null || value.isBlank()) {
            throw new IllegalStateException("Missing required setting: " + name + " (env var or -D" + name + "=...)");
        }
        return value;
    }

    private static String setting(String name) {
        String propertyValue = System.getProperty(name);
        if (propertyValue != null && !propertyValue.isBlank()) {
            return propertyValue;
        }
        return System.getenv(name);
    }

    private static long longSetting(String name, long defaultValue) {
        String value = setting(name);
        return value == null || value.isBlank() ? defaultValue : Long.parseLong(value);
    }

    private static Long longSettingNullable(String name) {
        String value = setting(name);
        return value == null || value.isBlank() ? null : Long.parseLong(value);
    }

    private static String blankToNull(String value) {
        return value == null || value.isBlank() ? null : value;
    }
}
