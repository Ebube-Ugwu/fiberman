# FiberMan Progress Report

Date: 2026-07-13

## Executive Summary

FiberMan has moved from an initial one-path MVP into a multi-surface Fiber infrastructure playground. The original build focused on a plain Java SDK, a Spring Boot backend, and an Angular frontend. That stack is now implemented and still present. On top of that, the project has expanded into a Go SDK, a Go backend intended to preserve the same frontend contract, a Wails desktop packaging path, and a containerized all-in-one deployment path that bundles the official Fiber runtime.

The repo reflects two distinct stages of progress:

1. A completed Java-centered MVP for hackathon demo flow.
2. A second-stage runtime and packaging expansion toward Go parity, desktop delivery, and judge-friendly deployment.

The core product direction has remained stable throughout: make Fiber easier to explore, test, verify, and integrate by shortening the loop between a live RPC call and reusable implementation code.

## Original Objective

The original implementation plan in [plans/fiberman-implementation-plan.md](/home/ebube/projects/portfolio/fiberman/plans/fiberman-implementation-plan.md:1) set out to build:

- a plain Java `fiber-java-sdk`
- a Spring Boot backend that wraps selected Fiber JSON-RPC methods
- an Angular frontend for RPC exploration, invoice generation, QR rendering, history, and generated code snippets

The main demo flow was:

- execute real Fiber RPC calls
- inspect raw JSON responses
- copy equivalent `curl`
- copy equivalent Java SDK code
- generate an invoice and QR code
- replay or inspect recent history

That initial scope is largely implemented.

## Product Direction We Chose

The project consistently chose the infrastructure-playground direction instead of building an end-user wallet or merchant app. This is aligned with the hackathon category around node, routing, and diagnostics infrastructure.

The practical problem being addressed is not just RPC access. It is the lack of a fast workflow for:

- validating request shapes against a real node
- understanding response and failure payloads
- turning successful experiments into code a developer can actually reuse
- testing invoice and payment primitives without writing all the plumbing first

## Major Technical Decisions

### 1. Keep the SDKs framework-independent

This was one of the strongest and most consistent decisions in the repo.

- `fiber-java-sdk` is plain Java and uses Maven, not Spring.
- `fiber-go-sdk` is plain Go and does not depend on a web framework.
- The backend applications depend on the SDKs, not the other way around.

Why this mattered:

- the SDKs remain reusable outside the demo app
- the hackathon output becomes more than just one monolithic demo
- the backend stays thin and demo-oriented instead of becoming the integration surface itself

### 2. Use real JSON-RPC over HTTP instead of mocking Fiber behavior

Both SDK lines are built around direct JSON-RPC transport:

- Java uses `java.net.http.HttpClient`
- Go uses `net/http`

Why this mattered:

- the project demonstrates real integration behavior, not simulated behavior
- request/response shapes can be validated against actual node expectations
- copied snippets are grounded in executed calls

### 3. Favor thin modeling and raw JSON over exhaustive domain modeling

The SDKs and backends intentionally return flexible JSON-compatible structures in many cases instead of trying to fully model every Fiber response upfront.

Why this mattered:

- Fiber response shapes can evolve or be larger than expected
- hackathon delivery speed mattered more than complete schema coverage
- the playground benefits from showing raw payloads directly

Tradeoff:

- weaker compile-time guarantees for result structures
- more frontend extraction logic for arrays and nested fields

### 4. Centralize generated code artifacts in the backend

The backend, not the frontend, generates:

- `curl`
- Java SDK snippets
- and now Go SDK snippets

This is implemented in the Java backend through `CodeArtifactService` and mirrored in the Go path.

Why this mattered:

- copied code stays aligned with the exact request the backend executed
- history replay can reuse the same generated artifacts
- the frontend stays focused on display and interaction rather than code synthesis rules

### 5. Preserve lightweight session-scoped history instead of adding a database

The Java backend stores history in `HttpSession` via `FiberHistoryService`. The Go backend mirrors the same idea with an in-memory session store.

Why this mattered:

- zero database setup for judges and demo use
- enough persistence for a single-user or single-session playground
- lower implementation cost

