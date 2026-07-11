package com.fiberman.fiberman_java_backend.dto;

import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotBlank;

public record InvoiceQrCodeRequest(
        @NotBlank String invoice,
        @Min(64) Integer size
) {
}
