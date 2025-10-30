# osdd-cli

OpenSDD Command Line Interface - A Go-based CLI for accessing OpenSDD flows with a beautiful branded terminal experience.

## Features

- **Cobra-based CLI** with intuitive command structure
- **Rainbow-colored ASCII art** version display
- **Graceful fallback** to plain ASCII art in non-TTY or limited terminal environments
- **Build-time version injection** via Go ldflags
- **Cross-platform support** for Linux, macOS, and Windows
- **Automated releases** via GitHub Actions for all major platforms

## Installation

### From GitHub Releases

Download the latest release for your platform from the [Releases page](https://github.com/opensdd/osdd-cli/releases):

- **Linux**: `osdd-linux-x64`
- **macOS**: `osdd-macos-x64`
- **Windows**: `osdd-windows-x64.exe`

#### Linux/macOS

```bash
# Download the binary (replace VERSION with the latest release)
curl -L -o osdd https://github.com/opensdd/osdd-cli/releases/download/VERSION/osdd-linux-x64  # or osdd-macos-x64

# Make it executable
chmod +x osdd

# Move to a directory in your PATH (optional)
sudo mv osdd /usr/local/bin/

# Verify installation
osdd version
```

#### Windows

```powershell
# Download the binary from Releases page
# Or use curl in PowerShell 7+
curl -L -o osdd.exe https://github.com/opensdd/osdd-cli/releases/download/VERSION/osdd-windows-x64.exe

# Run the CLI
.\osdd.exe version
```

### From Source

Requires Go 1.25.1 or later.

```bash
# Clone the repository
git clone https://github.com/opensdd/osdd-cli.git
cd osdd-cli

# Build for development (version will show as "dev")
make build-dev

# Or build with a specific version
make build VERSION=v1.0.0

# Run the CLI
./osdd version
```

## Usage

### Available Commands

```bash
# Display help
osdd
osdd --help

# Display version with animated ASCII art
osdd version

# Get help for a specific command
osdd [command] --help
```

### Version Command

The `osdd version` command displays the OpenSDD ASCII art logo with rainbow colors (on supported terminals) followed by the version information.

**Colored output** (on terminals with 256+ color support and TTY):
- Rainbow gradient from red → orange → yellow → green → cyan → blue → purple
- Static display (no animation)

**Plain output** (on limited terminals, CI/CD, or when piped):
- Simple block ASCII art
- Version text
- Warning if version is "dev" (not set at build time)

The display automatically detects terminal capabilities and falls back gracefully to ensure consistent behavior across all environments.

## Building Locally

### Prerequisites

- Go 1.25.1 or later
- Make (optional, but recommended)

### Build Commands

#### Release Build

Builds the CLI with a specific version injected via ldflags:

```bash
make build VERSION=v1.0.0
./osdd version
# Output: OpenSDD CLI version v1.0.0
```

#### Development Build

Builds the CLI with version set to "dev":

```bash
make build-dev
./osdd version
# Output: OpenSDD CLI version dev
# WARNING: Version not set at build time
```

#### Manual Build

Build without version injection:

```bash
go build -o osdd ./cmd/osdd
./osdd version
# Output: OpenSDD CLI version dev
# WARNING: Version not set at build time
```

#### Custom Build with ldflags

```bash
go build -ldflags "-X 'github.com/opensdd/osdd-cli/internal/version.Version=v2.1.0'" -o osdd ./cmd/osdd
./osdd version
# Output: OpenSDD CLI version v2.1.0
```

### Build Script

The `build.sh` script provides a portable way to build the CLI with version injection:

```bash
# Requires VERSION environment variable
VERSION=v1.0.0 ./build.sh

# Version format validation (must be vX.Y.Z or "dev")
VERSION=invalid ./build.sh
# Error: VERSION must be 'dev' or in format vX.Y.Z (e.g., v1.0.0)
```

### Makefile Targets

```bash
make help        # Show all available targets
make build       # Build with VERSION (required)
make build-dev   # Build for development
make clean       # Remove built binaries
make test        # Run all tests
```

## Releasing

Releases are created via a GitHub Actions workflow that builds the CLI for all major platforms and publishes a GitHub Release.

### Workflow Trigger

The release workflow is **manually triggered** using the GitHub Actions UI:

1. **Navigate to Actions**: Go to the repository → Actions tab
2. **Select workflow**: Click "Release" in the workflows list
3. **Run workflow**: Click "Run workflow" button
4. **Enter version**: Provide version in format `vX.Y.Z` (e.g., `v1.0.0`)
5. **Confirm**: Click "Run workflow" to start the build

### Workflow Behavior

The workflow performs the following steps:

1. **Validate VERSION format**: Ensures version is in `vX.Y.Z` format (fails fast if invalid)
2. **Build for 3 platforms**: ubuntu-latest, macos-latest, windows-latest (in parallel)
3. **Inject version via ldflags**: Sets version in each binary
4. **Set executable permissions**: Ensures Unix binaries have executable bit
5. **Verify builds**: Runs `osdd version` on each platform to confirm success
6. **Upload artifacts**: Stores binaries with platform-specific names
7. **Create release**: Publishes GitHub Release **only after all builds succeed**

**Atomic Release Guarantee**: The release is only published if ALL platform builds succeed. If any build fails, no release is created, preventing partial/incomplete releases.

### Artifact Naming

Release artifacts are named with platform identifiers:

- `osdd-linux-x64` - Linux binary
- `osdd-macos-x64` - macOS binary
- `osdd-windows-x64.exe` - Windows binary

### Required Permissions

The workflow requires the following repository permissions (configured in Settings → Actions → General):

- **Workflow permissions**: "Read and write permissions"
- **Allow GitHub Actions to create and approve pull requests**: Enabled (optional)

These permissions allow the workflow to:
- Create releases and tags (`contents: write`)
- Upload release assets

The `GITHUB_TOKEN` is automatically provided by GitHub Actions and scoped to the repository.

## Development

### Project Structure

```
osdd-cli/
├── .github/
│   └── workflows/
│       └── release.yml       # Multi-platform release workflow
├── cmd/
│   └── osdd/
│       └── main.go           # CLI entrypoint with Cobra wiring
├── internal/
│   ├── version/
│   │   ├── version.go        # Version variable and command
│   │   └── version_test.go   # Unit tests
│   └── ui/
│       └── renderer.go       # ASCII art rendering with animation
├── build.sh                  # Build script with version injection
├── Makefile                  # Build targets
├── go.mod                    # Go module definition
├── go.sum                    # Dependency checksums
├── LICENSE                   # License file
└── README.md                 # This file
```

### Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework for command structure and routing |
| `github.com/charmbracelet/lipgloss` | Terminal styling and color rendering |
| `github.com/muesli/termenv` | Terminal capability detection |
| `golang.org/x/term` | TTY detection |

Install dependencies:

```bash
go mod download
go mod tidy
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests for specific package
go test ./internal/version/

# Run tests with coverage
go test -cover ./...
```

### Terminal Capability Detection

The UI package detects terminal capabilities using three checks:

1. **TTY detection**: Checks if stdout is a terminal (not a pipe or redirect)
2. **NO_COLOR environment variable**: Respects user preference to disable colors
3. **Color profile**: Requires ANSI256+ support for rainbow gradient

If any check fails, the CLI falls back to plain ASCII art without colors.

### Rainbow Color Details

- **Color scheme**: 7-color rainbow gradient (red → orange → yellow → green → cyan → blue → purple)
- **Display**: Static colored text (no animation)
- **Color values**: #FF0000, #FF7F00, #FFFF00, #00FF00, #00FFFF, #0000FF, #8B00FF
- **Gradient mapping**: Colors distributed evenly across the ASCII art width
- **Panic recovery**: Rendering errors trigger graceful fallback to plain ASCII

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

### Development Workflow

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Build and verify: `make build-dev && ./osdd version`
6. Submit a pull request

## License

This project is licensed under the terms found in the [LICENSE](LICENSE) file.

## Support

For issues, questions, or feature requests, please file an issue on the [GitHub Issues page](https://github.com/opensdd/osdd-cli/issues).

For more information about OpenSDD, visit [https://opensdd.ai](https://opensdd.ai).
