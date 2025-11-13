# Project Context

## Purpose

**minit** is a lightweight init daemon (PID 1) designed specifically for containers. It provides process management, zombie process reaping, and flexible unit-based configuration for running services inside containers. The project aims to make container initialization simple and powerful without requiring full-featured init systems like systemd.

### Key Goals
- Serve as PID 1 in container environments with proper zombie process handling
- Support multiple unit types: `render` (template rendering), `once` (one-shot commands), `daemon` (long-running services), and `cron` (scheduled tasks)
- Provide flexible configuration through YAML files, environment variables, and command arguments
- Enable advanced container features: resource limits (ulimit), kernel parameters (sysctl), THP configuration, and built-in WebDAV server
- Maintain simplicity and minimal resource footprint

## Tech Stack

### Core Technologies
- **Go 1.24+** - Primary programming language with toolchain go1.24.5
- **YAML** - Configuration format for unit files (gopkg.in/yaml.v3)
- **Cron** - Scheduling library (github.com/robfig/cron/v3)

### Key Dependencies
- `github.com/yankeguo/rg` (v1.3.1) - Error handling utilities
- `github.com/robfig/cron/v3` (v3.0.1) - Cron scheduling with extended syntax
- `golang.org/x/text` (v0.31.0) - Text encoding/transcoding (GBK, GB18030)
- `golang.org/x/net` (v0.47.0) - WebDAV server support
- `golang.org/x/sys` (v0.38.0) - System calls for Linux-specific features

### Testing
- `github.com/stretchr/testify` (v1.11.1) - Test assertions and suite support

### Build & Deployment
- **Docker** - Multi-stage builds with busybox base image
- **CGO_ENABLED=0** - Static binary compilation for portability
- **Cocogitto (cog)** - Conventional commits and changelog management

## Project Conventions

### Code Style

#### Package Structure
- **Internal packages** follow `m*` naming convention (munit, mexec, mlog, menv, etc.)
- Each internal package is focused and self-contained with clear responsibilities
- Test files use `_test.go` suffix and live alongside implementation
- Testdata directories contain fixtures for integration tests

#### Naming Conventions
- **Packages**: Short, lowercase, single-word when possible (e.g., `munit`, `mexec`)
- **Types**: PascalCase (e.g., `Runner`, `Manager`, `Unit`)
- **Functions**: camelCase for private, PascalCase for exported
- **Constants**: PascalCase for exported, camelCase for private
- **Environment variables**: `MINIT_` prefix for all configuration

#### Code Organization
- Each package has clear entry points (e.g., `Load()`, `Create()`, `Setup()`)
- Use constructor functions returning `(result, error)` or `(result1, result2, error)`
- Leverage `rg.Must`, `rg.Must0`, `rg.Must2` for error handling
- Defer error handling with `defer rg.Guard(&err)`

### Architecture Patterns

#### Modular Internal Packages
The codebase is organized into focused internal packages:

- **munit**: Unit definition, loading, and filtering
  - Load units from files (`/etc/minit.d`), environment variables, and command args
  - Support multiple unit types (render, once, daemon, cron)
  - Enable/disable filtering with allowlist/denylist mode

- **mrunners**: Runner implementations for different unit types
  - Convert units to executable runners
  - Split runners into "short" (render, blocking once) and "long" (daemon, cron, non-blocking once)
  - Manage execution lifecycle and error handling

- **mexec**: Process execution manager
  - Track child processes for signal propagation
  - Handle zombie process reaping (when running as PID 1)
  - Platform-specific implementations (Linux vs non-Linux)

- **mlog**: Logging infrastructure
  - Rotating file logs with automatic cleanup
  - Per-unit log files when `MINIT_LOG_DIR` is set
  - Charset transcoding support (UTF-8, GBK, GB18030)

- **menv**: Environment variable handling
  - Parse and construct environment from multiple sources
  - Template rendering support for `MINIT_ENV_*` prefixed variables
  - Merge environment maps with precedence rules

- **mtmpl**: Template execution
  - Go template engine with custom functions
  - Functions generated from `funcs.gen.py` script
  - Support for file and directory rendering

- **msetups**: System setup tasks
  - Resource limits (rlimit)
  - Kernel parameters (sysctl)
  - Transparent Huge Pages (THP)
  - WebDAV server
  - Zombie reaping registration
  - Banner display

- **merrs**: Error handling utilities
  - Aggregate multiple errors
  - Error wrapping and unwrapping

- **pkg/shellquote**: Shell command quoting/unquoting utilities

