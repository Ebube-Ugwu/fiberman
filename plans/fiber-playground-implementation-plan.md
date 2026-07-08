# Fiber Playground Implementation Plan

## Objective

Build a demo-ready Fiber Playground consisting of:

- a plain Java `fiber-java-sdk` for talking to a Fiber node over JSON-RPC
- a Spring Boot backend that exposes selected SDK functionality over REST
- an Angular frontend for exploring RPC methods, creating invoices, viewing responses, and copying runnable code snippets

The project should optimize for demo clarity, delivery speed, and standalone SDK value.

## Product Goals

- Prove that Java developers can integrate with Fiber without needing Spring or custom RPC plumbing.
- Provide a visual playground for judges and developers to inspect Fiber RPC behavior quickly.
- Turn every successful action into reusable code examples, especially `cURL` and, if time permits, Java SDK snippets.
- Keep the implementation thin and direct so it can be completed in one week.

## Core Demo Flow

1. Open the RPC Explorer.
2. Pick a method and fill in parameters.
3. Execute the call against a real Fiber testnet node.
4. Inspect the raw JSON response.
5. Copy the equivalent `cURL` command.
6. Copy the equivalent Java SDK snippet.
7. Build an invoice, render its QR code, and reuse it in follow-up actions.
8. Reopen recent calls from session history and copy code again retroactively.

## Architecture

### 1. `fiber-java-sdk`

Responsibilities:

- send JSON-RPC requests over HTTP
- model request and response envelopes
- handle RPC and transport errors cleanly
- expose typed or minimally structured Java methods for common Fiber operations
- stay independent of Spring Boot

Suggested module layout:

- `fiber-java-sdk`
  - `client/`
  - `model/`
  - `exception/`
  - `examples/` or simple test `main()`

### 2. `fiber-playground` backend

Responsibilities:

- depend on the SDK
- expose REST endpoints that delegate to SDK calls
- return raw JSON plus generated “copy as code” artifacts
- maintain lightweight session-based history
- optionally generate invoice QR images or payloads

### 3. Angular frontend

Responsibilities:

- render the RPC Explorer UI
- generate or submit dynamic parameter forms
- display raw JSON responses
- expose “Copy as cURL” and “Copy as Java” actions
- provide Invoice Builder, QR display, and recent call history

## Delivery Plan

## Day 1: SDK Foundation

Goals:

- create the `fiber-java-sdk` repository or module
- keep it plain Java with no Spring dependency
- get a local Fiber testnet node running
- implement the base JSON-RPC client
- verify real connectivity against the node

Tasks:

- set up `pom.xml` and base package structure
- implement HTTP POST transport for JSON-RPC
- define request and response envelope models
- add error handling for:
  - HTTP failures
  - invalid JSON
  - JSON-RPC error objects
  - timeouts and unreachable node scenarios
- implement the first core RPC methods:
  - `node_info`
  - `create_invoice`
  - `send_payment`
  - `list_channels`
- write a tiny standalone `main()` smoke test that calls each method against the real node

Exit criteria:

- each method successfully executes against the Fiber testnet node
- failures return useful errors instead of opaque stack traces

## Day 2: Finish SDK and Start Spring Boot App

Goals:

- complete the first public SDK surface
- package it cleanly
- start the backend that delegates directly to SDK calls

Tasks:

- add remaining exposed RPC methods such as:
  - payment status
  - channel details
  - peers
- refine method signatures and DTOs
- write SDK README with:
  - installation
  - configuration
  - supported methods
  - minimal usage example
- scaffold the `fiber-playground` Spring Boot app
- add the SDK as a dependency
- implement REST endpoints that:
  - map 1:1 to selected SDK calls
  - return raw JSON responses
  - preserve request params for later code generation

Exit criteria:

- backend can call the SDK successfully
- a REST client can hit the backend and get raw Fiber responses

## Day 3: Angular Frontend - RPC Explorer

Goals:

- build the main demo screen first
- make raw method execution stable before adding secondary features

Tasks:

- scaffold the Angular frontend
- build a method picker
- build dynamic parameter forms per method
- submit requests to the backend
- display raw JSON response bodies clearly
- handle loading, validation, and error states

Exit criteria:

- the RPC Explorer is fully functional end to end
- a user can select a method, provide params, execute it, and inspect the response

## Day 4: Copy-as-Code and Invoice/QR

Goals:

- add the feature with the highest demo leverage after the Explorer
- make successful calls reusable as real commands

Tasks:

- on every successful backend call, return enough metadata to generate:
  - equivalent `cURL`
  - optional Java SDK snippet
