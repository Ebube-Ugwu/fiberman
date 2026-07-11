# Step 1 Implementation Plan: `fiber-java-sdk`

## Objective

Build the first working version of `fiber-java-sdk` as a plain Java library that can talk to a real Fiber node over JSON-RPC without any Spring dependency.

This step is complete when the SDK can successfully call the selected core Fiber methods against a real testnet node and surface both successful responses and failures cleanly.

## Scope

Included in Step 1:

- SDK module setup
- JSON-RPC transport layer
- request and response envelope handling
- core RPC method support
- configuration for node URL and auth headers
- error handling
- a runnable smoke test via `main()`

Not included in Step 1:

- Spring Boot backend
- REST controllers
- Angular frontend
- QR generation
- history
- deployment
- publishing the SDK to Maven Central

## Deliverable

Create a standalone Maven module named `fiber-java-sdk` with:

- a plain Java client API
- minimal DTOs for the selected methods
- transport and error handling
- a small executable example that validates real node connectivity

## Technical Decisions

### Language and packaging

- Use plain Java with Maven.
- Do not add Spring Boot or Spring Framework dependencies.
- Target a stable modern JDK version already used in your environment or intended deployment stack. Default to Java 21 unless the wider project requires lower compatibility.

### HTTP and JSON stack

- Use Java’s built-in `java.net.http.HttpClient` for transport.
- Use Jackson for JSON serialization and deserialization.
- Keep the transport dependency surface small: one HTTP client and one JSON library.

### API style

- Expose a top-level `FiberClient`.
- Keep method names aligned with Fiber RPC concepts.
- Prefer a typed request object per method where params are non-trivial.
- Return either typed response wrappers or `JsonNode`/structured DTOs depending on method complexity.
- For Step 1, favor thin DTOs over exhaustive domain modeling.

### Configuration

- `FiberClient` should accept:
  - base node URL
  - optional auth token or credentials if required
  - optional timeout configuration
- Keep configuration constructor-based or builder-based. Default to a builder if more than two optional fields are needed.

## Implementation Plan

### 1. Create the SDK module

- Add a new Maven module or standalone project directory for `fiber-java-sdk`.
- Create a clean package structure:
  - `client`
  - `model`
  - `transport`
  - `exception`
  - `example`
- Add dependencies for:
  - Jackson
  - test library if you want lightweight unit coverage in this step

### 2. Define JSON-RPC base models

- Create request envelope model with:
  - `jsonrpc`
  - `id`
  - `method`
  - `params`
- Create response envelope model with:
  - `jsonrpc`
  - `id`
  - `result`
  - `error`
- Create JSON-RPC error model with:
  - `code`
  - `message`
  - optional `data`
- Keep the envelope generic so the same transport path can be reused for every RPC method.

### 3. Implement transport

- Create a transport class responsible for:
  - serializing JSON-RPC requests
  - sending HTTP POST requests
  - attaching required headers
  - reading the response body
  - deserializing the JSON-RPC response
- Default headers:
  - `Content-Type: application/json`
  - `Accept: application/json`
- If the Fiber node needs authentication, support it through explicit configuration rather than hardcoding.

### 4. Implement SDK error handling

- Add SDK exceptions for:
  - transport failures
  - non-2xx HTTP responses
  - invalid JSON or schema mismatch
  - JSON-RPC error responses
  - timeout or interrupted request execution
- Ensure exception messages include enough context:
  - RPC method name
  - HTTP status when applicable
  - JSON-RPC error code and message when applicable
- Do not expose unreadable low-level exceptions directly without wrapping.

### 5. Build the public client API

- Implement `FiberClient` as the main entry point.
- Add a shared internal `call(method, params, responseType)` path used by all SDK methods.
- Implement these first public methods:
  - `nodeInfo()`
  - `createInvoice(CreateInvoiceRequest request)`
  - `sendPayment(SendPaymentRequest request)`
  - `listChannels()`
- Keep method names Java-idiomatic even if the underlying RPC name uses snake_case.
- Map each public method to the exact Fiber RPC method name internally.

### 6. Define minimal request/response DTOs

- Add request DTOs for methods that need params:
  - `CreateInvoiceRequest`
  - `SendPaymentRequest`
- Add minimal response DTOs only for fields needed to validate correctness in Step 1.
- Where the response schema is large or uncertain, it is acceptable in Step 1 to deserialize the result as `JsonNode` as long as:
  - method execution is stable
  - callers can inspect the returned data
  - the public API remains coherent

### 7. Add a runnable smoke test

- Create a small `main()` example under `example` or similar.
- Load node connection settings from environment variables or JVM properties.
- Execute the four core methods against a real Fiber testnet node.
- Print:
  - success status
  - method name
  - result summary or response JSON
- Fail fast on exceptions with readable output.

Recommended environment variables:

- `FIBER_NODE_URL`
- `FIBER_NODE_AUTH_TOKEN` if needed
- `FIBER_NODE_TIMEOUT_SECONDS` optional

### 8. Validate against a real node

- Confirm the exact payload shape expected by the node for each selected method.
- Verify:
  - node connectivity
  - authentication behavior
  - field naming conventions
  - response structure
- Adjust DTOs or raw result handling based on real responses, not assumptions.

## Suggested Public Surface

Example target API:

```java
FiberClient client = FiberClient.builder()
    .baseUrl(System.getenv("FIBER_NODE_URL"))
    .authToken(System.getenv("FIBER_NODE_AUTH_TOKEN"))
    .build();

var nodeInfo = client.nodeInfo();
var invoice = client.createInvoice(new CreateInvoiceRequest(...));
var payment = client.sendPayment(new SendPaymentRequest(...));
var channels = client.listChannels();
```

## Test and Validation Checklist

- `nodeInfo()` succeeds against the real node.
- `createInvoice(...)` succeeds with a valid test payload.
- `sendPayment(...)` executes successfully or returns a clean, inspectable RPC error if the test conditions do not permit payment.
- `listChannels()` returns a parsed result without transport failure.
- invalid node URL produces a transport exception with readable context.
- non-JSON or malformed JSON response produces a parse exception with method context.
- JSON-RPC error object is surfaced as an SDK-specific exception.
- missing auth or bad auth is surfaced clearly.

## Risks and Mitigations

- Fiber RPC docs may not match the actual node behavior.
  - Mitigation: validate payload and response structure against the real node early.
- Response schemas may be larger or less stable than expected.
  - Mitigation: keep Step 1 DTO modeling thin and allow `JsonNode` where necessary.
- `sendPayment` may be harder to validate than read-only calls.
  - Mitigation: treat clean RPC failure handling as acceptable proof of integration if payment prerequisites are not satisfied.

## Definition of Done

Step 1 is done when:

- `fiber-java-sdk` builds successfully
- the SDK has no Spring dependency
- the SDK exposes the 4 core methods through `FiberClient`
- the transport layer correctly sends JSON-RPC over HTTP POST
- failures are wrapped in usable SDK exceptions
- the `main()` smoke test runs against a real Fiber testnet node
- at least `nodeInfo`, `createInvoice`, and `listChannels` are confirmed working end to end, with `sendPayment` either working or failing through clean RPC-level handling
