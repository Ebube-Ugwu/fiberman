package com.fiberman.fiberman_java_backend.service;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fiberman.fiberman_java_backend.dto.CreateInvoiceApiRequest;
import com.fiberman.fiberman_java_backend.dto.FiberCallResponse;
import com.fiberman.fiberman_java_backend.dto.GetChannelApiRequest;
import com.fiberman.fiberman_java_backend.dto.GetPaymentStatusApiRequest;
import com.fiberman.fiberman_java_backend.dto.SendPaymentApiRequest;
import com.fiberman.sdk.client.FiberClient;
import com.fiberman.sdk.exception.FiberHttpException;
import com.fiberman.sdk.exception.FiberRpcException;
import com.fiberman.sdk.exception.FiberSerializationException;
import com.fiberman.sdk.exception.FiberTimeoutException;
import com.fiberman.sdk.exception.FiberTransportException;
import com.fiberman.sdk.model.CreateInvoiceRequest;
import com.fiberman.sdk.model.SendPaymentRequest;
import jakarta.servlet.http.HttpSession;
import java.time.Instant;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.UUID;
import java.util.function.Supplier;
import org.springframework.stereotype.Service;

@Service
public class FiberGatewayService {
    private final FiberClient fiberClient;
    private final ObjectMapper objectMapper;
    private final CodeArtifactService codeArtifactService;
    private final FiberHistoryService fiberHistoryService;
    private final InvoiceQrCodeService invoiceQrCodeService;

    public FiberGatewayService(
            FiberClient fiberClient,
            ObjectMapper objectMapper,
            CodeArtifactService codeArtifactService,
            FiberHistoryService fiberHistoryService,
            InvoiceQrCodeService invoiceQrCodeService
    ) {
        this.fiberClient = fiberClient;
        this.objectMapper = objectMapper;
        this.codeArtifactService = codeArtifactService;
        this.fiberHistoryService = fiberHistoryService;
        this.invoiceQrCodeService = invoiceQrCodeService;
    }

    public FiberCallResponse nodeInfo(HttpSession session) {
        return execute(
                session,
                "/api/fiber/node-info",
                "node_info",
                Map.of(),
                "client.nodeInfo()",
                fiberClient::nodeInfo);
    }

    public FiberCallResponse listChannels(HttpSession session) {
        return execute(
                session,
                "/api/fiber/channels",
                "list_channels",
                Map.of(),
                "client.listChannels()",
                fiberClient::listChannels);
    }

    public FiberCallResponse getChannel(HttpSession session, GetChannelApiRequest request) {
        return execute(
                session,
                "/api/fiber/channels/details",
                "get_channel",
                request,
                "client.getChannel(\"%s\")".formatted(escapeJava(request.channelId())),
                () -> fiberClient.getChannel(request.channelId()));
    }

    public FiberCallResponse listPeers(HttpSession session) {
        return execute(
                session,
                "/api/fiber/peers",
                "list_peers",
                Map.of(),
                "client.listPeers()",
                fiberClient::listPeers);
    }

    public FiberCallResponse createInvoice(HttpSession session, CreateInvoiceApiRequest request) {
        CreateInvoiceRequest sdkRequest = new CreateInvoiceRequest(
                request.amount(),
                request.currency(),
                request.description(),
                request.expirySeconds());
        return execute(
                session,
                "/api/fiber/invoices",
                "new_invoice",
                request,
                """
                        client.createInvoice(new CreateInvoiceRequest(
                            %sL,
                            "%s",
                            %s,
                            %s
                        ))
                        """.formatted(
                        request.amount(),
                        escapeJava(request.currency()),
                        quoteNullable(request.description()),
                        request.expirySeconds() == null ? "null" : request.expirySeconds() + "L"),
                () -> fiberClient.createInvoice(sdkRequest));
    }

    public FiberCallResponse sendPayment(HttpSession session, SendPaymentApiRequest request) {
        SendPaymentRequest sdkRequest = new SendPaymentRequest(
                request.invoice(),
                request.amount(),
                request.timeoutSeconds());
        return execute(
                session,
                "/api/fiber/payments",
                "send_payment",
                request,
                """
                        client.sendPayment(new SendPaymentRequest(
                            "%s",
                            %s,
                            %s
                        ))
                        """.formatted(
                        escapeJava(request.invoice()),
                        request.amount() == null ? "null" : request.amount() + "L",
                        request.timeoutSeconds() == null ? "null" : request.timeoutSeconds() + "L"),
                () -> fiberClient.sendPayment(sdkRequest));
    }

