# FiberMan

FiberMan is a desktop-first Fiber operations workbench for developers, wallet teams, merchants, and node operators. It helps you test live Fiber RPC flows, inspect node behavior, generate runnable code, and move faster from exploration to integration. It also feature reusable SDKs in Java and Golang (for now🙃), since the architecture is decoupled, the sdk can easily be reused to develop other application and/or libraries.

The web app is the public demo surface. The main product direction is the desktop app.

## Live Demo

Try it live:

- `https://fiberman.ebubeugwu.dev`

## What It Does

- explores live Fiber RPC methods against a real node
- generates runnable `cURL` commands from executed requests
- generates Java SDK snippets
- generates Go SDK snippets
- supports node info, peers, channels, invoices, payments, and request history
- generates invoice QR codes for payment flows
- allows runtime settings updates without rebuilding the app

## Product Direction

FiberMan is being built as a desktop-first product for local-node-first workflows, operator tooling, and secure developer environments.

Current surfaces:

- `fiberman-wails`: desktop shell built with Wails
- `fiberman-frontend`: Angular demo UI
- `fiberman-go-backend`: lighter Go runtime path
- `fiberman-java-backend`: Java runtime path

The web UI exists to demonstrate the workflows quickly. The desktop app is the primary long-term product direction.

## Screenshots

![FiberMan Screenshot 1](./assets/Screenshot%20From%202026-07-17%2000-01-44.png)
![FiberMan Screenshot 2](./assets/Screenshot%20From%202026-07-17%2000-06-42.png)
![FiberMan Screenshot 3](./assets/Screenshot%20From%202026-07-17%2000-07-36.png)
![FiberMan Screenshot 4](./assets/Screenshot%20From%202026-07-17%2000-07-44.png)
![FiberMan Screenshot 5](./assets/Screenshot%20From%202026-07-17%2000-08-00.png)

## Video Demo

<video src="./assets/Fiberman-demo.mp4" controls width="640" height="360">
  Your browser does not support the video tag.
</video>

[If it doesn't display use this](./assets/Fiberman-demo.mp4)

## Repository Layout

- [fiber-java-sdk](fiber-java-sdk): Java SDK for Fiber JSON-RPC
- [fiber-go-sdk](fiber-go-sdk): Go SDK for Fiber JSON-RPC
- [fiberman-java-backend](fiberman-java-backend): Spring Boot backend
- [fiberman-go-backend](fiberman-go-backend): Go backend
- [fiberman-frontend](fiberman-frontend): Angular UI
- [fiberman-wails](fiberman-wails): desktop packaging target
- [deploy](deploy): deployment assets

## Local Run

Prerequisites:

- Java 21 for the Java path
- Go 1.26 for the Go path
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

Then open `http://localhost:4200`.

## Deployment

FiberMan currently supports two practical deployment paths.

### 1. All-In-One Container

The root [Dockerfile](Dockerfile) and [Containerfile](Containerfile):

- build the frontend
- package the backend
- bundle the official Fiber runtime
- run `fnn` and the app together

Example:

```bash
docker build -t fiberman .
docker run --rm \
  -p 9010:9010 \
  -p 8228:8228 \
  -v "$(pwd)/var/fiber-node:/fiber" \
  fiberman
```

The repo also includes [compose.yaml](compose.yaml).

### 2. Lightweight Native Deploy

For smaller servers, the native path avoids the overhead of a full container daemon.

Manual deployment assets:

- [deploy/manual/fnn-run.sh](deploy/manual/fnn-run.sh)
- [deploy/manual/fnn.service](deploy/manual/fnn.service)
- [deploy/manual/fiberman-go.service](deploy/manual/fiberman-go.service)
- [deploy/manual/fiberman.nginx.conf](deploy/manual/fiberman.nginx.conf)

## Desktop Build

The desktop target lives in [fiberman-wails](fiberman-wails).

Linux build:

```bash
cd fiberman-wails
sudo dnf install gcc-c++ pkgconf-pkg-config glib2-devel gtk3-devel webkit2gtk4.0-devel
wails build -platform linux/amd64
```

Output:

```bash
fiberman-wails/build/bin/fiberman
```

## SDK Roadmap

Current SDKs:

- Java
- Go

Planned SDKs:

- TypeScript
- Python
- Rust
- C#

Planned approach:

- stabilize a canonical Fiber RPC contract
- share request and model definitions where possible
- generate model surfaces where generation is reliable
- hand-curate ergonomics, examples, and error handling per language

## Future Features

If FiberMan is funded and expanded, the next major features would include:

- first-class SDKs for TypeScript, Python, Rust, and C#
- snippet switching across all supported languages
- richer diagnostics and node health workflows
- channel and payment lifecycle visualization
- saved workspaces and persistent operator sessions
- improved team workflows and shared environment profiles
- broader desktop packaging beyond Linux
- embedded docs, onboarding, and integration guides

## Notes

- the web app is the demo surface
- the desktop app is the main product direction
- the Go runtime path is the lighter long-term execution path

## License

- [MIT](./LICENSE)