Tradeoff:

- history is intentionally ephemeral
- not suitable for multi-user production analytics or audit requirements

### 6. Add runtime-editable node settings instead of hardcoding one environment

This is a key product and technical decision that appears after the initial MVP plan and meaningfully improved usability.

The app allows runtime updates for:

- node URL
- auth token
- request timeout
- default invoice currency

In Java this is handled by `RuntimeFiberSettingsService`, which creates a fresh `FiberClient` per request using the current settings.

Why this mattered:

- the app can be repointed without restart
- demo operators are not locked to one local-node assumption
- generated artifacts can reflect the active runtime configuration

### 7. Keep the frontend contract stable while exploring a Go runtime migration

This is one of the clearest later-stage architecture decisions in the repo.

The Go backend explicitly aims to preserve the existing Angular contract:

- same `/api/fiber/*` routes
- same `/api/settings`
- same JSON response structure

Why this mattered:

- migration risk stays low
- the Angular UI can continue to work during backend substitution
- the project can evolve toward Go and Wails without rewriting the whole frontend

### 8. Package a judge-friendly all-in-one deployment

The root `Dockerfile`, `Containerfile`, `compose.yaml`, and `deploy/fiberman-container-entrypoint.sh` show a deliberate move toward a self-contained deployment experience.

Key decisions:

- build the Java SDK and backend in the image
- build the Angular frontend into Spring static assets
- include the official `fnn` runtime in the final image
- start both FNN and the app in one container
- keep FNN RPC on loopback inside the container
- persist node data under `var/fiber-node`

Why this mattered:

- judges do not need a separate Fiber installation
- the demo can run as a single artifact
- loopback-only RPC avoids exposing a public node RPC surface in the default setup

### 9. Explore desktop packaging through Wails rather than Electron

The repo includes `fiberman-wails`, which wraps:

- the Angular frontend
- the Go backend
- a loopback API server for the desktop runtime

Why this mattered:

- lower overhead than a browser-only plus external-backend workflow
- better local-node-first operator experience
- a desktop path is credible for node tooling and diagnostics use cases

## What We Planned vs What We Actually Built

### Planned in the original Java MVP

- Java SDK
- Spring Boot backend
- Angular frontend
- RPC explorer
- invoice builder
- QR generation
- history
- `curl` generation
- Java snippet generation

### Actually built in the Java path

- `fiber-java-sdk`
- Spring Boot backend under `fiberman-java-backend`
- Angular app under `fiberman-frontend`
- RPC explorer for selected methods
- session history
- invoice creation
- QR generation
- `curl` generation
- Java snippet generation
- runtime settings updates

### Built beyond the original plan

- `fiber-go-sdk`
- `fiberman-go-backend`
- Go snippet generation
- Wails desktop shell
- single-container Fiber + app runtime
- `docker compose` and Podman-oriented deployment path

## What We Tried So Far

This section reflects the implementation history visible in the repo structure, docs, commit history, and current working tree.

### Phase 1: Java-first MVP

The first meaningful shipped phase was:

- `fiber-java-sdk`
- Spring Boot backend
- Angular frontend

The git history shows this in commit `1719ed8` with the message `Add Fiber Java SDK and MVP backend`.

What this phase proved:

- Java can be a viable integration surface for Fiber JSON-RPC
- a thin backend can convert SDK calls into a demo-friendly REST surface
- the explorer/invoice/history/code-copy demo loop is feasible

### Phase 2: Runtime flexibility and better demo ergonomics

The Java backend evolved beyond a static configuration approach.

Visible changes include:

- runtime settings support
- backend creation of clients from mutable settings
- environment-driven `application.properties`
- frontend settings UI

This addressed a practical issue: a playground tied to one baked-in node is fragile during demos and deployment.

### Phase 3: Multi-language code generation

The repo now includes generated Go snippets in addition to Java snippets and `curl`.

This is visible in:

- `CodeArtifacts` gaining `goSnippet`
- `CodeArtifactService` generating Go examples
- the frontend snippet selector exposing Java and Golang while reserving TypeScript and Rust as planned-but-unimplemented targets

