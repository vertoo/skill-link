# skill-link

`skill-link` is a CLI tool built in Go for managing local agent skills (templates/scripts) via symlinks or copies.

## Features

- **CLI Interface**: Perform quick actions via CLI.
- **TUI Interface**: A rich terminal user interface to manage skills interactively.
- **Manifests**: Track installed skills through `.skill-link-lock.json`.

## Installation

Ensure you have [Go](https://go.dev/) installed (version 1.25.0 or later).

Clone the repository and build the binary:

```bash
# Build for your host OS
make build
```

### Cross-compilation

You can cross-compile the binary for other operating systems:

```bash
make build-linux
make build-darwin
make build-windows

# Build for all platforms
make cross-compile
```

The compiled binaries will be placed in the `build/` directory.

## Usage

```bash
skill-link <command> [arguments]
```

### Available Commands:

- `init`: Initialize a new local agent container. Creates a `.skill-link-lock.json` manifest.
- `create`: Create a new skill template globally.
- `manage`: Open the Terminal UI (TUI) to install, remove, and manage skills interactively.
- `restore`: Restore skills based on the local `.skill-link-lock.json` lockfile.
