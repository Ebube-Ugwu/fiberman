package com.fiberman.fiberman_java_backend.dto;

public record PlaygroundSettingsResponse(
        String nodeUrl,
        String authToken,
        long timeoutSeconds,
        String defaultInvoiceCurrency,
        String playgroundBaseUrl
) {
}
