# Fiber Hackathon Submission Pack

## Project Summary

FiberMan is a developer-facing Fiber infrastructure playground that makes node operations easier to test, understand, and integrate. It combines:

- a plain Java SDK for Fiber JSON-RPC access
- a Java backend that turns common Fiber operations into stable REST/demo endpoints
- an Angular frontend for exploring node info, peers, channels, invoices, payments, call history, and generated integration snippets

The project is intended to reduce the friction between "I want to try a Fiber operation" and "I now have working code and a verified request shape that I can reuse in my own wallet, merchant, or node tooling."

## Selected Category

Node, Routing, Cross-Chain, and Diagnostics Infrastructure

## Team Members

Fill in final team/member details before submission.

## Repository Link

Fill in the public repository URL before submission.

## Demo Link or Runnable Demo Instructions

Hosted demo link:

- Fill in deployed URL before submission.

Runnable demo instructions:

1. Start a reachable Fiber/FNN node.
2. Configure the node URL and auth token in the FiberMan settings page.
3. Run the backend:
   `cd fiberman-java-backend && ./gradlew bootRun`
4. Run the frontend:
   `cd fiberman-frontend && npm start`
5. Open `http://localhost:4200`.
6. Use `RPC Explorer`, `Invoice Builder`, and `Payments History` to exercise the integration flow.

## Video Demonstration

Fill in the final recorded demo link before submission.

Suggested demo flow:

1. Show node connectivity in the status pill.
2. Open `RPC Explorer`.
3. Run `node_info` and `list_channels`.
4. Generate an invoice and show the QR code.
5. Copy the generated `cURL` command.
6. Copy the generated Java SDK snippet.
7. Reopen the action from history and replay the code-copy flow.
8. Open `Settings` and show live runtime node reconfiguration.

## Technical Breakdown

### Architecture

- `fiber-java-sdk`
  - plain Java Fiber client using `java.net.http.HttpClient`
  - wraps transport, serialization, JSON-RPC envelope handling, and Fiber-specific exceptions
  - exposes reusable methods such as `nodeInfo`, `listChannels`, `listPeers`, `createInvoice`, `sendPayment`, `getChannel`, and `getPayment`
- `fiberman-java-backend`
  - uses the SDK to talk to a real Fiber node
  - exposes a thinner demo-facing API under `/api/fiber`
  - persists recent calls in session history
  - generates runnable `cURL` and Java SDK snippets from the exact executed action
  - exposes runtime settings so the playground can be repointed without restart
- `fiberman-frontend`
  - Angular-based explorer and invoice tooling
  - surfaces raw JSON responses, generated snippets, QR output, and recent history
  - adds a UI layer around live runtime settings for node URL, auth token, timeout, and invoice currency defaults

### Key Implementation Decisions

- Kept the SDK independent from Spring so it remains usable outside the demo app.
- Returned raw Fiber JSON where exhaustive domain modeling would slow delivery.
- Centralized code-generation in the backend so copied commands stay aligned with real executed requests.
- Added runtime settings to avoid baking in a single local-node setup.
- Preserved a lightweight demo-friendly history model instead of introducing a database.

### Current State

Fully working:

- live node connectivity
- RPC execution flow
- invoice generation flow
- QR generation
- copy-as-`cURL`
- copy-as-Java SDK
- session history
- live runtime settings updates

Intentionally deferred in this Java version:

- richer topology diagnostics
- deeper logs/observability workflows
- production auth, tenancy, and secret management

## Fiber Infrastructure Gap Addressed

Fiber currently needs better surrounding infrastructure for developers and operators who are trying to integrate or debug real node behavior. The gap is not only RPC access itself, but the lack of a fast feedback loop for:

- validating request shapes
- understanding failure responses
- inspecting peers/channels/node health
- turning successful experiments into reusable application code

FiberMan addresses that gap by acting as an educational and operational bridge between low-level Fiber RPC and reusable integration workflows. It is especially useful for:

- SDK consumers validating Fiber request/response behavior
- wallet or merchant teams exploring invoice and payment primitives
- node operators wanting a quick diagnostics surface for peers, channels, and recent actions
- developers who need working examples before embedding Fiber in a larger system

## Future Roadmap

- replace the Java backend with a Go SDK and Go-native app service
- ship a Wails desktop app for local-node-first usage
- add richer diagnostics for channel state, payment failures, and route confidence
- add exportable logs and structured operational events
- expand code generation beyond Java to Go and TypeScript
- deepen node configuration profiles and environment switching

## AI Allowance Claim

Fill in the final AI tooling claim details, if submitted.

## Submission Notes

- Be explicit during submission about which surfaces are production-ready versus demo-ready.
- Include the final hosted demo URL and video link before publishing.
- If any Fiber features are node-dependent, mention that clearly in the walkthrough.
