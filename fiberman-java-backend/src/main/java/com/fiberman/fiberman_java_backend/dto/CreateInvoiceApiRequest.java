package com.fiberman.fiberman_java_backend.dto;

import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotBlank;

public record CreateInvoiceApiRequest(
        @Min(1) Long amount,
        @NotBlank String currency,
        String description,
        @Min(1) Long expirySeconds
) {
}