    public FiberCallResponse paymentStatus(HttpSession session, GetPaymentStatusApiRequest request) {
        return execute(
                session,
                "/api/fiber/payments/status",
                "get_payment",
                request,
                "client.getPayment(\"%s\")".formatted(escapeJava(request.paymentId())),
                () -> fiberClient.getPayment(request.paymentId()));
    }

    public FiberCallResponse invoiceQrCode(String invoice, Integer size) {
        Map<String, Object> params = new LinkedHashMap<>();
        params.put("invoice", invoice);
        params.put("size", size);
        return new FiberCallResponse(
                UUID.randomUUID().toString(),
                "invoice_qr",
                "/api/fiber/invoices/qr",
                Instant.now(),
                true,
                params,
                null,
                null,
                null,
                invoiceQrCodeService.generate(invoice, size));
    }

    private FiberCallResponse execute(
            HttpSession session,
            String backendPath,
            String method,
            Object requestBody,
            String javaMethodCall,
            Supplier<JsonNode> supplier
    ) {
        Map<String, Object> params = toMap(requestBody);
        var codeArtifacts = codeArtifactService.generate(
                backendPath,
                javaMethodCall,
                requestBody instanceof Map<?, ?> map && map.isEmpty() ? null : requestBody);

        try {
            Object result = unwrap(supplier.get());
            FiberCallResponse response = new FiberCallResponse(
                    UUID.randomUUID().toString(),
                    method,
                    backendPath,
                    Instant.now(),
                    true,
                    params,
                    result,
                    null,
                    codeArtifacts,
                    extractInvoiceQrCode(result));
            return fiberHistoryService.append(session, response);
        } catch (RuntimeException exception) {
            FiberCallResponse response = new FiberCallResponse(
                    UUID.randomUUID().toString(),
                    method,
                    backendPath,
                    Instant.now(),
                    false,
                    params,
                    null,
                    toErrorPayload(exception),
                    codeArtifacts,
                    null);
            fiberHistoryService.append(session, response);
            throw exception;
        }
    }

    private Map<String, Object> toMap(Object requestBody) {
        if (requestBody instanceof Map<?, ?> map) {
            Map<String, Object> converted = new LinkedHashMap<>();
            map.forEach((key, value) -> converted.put(String.valueOf(key), value));
            return converted;
        }
        return objectMapper.convertValue(requestBody, objectMapper.getTypeFactory()
                .constructMapType(LinkedHashMap.class, String.class, Object.class));
    }

    private Object unwrap(JsonNode node) {
        return objectMapper.convertValue(node, Object.class);
    }

    private Map<String, Object> toErrorPayload(RuntimeException exception) {
        if (exception instanceof FiberRpcException fiberRpcException) {
            Map<String, Object> error = new LinkedHashMap<>();
            error.put("type", "fiber_rpc_error");
            error.put("message", fiberRpcException.getMessage());
            error.put("code", fiberRpcException.getCode());
            error.put("data", fiberRpcException.getData() == null ? null : unwrap(fiberRpcException.getData()));
            return error;
        }
        if (exception instanceof FiberHttpException fiberHttpException) {
            Map<String, Object> error = new LinkedHashMap<>();
            error.put("type", "fiber_http_error");
            error.put("message", fiberHttpException.getMessage());
            error.put("status", fiberHttpException.getStatusCode());
            error.put("responseBody", fiberHttpException.getResponseBody());
            return error;
        }
        if (exception instanceof FiberTimeoutException) {
            Map<String, Object> error = new LinkedHashMap<>();
            error.put("type", "fiber_timeout");
            error.put("message", exception.getMessage());
            return error;
        }
        if (exception instanceof FiberTransportException || exception instanceof FiberSerializationException) {
            Map<String, Object> error = new LinkedHashMap<>();
            error.put("type", "fiber_backend_error");
            error.put("message", exception.getMessage());
            return error;
        }
        Map<String, Object> error = new LinkedHashMap<>();
        error.put("type", "internal_error");
        error.put("message", exception.getMessage());
        return error;
    }

    private com.fiberman.fiberman_java_backend.dto.InvoiceQrCodeResponse extractInvoiceQrCode(Object result) {
        if (result instanceof Map<?, ?> resultMap) {
            Object invoiceValue = resultMap.get("invoice_address");
            if (invoiceValue == null) {
                invoiceValue = resultMap.get("invoice");
            }
            if (invoiceValue instanceof String invoice && !invoice.isBlank()) {
                return invoiceQrCodeService.generate(invoice, null);
            }
        }
        return null;
    }

    private String escapeJava(String value) {
        return value.replace("\\", "\\\\").replace("\"", "\\\"");
    }

    private String quoteNullable(String value) {
        return value == null ? "null" : "\"" + escapeJava(value) + "\"";
    }
}