#### Execution Flow
1. **Initialization**: Load configuration, parse environment, setup logging
2. **Setup phase**: Run system setups (banner, rlimits, sysctl, THP, WebDAV, zombies)
3. **Unit loading**: Load and filter units from all sources (files → env → args)
4. **Runner creation**: Convert units to runners, split into short/long runners
5. **Short execution**: Run render and blocking once units sequentially
6. **Long execution**: Start daemon and cron units concurrently
7. **Signal handling**: Wait for SIGTERM/SIGINT or critical errors
8. **Graceful shutdown**: Cancel context, wait 3 seconds, propagate signals, wait for cleanup

### Testing Strategy

#### Test Coverage
- Each internal package has comprehensive test files (`*_test.go`)
- Unit tests for individual functions and methods
- Integration tests using testdata fixtures
- Table-driven tests for multiple scenarios

#### Test Conventions
- Use `testify/assert` and `testify/require` for assertions
- Place test fixtures in `testdata/` subdirectories
- Mock external dependencies when necessary
- Test error cases and edge conditions
- Maintain test isolation (no shared state between tests)

#### Running Tests
```bash
go test ./...                    # Run all tests
go test -v ./internal/munit      # Run specific package tests
go test -cover ./...             # Run with coverage
```

### Git Workflow

#### Commit Conventions
- **Conventional Commits** with scope (enforced by user preference)
- Format: `<type>(<scope>): <description>`
- Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`
- Examples:
  - `feat(munit): add support for replicas in cron units`
  - `fix(mexec): handle zombie processes on non-Linux systems`
  - `docs(readme): update installation instructions`

#### Changelog Management
- Automated via **Cocogitto (cog)**
- Configuration in `cog.toml`
- Changelog maintained in `CHANGELOG.md`
- Version tags prefixed with `v` (e.g., `v1.2.3`)

#### Branching Strategy
- `main` branch is the primary development branch
- Feature branches for new capabilities
- Tag releases with semantic versioning

## Domain Context

### Container Init Systems
- **PID 1 responsibilities**: The first process (PID 1) in a container/system has special responsibilities:
  - Must reap zombie processes (wait on orphaned child processes)
  - Receives signals from the container runtime
  - Should handle graceful shutdown
  - Cannot be killed except by special signals

### Unit Types & Lifecycle
- **render**: Template rendering, executes first, used for configuration generation
- **once**: One-shot commands, can be blocking (default) or non-blocking
- **daemon**: Long-running services that restart on failure
- **cron**: Scheduled tasks with cron expressions, support extended syntax

### Unit Loading Order
1. Source order: files → environment → arguments
2. Type order: render → once → daemon/cron
3. Override with `order` field (negative = earlier, positive = later)

### Container Privileges
- **Normal mode**: Standard container, limited system access
- **Privileged mode**: Required for rlimit, sysctl, THP features
- **Host /sys mount**: Required for THP configuration

## Important Constraints

### Platform Support
- Primary target: **Linux containers**
- Platform-specific code isolated in `*_linux.go` and `*_nolinux.go` files
- Zombie reaping only functional on Linux (no-op on other platforms)
- Resource limits (rlimit), sysctl, THP require privileged containers

### Runtime Requirements
- Must be run as **PID 1** for proper signal handling and zombie reaping
- CGO disabled (`CGO_ENABLED=0`) for static binary compilation
- Minimal base image (busybox) for production containers

### Breaking Changes
- Configuration format (YAML schema) changes require migration path
- Environment variable naming must maintain `MINIT_` prefix
- Unit type changes affect existing deployments
- Signal handling behavior must remain predictable

### Performance Considerations
- Minimize memory footprint (init system runs for container lifetime)
- Fast startup time (containers should start quickly)
- Efficient log rotation (avoid disk space exhaustion)
- Graceful shutdown within reasonable timeframe (3s delay + cleanup)

## External Dependencies

### Core Dependencies
- **Go standard library**: Core functionality, no external runtime dependencies
- **Cron library** (robfig/cron): Industry-standard cron expression parser
- **Text encoding** (golang.org/x/text): Character set transcoding for legacy systems
- **WebDAV** (golang.org/x/net/webdav): Optional built-in file server

### Build-Time Dependencies
- **Docker**: Multi-stage builds for minimal images
- **Go toolchain**: Version 1.24+ with go1.24.5 toolchain
- **Cocogitto**: Changelog and version management (optional)

### Container Runtime
- Compatible with Docker, containerd, Kubernetes, and other OCI runtimes
- No specific runtime dependencies beyond standard container features
- Privileged mode optional (only for advanced features)

### File System Expectations
- `/etc/minit.d/`: Default unit file directory (configurable via `MINIT_UNIT_DIR`)
- `/etc/banner.minit.txt`: Optional startup banner file
- `MINIT_LOG_DIR`: Optional log directory (no file logging by default)
- `/sys/`: Host mount required for THP configuration (privileged only)
