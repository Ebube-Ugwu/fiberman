package com.fiberman.sdk.transport;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fiberman.sdk.exception.FiberHttpException;
import com.fiberman.sdk.exception.FiberRpcException;
import com.fiberman.sdk.exception.FiberSerializationException;
import com.fiberman.sdk.exception.FiberTimeoutException;
import com.fiberman.sdk.exception.FiberTransportException;
import com.fiberman.sdk.model.JsonRpcRequest;
import com.fiberman.sdk.model.JsonRpcResponse;
import java.io.IOException;
import java.net.URI;
import java.net.http.HttpTimeoutException;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.Duration;
import java.util.Map;
import java.util.UUID;

public final class FiberTransport {
    private static final String CONTENT_TYPE = "application/json";

    private final HttpClient httpClient;
    private final ObjectMapper objectMapper;
    private final URI endpoint;
    private final Duration requestTimeout;
    private final Map<String, String> headers;

    public FiberTransport(
            HttpClient httpClient,
            ObjectMapper objectMapper,
            URI endpoint,
            Duration requestTimeout,
            Map<String, String> headers
    ) {
        this.httpClient = httpClient;
        this.objectMapper = objectMapper;
        this.endpoint = endpoint;
        this.requestTimeout = requestTimeout;
        this.headers = headers;
    }

    public <T> T call(String method, Object params, Class<T> responseType) {
        JsonRpcRequest request = JsonRpcRequest.of(UUID.randomUUID().toString(), method, params);
        String body = serializeRequest(method, request);
        HttpRequest httpRequest = buildHttpRequest(body);
        HttpResponse<String> httpResponse = send(method, httpRequest);

        if (httpResponse.statusCode() < 200 || httpResponse.statusCode() >= 300) {
            throw new FiberHttpException(
                    "Fiber RPC call failed for method '%s' with HTTP status %d".formatted(method, httpResponse.statusCode()),
                    httpResponse.statusCode(),
                    httpResponse.body());
        }

        JsonRpcResponse rpcResponse = deserializeResponse(method, httpResponse.body());
        if (rpcResponse.error() != null) {
            throw new FiberRpcException(
                    "Fiber RPC error for method '%s': [%d] %s"
                            .formatted(method, rpcResponse.error().code(), rpcResponse.error().message()),
                    rpcResponse.error().code(),
                    rpcResponse.error().data());
        }

        return convertResult(method, rpcResponse.result(), responseType);
    }

    private String serializeRequest(String method, JsonRpcRequest request) {
        try {
            return objectMapper.writeValueAsString(request);
        } catch (JsonProcessingException exception) {
            throw new FiberSerializationException(
                    "Failed to serialize Fiber RPC request for method '%s'".formatted(method),
                    exception);
        }
    }

    private JsonRpcResponse deserializeResponse(String method, String responseBody) {
        try {
            return objectMapper.readValue(responseBody, JsonRpcResponse.class);
        } catch (JsonProcessingException exception) {
            throw new FiberSerializationException(
                    "Failed to deserialize Fiber RPC response for method '%s'".formatted(method),
                    exception);
        }
    }

    private HttpRequest buildHttpRequest(String body) {
        HttpRequest.Builder builder = HttpRequest.newBuilder(endpoint)
                .timeout(requestTimeout)
                .header("Content-Type", CONTENT_TYPE)
                .header("Accept", CONTENT_TYPE)
                .POST(HttpRequest.BodyPublishers.ofString(body));

        headers.forEach(builder::header);
        return builder.build();
    }

    private <T> T convertResult(String method, JsonNode result, Class<T> responseType) {
        if (JsonNode.class.equals(responseType)) {
            return responseType.cast(result);
        }

        try {
            return objectMapper.treeToValue(result, responseType);
        } catch (JsonProcessingException exception) {
            throw new FiberSerializationException(
                    "Failed to deserialize Fiber RPC result for method '%s'".formatted(method),
                    exception);
        }
    }

    private HttpResponse<String> send(String method, HttpRequest request) {
        try {
            return httpClient.send(request, HttpResponse.BodyHandlers.ofString());
        } catch (HttpTimeoutException exception) {
            throw new FiberTimeoutException(
                    "Fiber RPC call timed out for method '%s' after %s".formatted(method, requestTimeout),
                    exception);
        } catch (IOException exception) {
            throw new FiberTransportException(
                    "Transport failure calling Fiber RPC method '%s'".formatted(method),
                    exception);
        } catch (InterruptedException exception) {
            Thread.currentThread().interrupt();
            throw new FiberTransportException(
                    "Fiber RPC call interrupted for method '%s'".formatted(method),
                    exception);
        }
    }
}
