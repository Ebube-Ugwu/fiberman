package com.fiberman.fiberman_java_backend.dto;

import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotBlank;

public record SendPaymentApiRequest(
        @NotBlank String invoice,
        @Min(1) Long amount,
        @Min(1) Long timeoutSeconds
) {
}
