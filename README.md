# FiberMan

FiberMan is a demo-ready Fiber infrastructure playground for developers, wallet teams, merchants, and node operators who need a faster way to explore Fiber RPC operations and convert successful calls into reusable code.

It includes:

- `fiber-java-sdk`: a plain Java SDK for Fiber JSON-RPC
- `fiberman-java-backend`: a Spring Boot backend that delegates to the SDK and generates code artifacts
- `fiber-go-sdk`: a plain Go SDK for Fiber JSON-RPC
- `fiberman-go-backend`: a Go backend that preserves the frontend contract and powers the desktop path
- `fiberman-frontend`: an Angular UI for RPC exploration, invoice building, QR generation, runtime settings, and session history
- `fiberman-wails`: a Linux desktop packaging target built on Wails

## Why It Exists

FiberMan targets a practical infrastructure gap: developers often need to validate real node behavior before they can safely embed Fiber into wallets, payment tooling, or operational dashboards. This repo reduces that friction by giving them:

- a visual RPC explorer
- runnable `cURL` generated from real calls
- Java SDK snippets generated from the same executed request
- a quick node-health, invoice, payment, and history workflow
- runtime node settings that can be updated without rebuilding the app

## Project Structure

- [fiber-java-sdk](fiber-java-sdk)
- [fiber-go-sdk](fiber-go-sdk)
- [fiberman-java-backend](fiberman-java-backend)
- [fiberman-go-backend](fiberman-go-backend)
- [fiberman-frontend](fiberman-frontend)
- [fiberman-wails](fiberman-wails)
- [docs/fiber-hackathon-submission.md](docs/fiber-hackathon-submission.md)

## Local Run

Prerequisites:

- Java 21
- Node.js and npm
- a reachable Fiber/FNN node

Backend:

```bash
cd fiberman-java-backend
./gradlew bootRun
```

Frontend:

```bash
cd fiberman-frontend
npm start
```

Open `http://localhost:4200`.

Go backend:

```bash
cd fiberman-go-backend
go run ./cmd/server
```

## Runtime Settings

Use the `Settings` page to configure:

- node URL
- auth token
- request timeout
- default invoice currency

These settings are applied by the running backend immediately for future RPC calls.

Recommendation:

- do not hardcode a currency like `FIBD` unless your node explicitly supports it
- leave the default invoice currency blank until you verify the asset code your node expects

## Deployment

A single-container deployment path is provided through the root [Dockerfile](Dockerfile) and [Containerfile](Containerfile). They:

- builds the Java SDK
- builds the Angular frontend
- packages the frontend into Spring static assets
- copies the official Fiber runtime into the final image
- starts FNN and the Java backend together in one container
- serves both the API and the web UI from that same container

### Option 1: Docker

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

On first boot the container:

- creates `/fiber/ckb/key` if one does not already exist
- copies the bundled Fiber testnet config into `/fiber/config.yml`
- starts FNN on loopback RPC `127.0.0.1:8227` for the Java backend
- exposes the Fiber peer port on `8228`

### Option 2: Full Judge Stack With `docker compose`

The repo includes [compose.yaml](compose.yaml).

Run:

```bash
docker compose up --build
```

Then open `http://localhost:9010`.

This stack now starts:

- one Fiberman container
- the Java backend and frontend inside that container
- the official Fiber runtime inside that same container

The bundled FNN runtime:

- persists node data under `var/fiber-node/`
- generates a disposable testnet private key automatically if `var/fiber-node/ckb/key` does not exist
- keeps RPC private on `127.0.0.1:8227` inside the container, which avoids the public RPC biscuit-key requirement
- uses SELinux-compatible bind mount labels so the same stack works with Podman on Fedora-class hosts

Override runtime settings by exporting env vars first, for example:

```bash
export FIBER_SECRET_KEY_PASSWORD=replace-this-for-shared-demo-use
export FIBER_PLAYGROUND_BASE_URL=http://your-server:9010
docker compose up --build
```

Notes:

- The default generated key is suitable for a disposable demo node only.
- To preserve identity or use a funded testnet account, replace `var/fiber-node/ckb/key` with your own testnet private key before starting the stack.
- The app talks to the bundled node at `http://127.0.0.1:8227` by default, so judges do not need a separate host-level Fiber installation.
- Channel creation and real payments still require testnet funds and liquidity on the node key you provide.

### Runtime Environment Variables

- `FIBER_NODE_URL`
- `FIBER_NODE_AUTH_TOKEN`
- `FIBER_NODE_TIMEOUT_SECONDS`
- `FIBER_PLAYGROUND_BASE_URL`
- `SERVER_PORT`
- `FIBER_SECRET_KEY_PASSWORD`
- `FIBER_RUST_LOG`

Notes:

- `FIBER_PLAYGROUND_BASE_URL` controls the base URL used in generated `cURL` snippets.
- If you deploy behind a public hostname, set `FIBER_PLAYGROUND_BASE_URL` to that public URL so generated snippets are correct.
- The container serves both the frontend and backend on the same port.

### Option 3: Podman

Build:

```bash
podman build -f Containerfile -t fiberman .
```

Run:

```bash
podman run --rm \
  -p 9010:9010 \
  -p 8228:8228 \
  -v "$(pwd)/var/fiber-node:/fiber:Z" \
  fiberman
```

Podman notes:

- `Containerfile` defaults `FIBER_NODE_URL` to `http://127.0.0.1:8227` because FNN runs in the same container.
- The `:Z` bind mount label is important on SELinux-enabled hosts.
- You can also run the same all-in-one image through [compose.yaml](compose.yaml) if you prefer a compose workflow.

## Linux Desktop Packaging

The Wails packaging target lives in [fiberman-wails](fiberman-wails).

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
- starts a loopback API server internally so generated `cURL` snippets still target a live local endpoint

## Submission Material

The hackathon submission pack lives in [docs/fiber-hackathon-submission.md](docs/fiber-hackathon-submission.md). Fill in the final repository URL, hosted demo link, team info, and video link before submission.

## Future Development Plans

We plan to expand FiberMan beyond the current Java and Go SDK paths.

Near-term work:

- finish the Go backend parity pass and move desktop packaging to Wails
- complete QR generation and remaining diagnostics gaps in the Go backend
- preserve the current Angular contract while migrating the runtime fully to Go

SDK expansion roadmap:

- keep first-class support for Java and Go
- add TypeScript SDK support for browser and Node-integrated application flows
- add Rust SDK support for systems integrations and performance-sensitive backend services
- extend the playground so users can switch snippet languages directly from the UI as new SDKs become available

The current explorer already reserves space for future language targets so the SDK surface can grow without redesigning the product.
