# FiberMan Wails Desktop

This package wraps the existing Angular frontend and Go backend into a Linux desktop app using Wails.

## Fedora Prerequisites

On Fedora, install the Linux desktop build dependencies first:

```bash
sudo dnf install gcc-c++ pkgconf-pkg-config glib2-devel gtk3-devel webkit2gtk4.0-devel
```

On Fedora 44 and newer systems where `webkit2gtk4.0-devel` is unavailable, use the 4.1 package and build with the Wails `webkit2_41` tag:

```bash
sudo dnf install gcc-c++ pkgconf-pkg-config glib2-devel gtk3-devel webkit2gtk4.1-devel
wails build -platform linux/amd64 -tags webkit2_41
```

## What It Does

- builds the Angular app from `../fiberman-frontend`
- embeds the built frontend into the desktop binary
- reuses the real Go backend from `../fiberman-go-backend`
- starts a local loopback API server so generated `cURL` snippets still point to a live endpoint

## Build For Linux

From this directory:

```bash
wails build -platform linux/amd64
```

The binary will be written to:

```bash
build/bin/fiberman
```

## Run In Desktop Mode

```bash
wails dev
```

Note:

- the Wails wrapper reuses the Angular production build for packaging
- for day-to-day UI work, continue using `../fiberman-frontend`
- runtime Fiber node configuration is still done inside the app from the `Settings` screen
- generated `cURL` snippets point at a loopback API server started by the desktop app at runtime