This indicates a deliberate move from "demo a Java SDK" to "demo Fiber integration patterns across languages."

### Phase 4: Go migration experiment and desktop path

The repo then expands into:

- `fiber-go-sdk`
- `fiberman-go-backend`
- `fiberman-wails`

This is not just parallel experimentation. The Go backend README states that it is meant to preserve the frontend contract while replacing the Java backend, and the Wails package reuses the Go backend as the desktop runtime.

This shows a concrete architectural exploration:

- keep the existing UI
- move runtime logic toward Go
- unlock desktop packaging around the Go server

### Phase 5: All-in-one deployment packaging

The root container assets show another explicit attempt:

- make the project runnable without asking judges to wire up Fiber separately
- bundle FNN and the app in one image
- support Docker, Podman, and `docker compose`

This is a substantial demo-readiness improvement, not just ops polish.

## Implemented Components

## 1. `fiber-java-sdk`

Location:

- [fiber-java-sdk/pom.xml](/home/ebube/projects/portfolio/fiberman/fiber-java-sdk/pom.xml:1)
- Java sources under `fiber-java-sdk/src/main/java`

Implemented characteristics:

- plain Java library
- Java 21 target
- Jackson-based JSON serialization/deserialization
- built-in Java HTTP client
- custom exception hierarchy for HTTP, RPC, timeout, transport, and serialization failures
- request DTOs for invoice and payment flows
- support for:
  - `nodeInfo`
  - `createInvoice`
  - `sendPayment`
  - `listChannels`
  - `listPeers`
  - `getChannel`
  - `getPayment`
  - ad hoc `invoke`

Validation status:

- `mvn test` passes in `fiber-java-sdk`
- current result: 12 tests, 0 failures, 0 errors

What this proves:

- the Java SDK surface is the most concretely validated part of the repo

## 2. `fiberman-java-backend`

Location:

- [fiberman-java-backend/build.gradle](/home/ebube/projects/portfolio/fiberman/fiberman-java-backend/build.gradle:1)
- controller and service sources under `fiberman-java-backend/src/main/java`

Implemented API endpoints:

- `GET /api/fiber/node-info`
- `GET /api/fiber/channels`
- `POST /api/fiber/channels/details`
- `GET /api/fiber/peers`
- `POST /api/fiber/invoices`
- `POST /api/fiber/invoices/qr`
- `POST /api/fiber/payments`
- `POST /api/fiber/payments/status`
- `GET /api/fiber/history`
- `DELETE /api/fiber/history`
- `GET /api/settings`
- `PUT /api/settings`

Implemented backend responsibilities:

- delegate to `fiber-java-sdk`
- normalize success and error payloads into `FiberCallResponse`
- record session-scoped history
- generate `curl`, Java, and Go code artifacts
- generate QR payloads for invoice strings
- expose runtime-editable settings

Notable implementation detail:

- the backend now creates a fresh `FiberClient` from runtime settings for each call instead of relying on one fixed configured bean

Why that is important:

- runtime settings are real, not cosmetic
- snippet generation can stay aligned with active runtime config

## 3. `fiberman-frontend`

Location:

- [fiberman-frontend/src/app/app.routes.ts](/home/ebube/projects/portfolio/fiberman/fiberman-frontend/src/app/app.routes.ts:1)

Implemented pages:

- Dashboard
- RPC Explorer
- Invoice Builder
- Network Topology
- Payments History
- Logs
- Settings

Implemented frontend capabilities:

- method selection and request execution for selected RPC methods
- display of raw response payloads
- copy generated snippets
- history browsing
- runtime settings editing
- invoice creation flow

Important current-state observation:

- the frontend is a mix of implemented product flow and presentation scaffolding

Examples:

- `RPC Explorer`, `Settings`, `Invoice Builder`, and history interactions are backed by real API calls.
- `Dashboard` includes derived values but also fallback and illustrative values such as default liquidity text and hardcoded latency data.
- `Network Topology` computes node placement visually but falls back to generated placeholder peers when real peer data is absent.
- `Logs` is derived from request history rather than a true backend log stream.

