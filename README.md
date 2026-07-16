# FiberMan

FiberMan is a desktop-first Fiber operations and integration workbench for developers, wallet teams, merchants, and node operators who need to inspect live node behavior, execute RPC flows safely, and turn successful calls into reusable SDK code.

The current public web UI is a demonstration surface. The main product direction is the desktop app, where local-node workflows, safer operator tooling, and packaged multi-SDK developer experiences make more sense than a browser-only deployment.

Live demo:

- `https://fiberman.ebubeugwu.dev`

## Product Direction

FiberMan is being built around three product assumptions:

- serious Fiber integrations need access to real node behavior, not mocked examples
- teams want generated code and diagnostics close to the executed action
- desktop delivery is the best long-term form factor for local-node-first workflows, operator tooling, and secure developer environments

That is why the repo contains both a web demo path and a Wails desktop packaging path. The browser deployment proves the workflows. The desktop app is the intended core product.

## What It Includes

- `fiber-java-sdk`: plain Java SDK for Fiber JSON-RPC
- `fiber-go-sdk`: plain Go SDK for Fiber JSON-RPC
- `fiberman-java-backend`: Spring Boot backend built around the Java SDK
- `fiberman-go-backend`: Go backend built around the Go SDK
- `fiberman-frontend`: Angular UI for RPC exploration and workflow demos
- `fiberman-wails`: Wails desktop shell that reuses the Go runtime path

## Main Features

- live Fiber RPC exploration against a real node
- generated `cURL` snippets from executed requests
- generated Java SDK snippets from the same request path
- generated Go SDK snippets from the same request path
- node info, peers, channels, invoice, payment, and history workflows
- direct invoice QR generation for payment handoff flows
- runtime settings updates without rebuilding the app
- session-scoped history for replaying and reviewing prior actions
- all-in-one demo packaging that can bundle `fnn` with the app
- desktop packaging path for local-node-first usage

## Why The Web UI Exists

The web UI is primarily a demonstration and onboarding surface.

It is useful for:

- hackathon demos
- quick evaluator access
- showing the end-to-end product idea without requiring a local install
- validating frontend contract and workflow design quickly

It is not the final product thesis.

The main product thesis is:

- a desktop application for operators and developers
- powered by local or controlled Fiber runtime access
- with packaged SDK generation, diagnostics, and workflow tooling

## Current Architecture

Today the repo supports two main runtime paths:

1. Java path
   - `fiber-java-sdk`
   - `fiberman-java-backend`
   - Angular frontend
2. Go path
   - `fiber-go-sdk`
   - `fiberman-go-backend`
   - Angular frontend
   - Wails desktop shell

The Go path is the lighter long-term runtime direction, especially for desktop packaging and low-overhead deployments.

## Project Structure

- [fiber-java-sdk](fiber-java-sdk)
- [fiber-go-sdk](fiber-go-sdk)
- [fiberman-java-backend](fiberman-java-backend)
- [fiberman-go-backend](fiberman-go-backend)
- [fiberman-frontend](fiberman-frontend)
- [fiberman-wails](fiberman-wails)
- [deploy](deploy)
- [docs/fiber-hackathon-submission.md](docs/fiber-hackathon-submission.md)

## Local Run

Prerequisites:

- Java 21 for the Java backend path
- Go 1.26 for the Go backend path
- Node.js and npm for the frontend
- a reachable Fiber/FNN node

Java backend:

```bash
cd fiberman-java-backend
./gradlew bootRun
```

Go backend:

```bash
cd fiberman-go-backend
go run ./cmd/server
```

Frontend:

```bash
cd fiberman-frontend
npm start
```

Open `http://localhost:4200`.

For the current Go migration path, the frontend proxies `/api` to the Go backend on port `9020`.

## Runtime Settings

Use the `Settings` page to configure:

- node URL
- auth token
- request timeout
- default invoice currency
- playground base URL for generated snippets

These settings are applied by the running backend immediately for future RPC calls.

Recommendation:

- do not hardcode a currency like `FIBD` unless your node explicitly supports it
- leave the default invoice currency blank until you verify the asset code your node expects

## Deployment

FiberMan currently supports two practical deployment styles:

1. all-in-one container deployment for demos and judges
2. manual native deployment for lightweight servers and desktop-adjacent runtime control

### Deployed Demo

The current demo deployment is:

- `https://fiberman.ebubeugwu.dev`

This is a demo environment, not the final intended product form.

### Option 1: All-in-One Container

The root [Dockerfile](Dockerfile) and [Containerfile](Containerfile):

- build the Java SDK
- build the Angular frontend
- package the frontend into Spring static assets
- copy the official Fiber runtime into the final image
- start `fnn` and the Java backend together in one container
- serve both API and UI from one runtime

Build:

```bash
docker build -t fiberman .
```

Run:

```bash
docker run --rm \
  -p 9010:9010 \
  -p 8228:8228 \
  -v "$(pwd)/var/fiber-node:/fiber" \
  fiberman
```

