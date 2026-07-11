package com.fiberman.fiberman_java_backend.dto;

public record InvoiceQrCodeResponse(
        String value,
        int size,
        String pngBase64,
        String dataUrl
) {
}
