package com.fiberman.fiberman_java_backend.dto;

import java.time.Instant;
import java.util.Map;

public record FiberCallResponse(
        String historyId,
        String method,
        String backendPath,
        Instant timestamp,
        boolean success,
        Map<String, Object> params,
        Object result,
        Map<String, Object> error,
        CodeArtifacts codeArtifacts,
        InvoiceQrCodeResponse invoiceQrCode
) {
}