Conclusion:

- the frontend is demo-ready for core workflow
- some screens are still presentational or derived rather than operationally deep

## 4. `fiber-go-sdk`

Location:

- [fiber-go-sdk/README.md](/home/ebube/projects/portfolio/fiberman/fiber-go-sdk/README.md:1)
- [fiber-go-sdk/client/client.go](/home/ebube/projects/portfolio/fiberman/fiber-go-sdk/client/client.go:1)

Implemented characteristics:

- plain Go SDK
- config-based client construction
- auth header support
- timeout configuration
- JSON-RPC transport abstraction
- support for:
  - `NodeInfo`
  - `CreateInvoice`
  - `SendPayment`
  - `ListChannels`
  - `ListPeers`
  - `GetChannel`
  - `GetPayment`
  - `Invoke`

This appears to be the Go counterpart to the Java SDK rather than a partial mock or placeholder.

## 5. `fiberman-go-backend`

Location:

- [fiberman-go-backend/README.md](/home/ebube/projects/portfolio/fiberman/fiberman-go-backend/README.md:1)
- [fiberman-go-backend/server.go](/home/ebube/projects/portfolio/fiberman/fiberman-go-backend/server.go:1)

Implemented characteristics:

- REST surface intentionally matching the frontend contract
- runtime settings
- session-based history
- QR generation
- generated `curl`, Java, and Go snippets
- validation logic for request payloads

Interpretation:

- this is a serious migration path, not a speculative stub

## 6. `fiberman-wails`

Location:

- [fiberman-wails/README.md](/home/ebube/projects/portfolio/fiberman/fiberman-wails/README.md:1)
- [fiberman-wails/desktop_runtime.go](/home/ebube/projects/portfolio/fiberman/fiberman-wails/desktop_runtime.go:1)

Implemented characteristics:

- desktop wrapper around the Angular build
- embedded loopback backend runtime
- backend port chosen dynamically on `127.0.0.1`
- desktop app lifecycle wiring via Wails

Interpretation:

- the project has moved beyond a browser-only demo and has a credible desktop-node-tooling direction

## 7. Containerized deployment

Location:

- [Dockerfile](/home/ebube/projects/portfolio/fiberman/Dockerfile:1)
- [compose.yaml](/home/ebube/projects/portfolio/fiberman/compose.yaml:1)
- [deploy/fiberman-container-entrypoint.sh](/home/ebube/projects/portfolio/fiberman/deploy/fiberman-container-entrypoint.sh:1)

Implemented characteristics:

- multi-stage build
- Java SDK built before backend packaging
- Angular build copied into Spring static assets
- official Fiber runtime bundled
- disposable testnet key auto-generation
- default testnet config copy-on-first-run
- loopback RPC hardening inside container
- coordinated shutdown between app and FNN processes

This is one of the strongest signs of project maturity from a demo-delivery standpoint.

## Current Worktree Changes

The current uncommitted worktree shows the project is still actively evolving.

Tracked modifications currently visible:

- `README.md`
- Java backend `CodeArtifacts`
- Java backend `CodeArtifactService`
- Java backend `FiberGatewayService`
- Java backend `application.properties`

New untracked additions currently visible include:

- root deployment files
- `docs/`
- `fiber-go-sdk/`
- `fiberman-go-backend/`
- `fiberman-frontend/`
- `fiberman-wails/`
- Java runtime settings classes and controller additions
- UI assets and supporting directories

Interpretation:

- the repository contains substantially more implemented work than what has been committed to `main` so far
- a large portion of the current project state is still in-progress or awaiting a consolidation commit

## Validation and Testing Status

### Confirmed

- `fiber-java-sdk` tests pass locally with Maven
- result observed: 12 passing tests

### Blocked during this review

- `fiberman-java-backend` tests could not be fully executed in this environment because:
  - the default Gradle wrapper path attempted to write under `~/.gradle`, which is read-only here
  - rerunning with a workspace-local `GRADLE_USER_HOME` then required downloading the Gradle distribution
  - the environment blocks outbound network access, so the wrapper download failed