- decide whether code generation lives in backend or frontend
  - backend is preferred if you want consistency and easier history replay
  - frontend is acceptable if payload shapes are simple
- add a “Copy as cURL” button next to each response
- build the Invoice Builder screen with:
  - amount
  - description
  - generated invoice string
- add QR generation for invoices
  - backend option: ZXing
  - frontend option: JavaScript QR library

Recommended implementation detail:

- generate `cURL` from the exact request payload and endpoint used, including headers and actual parameter values
- ensure the copied command is runnable with minimal edits

Exit criteria:

- every successful RPC call exposes a working `Copy as cURL` action
- invoices can be generated and shown as QR codes

## Day 5: History and Copy-as-Java

Goals:

- make the playground feel persistent and educational
- turn the SDK itself into part of the demo experience

Tasks:

- add session-based history of calls including:
  - method name
  - timestamp
  - status
  - params
  - response summary or raw response
  - stored code-generation artifacts
- make history entries reopenable
- add “Copy as Java” using actual SDK method signatures
- keep the generated Java snippet aligned with the SDK README examples
- do basic UI cleanup and improve readability

Exit criteria:

- a user can revisit earlier calls and still copy `cURL`
- Java code examples look credible and directly useful as SDK documentation

## Day 6: Demo, Docs, Deploy

Goals:

- make the project externally consumable
- package both the software and the story for judging

Tasks:

- deploy the Spring Boot backend and Angular frontend
- point the deployment at a shared testnet node
- record demo flow:
  - RPC Explorer
  - execute call
  - copy as `cURL`
  - copy as Java
  - Invoice Builder
  - QR code
  - history replay
- finish the SDK README as a standalone deliverable
- write the submission technical breakdown
- write the “Fiber infrastructure gap addressed” section

Exit criteria:

- hosted demo works
- documentation is good enough for someone else to try the SDK without hand-holding

## Day 7: Buffer and Polish

Goals:

- absorb deployment issues
- harden presentation quality before submission

Tasks:

- fix deployment or environment issues
- polish generated code formatting and indentation
- verify copied commands are runnable
- finalize submission write-up
- add a future roadmap section

Suggested roadmap items:

- JavaFX desktop client
- payment debugger
- “Copy as TypeScript”
- Prometheus metrics

## Copy-as-Code Design

This feature has outsized demo value and should be treated as a core requirement, not polish.

### `cURL`

Requirements:

- generated from the exact method, endpoint, headers, and params used
- uses the real JSON-RPC payload shape
- remains available from both the immediate response view and history view

### Java

Requirements:

- generated from the actual SDK API, not hand-wavy pseudocode
- mirrors the method names and parameter structure exposed by `fiber-java-sdk`
- doubles as living SDK documentation

Suggested response shape from backend:

```json
{
  "method": "create_invoice",
  "params": {
    "amount": 1000,
    "description": "Demo invoice"
  },
  "result": {},
  "generatedCode": {
    "curl": "curl ...",
    "java": "FiberClient client = ..."
  },
  "timestamp": "2026-07-08T12:00:00Z",
  "status": "success"
}
```

## Priorities

### Must Have

- working Java SDK
- real Fiber node connectivity
- RPC Explorer
- raw JSON responses
- Invoice Builder
- `Copy as cURL`

### Should Have

- QR generation
- session history
- deployment
- strong README and submission docs

### Nice to Have

- `Copy as Java`
- richer response formatting
- extra RPC coverage beyond demo-critical methods

## Cut List

If time slips:

1. drop `Copy as Java` before dropping `Copy as cURL`
2. drop QR generation before dropping RPC Explorer or Invoice Builder
3. keep backend responses thin and avoid over-modeling every RPC method
4. prefer session-only history over persistent storage

## Technical Notes

- Keep the SDK transport and method layer cleanly separated so the SDK remains reusable outside the playground.
- Avoid premature abstraction in the frontend. The method picker and param renderer only need to support the selected demo methods well.
- Prefer returning raw JSON from the backend early. Typed backend response models can be added later only where they improve code generation or UI clarity.
- Ensure copied commands do not drift from the actual request execution path. If necessary, centralize request serialization in one backend service and reuse it for both execution and code generation.

## Definition of Done

The implementation is done when:

- the SDK successfully calls a real Fiber testnet node
- the Spring Boot backend exposes those calls reliably
- the Angular frontend can execute RPC methods and show raw results
- successful calls can be copied as runnable `cURL`
- invoice generation is demoable end to end
- the project is deployed and documented well enough for judges and developers to understand the value quickly
