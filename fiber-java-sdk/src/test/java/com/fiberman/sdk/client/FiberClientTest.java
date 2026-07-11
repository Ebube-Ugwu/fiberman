package com.fiberman.sdk.client;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fiberman.sdk.exception.FiberHttpException;
import com.fiberman.sdk.exception.FiberRpcException;
import com.fiberman.sdk.exception.FiberSerializationException;
import com.fiberman.sdk.exception.FiberTimeoutException;
import com.fiberman.sdk.model.CreateInvoiceRequest;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;
import java.io.IOException;
import java.io.InputStream;
import java.net.InetSocketAddress;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.util.concurrent.Executors;
import java.util.concurrent.atomic.AtomicReference;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertInstanceOf;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class FiberClientTest {
    private static final ObjectMapper OBJECT_MAPPER = new ObjectMapper().findAndRegisterModules();

    @Test
    void nodeInfoUsesJsonRpcAndAuthHeader() throws Exception {
        AtomicReference<String> requestBody = new AtomicReference<>();
        AtomicReference<String> authHeader = new AtomicReference<>();

        try (TestServer server = new TestServer(exchange -> {
            authHeader.set(exchange.getRequestHeaders().getFirst("Authorization"));
            requestBody.set(readBody(exchange));
            writeJson(exchange, 200, """
                    {"jsonrpc":"2.0","id":"1","result":{"version":"0.4.0","node_name":"local-node"}}
                    """);
        })) {
            FiberClient client = FiberClient.builder()
                    .baseUrl(server.baseUrl())
                    .authToken("secret-token")
                    .requestTimeout(Duration.ofSeconds(1))
                    .build();

            JsonNode result = client.nodeInfo();

            JsonNode requestJson = OBJECT_MAPPER.readTree(requestBody.get());
            assertEquals("2.0", requestJson.get("jsonrpc").asText());
            assertEquals("node_info", requestJson.get("method").asText());
            assertTrue(requestJson.get("id").asText().length() > 10);
            assertEquals("Bearer secret-token", authHeader.get());
            assertEquals("0.4.0", result.get("version").asText());
            assertEquals("local-node", result.get("node_name").asText());
        }
    }

    @Test
    void nodeInfoReturnsActualJsonPayloadWhenUsingJsonNode() throws Exception {
        try (TestServer server = new TestServer(exchange ->
                writeJson(exchange, 200, """
                        {"jsonrpc":"2.0","id":"1","result":{"version":"0.4.0","features":["rpc","payments"]}}
                        """))) {
            FiberClient client = FiberClient.builder()
                    .baseUrl(server.baseUrl())
                    .requestTimeout(Duration.ofSeconds(1))
                    .build();

            JsonNode result = client.nodeInfo();

            assertEquals("0.4.0", result.get("version").asText());
            assertEquals("rpc", result.get("features").get(0).asText());
            assertTrue(result.isObject());
        }
    }

    @Test
    void createInvoiceMapsToNewInvoiceAndSerializesCliFields() throws Exception {
        AtomicReference<String> requestBody = new AtomicReference<>();

        try (TestServer server = new TestServer(exchange -> {
            requestBody.set(readBody(exchange));
            writeJson(exchange, 200, """
                    {"jsonrpc":"2.0","id":"1","result":{"invoice_address":"fiber1testinvoice"}}
                    """);
        })) {
            FiberClient client = FiberClient.builder()
                    .baseUrl(server.baseUrl())
                    .requestTimeout(Duration.ofSeconds(1))
                    .build();

            JsonNode result = client.createInvoice(new CreateInvoiceRequest(
                    4200L,
                    "FIBD",
                    "sdk test invoice",
                    3600L));

            JsonNode requestJson = OBJECT_MAPPER.readTree(requestBody.get());
            JsonNode params = requestJson.get("params");
            assertEquals("new_invoice", requestJson.get("method").asText());
            assertEquals(4200L, params.get("amount").asLong());
            assertEquals("FIBD", params.get("currency").asText());
            assertEquals("sdk test invoice", params.get("description").asText());
            assertEquals(3600L, params.get("expiry").asLong());
            assertEquals("fiber1testinvoice", result.get("invoice_address").asText());
        }
    }

    @Test
    void getChannelMapsChannelIdCorrectly() throws Exception {
        AtomicReference<String> requestBody = new AtomicReference<>();

        try (TestServer server = new TestServer(exchange -> {
            requestBody.set(readBody(exchange));
            writeJson(exchange, 200, """
                    {"jsonrpc":"2.0","id":"1","result":{"channel_id":"abc123","state":"OPEN"}}
                    """);
        })) {
            FiberClient client = FiberClient.builder()
                    .baseUrl(server.baseUrl())
                    .requestTimeout(Duration.ofSeconds(1))
                    .build();

            JsonNode result = client.getChannel("abc123");

            JsonNode requestJson = OBJECT_MAPPER.readTree(requestBody.get());
            assertEquals("get_channel", requestJson.get("method").asText());
            assertEquals("abc123", requestJson.get("params").get("channel_id").asText());
            assertEquals("OPEN", result.get("state").asText());
        }
    }

    @Test
    void getPaymentMapsPaymentIdCorrectly() throws Exception {
        AtomicReference<String> requestBody = new AtomicReference<>();

        try (TestServer server = new TestServer(exchange -> {
            requestBody.set(readBody(exchange));
            writeJson(exchange, 200, """
                    {"jsonrpc":"2.0","id":"1","result":{"payment_id":"pay_123","status":"SETTLED"}}
                    """);
        })) {
            FiberClient client = FiberClient.builder()
                    .baseUrl(server.baseUrl())
                    .requestTimeout(Duration.ofSeconds(1))
                    .build();

            JsonNode result = client.getPayment("pay_123");

            JsonNode requestJson = OBJECT_MAPPER.readTree(requestBody.get());
            assertEquals("get_payment", requestJson.get("method").asText());
            assertEquals("pay_123", requestJson.get("params").get("payment_id").asText());
            assertEquals("SETTLED", result.get("status").asText());
        }
    }

    @Test
    void listPeersUsesExpectedMethodName() throws Exception {
        AtomicReference<String> requestBody = new AtomicReference<>();

        try (TestServer server = new TestServer(exchange -> {
            requestBody.set(readBody(exchange));
            writeJson(exchange, 200, """
                    {"jsonrpc":"2.0","id":"1","result":[{"peer_id":"peer-1"}]}
                    """);
        })) {
            FiberClient client = FiberClient.builder()
                    .baseUrl(server.baseUrl())
                    .requestTimeout(Duration.ofSeconds(1))
                    .build();

            JsonNode result = client.listPeers();

            JsonNode requestJson = OBJECT_MAPPER.readTree(requestBody.get());
            assertEquals("list_peers", requestJson.get("method").asText());
            assertEquals("peer-1", result.get(0).get("peer_id").asText());
        }
    }

    @Test
    void invokeAllowsAdHocMethods() throws Exception {
        AtomicReference<String> requestBody = new AtomicReference<>();

        try (TestServer server = new TestServer(exchange -> {
            requestBody.set(readBody(exchange));
            writeJson(exchange, 200, """
                    {"jsonrpc":"2.0","id":"1","result":{"ok":true}}
                    """);
        })) {
            FiberClient client = FiberClient.builder()
                    .baseUrl(server.baseUrl())
                    .requestTimeout(Duration.ofSeconds(1))
                    .build();

            JsonNode result = client.invoke("custom_method", java.util.Map.of("flag", true));

            JsonNode requestJson = OBJECT_MAPPER.readTree(requestBody.get());
            assertEquals("custom_method", requestJson.get("method").asText());
            assertTrue(requestJson.get("params").get("flag").asBoolean());
            assertTrue(result.get("ok").asBoolean());
        }
    }

    @Test
    void non2xxResponsesBecomeFiberHttpException() throws Exception {
        try (TestServer server = new TestServer(exchange ->
                writeJson(exchange, 401, """
                        {"error":"unauthorized"}
                        """))) {
            FiberClient client = FiberClient.builder()
                    .baseUrl(server.baseUrl())
                    .requestTimeout(Duration.ofSeconds(1))
                    .build();

            FiberHttpException exception = assertThrows(FiberHttpException.class, client::nodeInfo);

            assertEquals(401, exception.getStatusCode());
            assertTrue(exception.getMessage().contains("node_info"));
            assertTrue(exception.getResponseBody().contains("unauthorized"));
        }
    }

    @Test
    void rpcErrorsBecomeFiberRpcException() throws Exception {
        try (TestServer server = new TestServer(exchange ->
                writeJson(exchange, 200, """
                        {"jsonrpc":"2.0","id":"1","error":{"code":-32001,"message":"bad invoice","data":{"field":"invoice"}}}
                        """))) {
            FiberClient client = FiberClient.builder()
                    .baseUrl(server.baseUrl())
                    .requestTimeout(Duration.ofSeconds(1))
                    .build();

            FiberRpcException exception = assertThrows(FiberRpcException.class, client::nodeInfo);

            assertEquals(-32001, exception.getCode());
            assertTrue(exception.getMessage().contains("node_info"));
            assertEquals("invoice", exception.getData().get("field").asText());
        }
    }

    @Test
    void malformedJsonBecomesFiberSerializationException() throws Exception {
        try (TestServer server = new TestServer(exchange -> {
            byte[] body = "not-json".getBytes(StandardCharsets.UTF_8);
            exchange.sendResponseHeaders(200, body.length);
            exchange.getResponseBody().write(body);
            exchange.close();
        })) {
            FiberClient client = FiberClient.builder()
                    .baseUrl(server.baseUrl())
                    .requestTimeout(Duration.ofSeconds(1))
                    .build();

            FiberSerializationException exception = assertThrows(FiberSerializationException.class, client::nodeInfo);

            assertTrue(exception.getMessage().contains("node_info"));
        }
    }

    @Test
    void requestTimeoutBecomesFiberTimeoutException() throws Exception {
        try (TestServer server = new TestServer(exchange -> {
            try {
                Thread.sleep(200);
            } catch (InterruptedException exception) {
                Thread.currentThread().interrupt();
            }
            writeJson(exchange, 200, """
                    {"jsonrpc":"2.0","id":"1","result":{"status":"late"}}
                    """);
        })) {
            FiberClient client = FiberClient.builder()
                    .baseUrl(server.baseUrl())
                    .requestTimeout(Duration.ofMillis(50))
                    .connectTimeout(Duration.ofMillis(50))
                    .build();

            FiberTimeoutException exception = assertThrows(FiberTimeoutException.class, client::nodeInfo);

            assertTrue(exception.getMessage().contains("node_info"));
            assertTrue(exception.getMessage().contains("PT0.05S"));
        }
    }

    @Test
    void missingBaseUrlFailsFast() {
        IllegalStateException exception = assertThrows(IllegalStateException.class, () -> FiberClient.builder().build());
        assertEquals("baseUrl must be provided", exception.getMessage());
    }

    private static String readBody(HttpExchange exchange) throws IOException {
        try (InputStream inputStream = exchange.getRequestBody()) {
            return new String(inputStream.readAllBytes(), StandardCharsets.UTF_8);
        }
    }

    private static void writeJson(HttpExchange exchange, int statusCode, String body) throws IOException {
        byte[] responseBody = body.getBytes(StandardCharsets.UTF_8);
        exchange.getResponseHeaders().add("Content-Type", "application/json");
        exchange.sendResponseHeaders(statusCode, responseBody.length);
        exchange.getResponseBody().write(responseBody);
        exchange.close();
    }

    private static final class TestServer implements AutoCloseable {
        private final HttpServer server;

        private TestServer(ThrowingHandler handler) throws IOException {
            this.server = HttpServer.create(new InetSocketAddress(0), 0);
            this.server.createContext("/", exchange -> handler.handle(exchange));
            this.server.setExecutor(Executors.newCachedThreadPool());
            this.server.start();
        }

        private String baseUrl() {
            return "http://127.0.0.1:" + server.getAddress().getPort() + "/";
        }

        @Override
        public void close() {
            server.stop(0);
        }
    }

    @FunctionalInterface
    private interface ThrowingHandler {
        void handle(HttpExchange exchange) throws IOException;
    }
}
