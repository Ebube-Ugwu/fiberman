package com.fiberman.fiberman_java_backend.dto;

import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotBlank;

public record UpdatePlaygroundSettingsRequest(
        @NotBlank String nodeUrl,
        String authToken,
        @Min(1) long timeoutSeconds,
        String defaultInvoiceCurrency
) {
}