### Additional source-level concern

There is at least one likely backend test breakage visible from source inspection:

- [fiberman-java-backend/src/test/java/com/fiberman/fiberman_java_backend/service/FiberHistoryServiceTest.java](/home/ebube/projects/portfolio/fiberman/fiberman-java-backend/src/test/java/com/fiberman/fiberman_java_backend/service/FiberHistoryServiceTest.java:1) still constructs `new CodeArtifacts("curl", "java")`
- `CodeArtifacts` now requires three fields: `curl`, `javaSnippet`, and `goSnippet`

Unless there is local unpublished code compensating for this, that test will not compile in the current tree.

What this means:

- the Java backend implementation has advanced
- backend test fixtures have not fully kept pace with the latest code-artifact shape changes

## What Is Working Best Right Now

The strongest completed pieces are:

- the Java SDK design and tests
- the Java backend’s core API shape and runtime settings approach
- the Angular explorer/history/settings workflow
- the multi-language snippet generation direction
- the all-in-one deployment packaging

These pieces form a coherent and credible hackathon submission story.

## What Is Still In Progress or Not Fully Hardened

### 1. Backend test hygiene after recent changes

The Go snippet addition appears to have changed shared response DTOs without fully updating all Java backend tests.

### 2. Frontend depth on diagnostics-oriented screens

Some pages are real operational surfaces, but others are still partially illustrative:

- dashboard metrics include fallback values and mock-style charts
- logs are derived from API history rather than backend telemetry
- topology is visual and useful for demoing peer layout, but not yet a deep diagnostics console

### 3. Full Go parity and migration completion

The Go backend and SDK are substantial, but the repo text still describes them as part of an ongoing replacement/parity effort rather than the fully finalized runtime default.

### 4. Production hardening

The code and docs themselves suggest that the current scope is demo-ready rather than production-ready in areas such as:

- secrets management
- multi-user persistence
- operational logging and observability
- tenancy and auth controls

## High-Confidence Narrative of Progress

The project did not stall after the first MVP. It expanded in a technically coherent way:

1. Build a Java SDK and Java-backed playground to prove the integration pattern.
2. Add runtime mutability so the app can be repointed during demos.
3. Add better generated artifacts so the playground teaches integration, not just execution.
4. Start a Go-native runtime path while preserving the frontend contract.
5. Package the project for desktop and all-in-one deployment.

That is a rational sequence. It improves demoability, reusability, and long-term maintainability without throwing away the earlier work.

## Recommended Positioning in External Communication

If this report is used to brief judges, reviewers, or collaborators, the most accurate framing is:

- FiberMan is already a functioning developer playground for live Fiber RPC exploration.
- The Java path is the most validated implementation today.
- The Go path, desktop wrapper, and containerized runtime represent active expansion toward broader portability and easier delivery.
- Some frontend surfaces are polished demo views rather than deep operational tooling.
- The project is strong on integration workflow and packaging, with remaining work mostly in hardening and parity cleanup.

## Suggested Immediate Next Steps

Based on the current state of the repo, the highest-value next steps are:

1. Fix Java backend test fixtures for the updated `CodeArtifacts` shape and rerun backend tests in a network-enabled environment.
2. Decide whether the primary judged runtime is the Java backend or the Go backend and document that explicitly.
3. Separate demo-only visualizations from fully live diagnostics in the README and submission pack.
4. Finish the Go parity pass only if it materially improves the judged submission; otherwise prioritize stability on the Java path.
5. Record one authoritative deployment path for judges, most likely the all-in-one container flow.

## Bottom Line

FiberMan is not just an idea or a wireframe set. It is an implemented Fiber infrastructure playground with:

- a working Java SDK
- a working Java backend pattern
- a functioning Angular demo workflow
- runtime configuration support
- generated `curl`, Java, and Go artifacts
- invoice QR support
- a Go migration path
- a Wails desktop path
- an all-in-one container deployment path

The main remaining gaps are not about whether the concept has been built. They are about stabilization, test alignment, and deciding which runtime path should be treated as the canonical submission surface.
