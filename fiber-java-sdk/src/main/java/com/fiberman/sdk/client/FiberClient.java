package com.fiberman.sdk.client;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fiberman.sdk.model.CreateInvoiceRequest;
import com.fiberman.sdk.model.SendPaymentRequest;
import com.fiberman.sdk.transport.FiberTransport;
import java.net.URI;
import java.net.http.HttpClient;
import java.time.Duration;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.Objects;

public final class FiberClient {
    private final FiberTransport transport;

    private FiberClient(Builder builder) {
        ObjectMapper objectMapper = builder.objectMapper != null
                ? builder.objectMapper.copy()
                : new ObjectMapper().findAndRegisterModules();
        HttpClient httpClient = builder.httpClient != null
                ? builder.httpClient
                : HttpClient.newBuilder()
                        .connectTimeout(builder.connectTimeout)
                        .build();
        this.transport = new FiberTransport(
                httpClient,
                objectMapper,
                URI.create(builder.baseUrl),
                builder.requestTimeout,
                Map.copyOf(builder.headers));
    }

    public static Builder builder() {
        return new Builder();
    }

    public JsonNode nodeInfo() {
        return call("node_info", Map.of(), JsonNode.class);
    }

    public JsonNode createInvoice(CreateInvoiceRequest request) {
        Objects.requireNonNull(request, "request must not be null");
        return call("new_invoice", request, JsonNode.class);
    }

    public JsonNode sendPayment(SendPaymentRequest request) {
        Objects.requireNonNull(request, "request must not be null");
        return call("send_payment", request, JsonNode.class);
    }

    public JsonNode listChannels() {
        return call("list_channels", Map.of(), JsonNode.class);
    }

    public JsonNode listPeers() {
        return call("list_peers", Map.of(), JsonNode.class);
    }

    public JsonNode getChannel(String channelId) {
        Objects.requireNonNull(channelId, "channelId must not be null");
        return call("get_channel", Map.of("channel_id", channelId), JsonNode.class);
    }

    public JsonNode getPayment(String paymentId) {
        Objects.requireNonNull(paymentId, "paymentId must not be null");
        return call("get_payment", Map.of("payment_id", paymentId), JsonNode.class);
    }

    public JsonNode invoke(String method, Object params) {
        Objects.requireNonNull(method, "method must not be null");
        Objects.requireNonNull(params, "params must not be null");
        return call(method, params, JsonNode.class);
    }

    private <T> T call(String method, Object params, Class<T> responseType) {
        return transport.call(method, params, responseType);
    }

    public static final class Builder {
        private String baseUrl;
        private Duration connectTimeout = Duration.ofSeconds(10);
        private Duration requestTimeout = Duration.ofSeconds(30);
        private final Map<String, String> headers = new LinkedHashMap<>();
        private ObjectMapper objectMapper;
        private HttpClient httpClient;

        private Builder() {
        }

        public Builder baseUrl(String baseUrl) {
            this.baseUrl = Objects.requireNonNull(baseUrl, "baseUrl must not be null");
            return this;
        }

        public Builder authToken(String authToken) {
            Objects.requireNonNull(authToken, "authToken must not be null");
            this.headers.put("Authorization", "Bearer " + authToken);
            return this;
        }

        public Builder header(String name, String value) {
            this.headers.put(
                    Objects.requireNonNull(name, "header name must not be null"),
                    Objects.requireNonNull(value, "header value must not be null"));
            return this;
        }

        public Builder connectTimeout(Duration connectTimeout) {
            this.connectTimeout = Objects.requireNonNull(connectTimeout, "connectTimeout must not be null");
            return this;
        }

        public Builder requestTimeout(Duration requestTimeout) {
            this.requestTimeout = Objects.requireNonNull(requestTimeout, "requestTimeout must not be null");
            return this;
        }

        public Builder objectMapper(ObjectMapper objectMapper) {
            this.objectMapper = Objects.requireNonNull(objectMapper, "objectMapper must not be null");
            return this;
        }

        public Builder httpClient(HttpClient httpClient) {
            this.httpClient = Objects.requireNonNull(httpClient, "httpClient must not be null");
            return this;
        }

        public FiberClient build() {
            if (baseUrl == null || baseUrl.isBlank()) {
                throw new IllegalStateException("baseUrl must be provided");
            }
            return new FiberClient(this);
        }
    }
}