Then open `http://localhost:9010`.

The repo also includes [compose.yaml](compose.yaml) for `docker compose` usage and a Podman-compatible [Containerfile](Containerfile).

### Option 2: Native Lightweight Deployment

For small servers, native deployment is often better than running a full container daemon.

The manual deployment assets live under [deploy/manual](deploy/manual):

- [deploy/manual/fnn-run.sh](deploy/manual/fnn-run.sh)
- [deploy/manual/fnn.service](deploy/manual/fnn.service)
- [deploy/manual/fiberman-go.service](deploy/manual/fiberman-go.service)
- [deploy/manual/fiberman.nginx.conf](deploy/manual/fiberman.nginx.conf)

This path is useful when you want:

- native `fnn` installation
- native `systemd` process management
- nginx in front of the static frontend and Go backend
- lower overhead than Docker on small machines

## Desktop Packaging

The desktop target lives in [fiberman-wails](fiberman-wails).

Build the Linux desktop binary:

```bash
cd fiberman-wails
sudo dnf install gcc-c++ pkgconf-pkg-config glib2-devel gtk3-devel webkit2gtk4.0-devel
wails build -platform linux/amd64
```

If `webkit2gtk4.0-devel` is not available on your Fedora release, use:

```bash
sudo dnf install gcc-c++ pkgconf-pkg-config glib2-devel gtk3-devel webkit2gtk4.1-devel
wails build -platform linux/amd64 -tags webkit2_41
```

The packaged binary will be written to:

```bash
fiberman-wails/build/bin/fiberman
```

The desktop wrapper:

- embeds the Angular frontend
- reuses the Go backend implementation
- starts a loopback API server internally
- preserves the same workflow model while moving the product closer to local-node-first usage

## SDK Expansion Plan

Current first-class SDKs:

- Java
- Go

Planned SDKs:

- TypeScript
- Python
- Rust
- C#

The plan for each is different because the target users are different.

### TypeScript SDK

Purpose:

- web apps
- Node.js backends
- Electron or hybrid wallet tooling

Expected shape:

- typed request and response models
- browser-safe transport where possible
- Node-oriented transport for backend usage
- parity with the playground snippet generator

### Python SDK

Purpose:

- automation
- operations scripting
- internal tools
- notebooks and rapid experimentation

Expected shape:

- simple client surface
- low-friction install
- examples for diagnostics and payment flows
- strong value for exchange, infra, and scripting teams

### Rust SDK

Purpose:

- systems integrations
- high-performance services
- wallet infrastructure
- security-sensitive and memory-conscious environments

Expected shape:

- typed models
- async-first transport
- strong error typing
- suitability for production backends and embedded infrastructure tooling

### C# SDK

Purpose:

- enterprise teams
- internal business tooling
- merchant infrastructure
- organizations already standardized on .NET

Expected shape:

- ergonomic .NET client API
- generated models and transport abstractions
- strong support for backend services and internal desktop tooling

### How We Plan To Build These SDKs

The SDK expansion strategy is not to hand-maintain unrelated clients forever.

The intended approach is:

- stabilize a canonical Fiber RPC schema and shared request model
- keep a thin shared transport and model contract where possible
- generate language-specific models and method surfaces where generation is reliable
- hand-curate ergonomics, errors, examples, and auth behavior per language
- keep the playground snippet generator aligned with real SDK implementations

In practice, FiberMan can become both:

- the product surface for interacting with Fiber
- the proving ground that validates new SDKs before they are published broadly

## Near-Term Roadmap

- complete Go backend parity with the Java path
- make the Wails desktop app the primary polished runtime
- improve node diagnostics and error inspection flows
- expand generated snippet coverage across more RPC methods
- improve history and replay workflows
- support richer channel and payment lifecycle visualization
- package a cleaner cross-platform desktop release flow

## Future Features If Funded

If FiberMan is funded and expanded, the product can move beyond a demo playground into a broader Fiber developer and operator platform.

Potential future features:

- first-class SDKs for TypeScript, Python, Rust, and C#
- one-click snippet switching across all supported languages
- packaged multi-network profiles for testnet, devnet, and custom operator networks
- richer node diagnostics, health scoring, and connectivity debugging
- channel lifecycle visualizations and topology mapping
- invoice and payment monitoring dashboards
- saved workspaces and persistent operator sessions
- team collaboration and shared environment profiles
- secure secret handling for desktop and managed deployments
- plugin or extension model for custom workflows
- desktop packaging for macOS and Windows in addition to Linux
- embedded documentation and guided onboarding for new Fiber teams
- test fixtures and mock-mode workflows for teams integrating before main node access
- CI-ready SDK smoke tooling and integration verification kits
- operator-grade logging export and structured audit trails

## Status

FiberMan is already credible as:

- a live Fiber demo
- a code-generation-assisted RPC explorer
- a Go migration path
- a desktop-first product direction

The next step is not to keep treating it as just a browser demo. The next step is to turn the desktop runtime and SDK expansion roadmap into the primary product.
