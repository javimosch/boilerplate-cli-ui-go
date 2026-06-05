# Boilerplate CLI UI Go

**In 2026, every app is a CLI, even UI based.**

This is a boilerplate for crafting UI-based applications packed as single binaries, designed to work as plugins for [SuperCLI](https://github.com/javimosch/supercli). It demonstrates how modern applications should be built: CLI-first with optional web interfaces, compiled to portable binaries.

| Stack | Repo | Binary |
|-------|------|--------|
| **Go + inline HTML** | **boilerplate-cli-ui-go** | **~5MB** |
| Go + Vue 3 CDN | [boilerplate-cli-ui-go-v2-vue](https://github.com/javimosch/boilerplate-cli-ui-go-v2-vue) | ~5MB |
| Go + React 18 CDN | [boilerplate-cli-ui-go-v2-react](https://github.com/javimosch/boilerplate-cli-ui-go-v2-react) | ~5MB |
| Deno + vanilla JS | [boilerplate-cli-ui-deno](https://github.com/javimosch/boilerplate-cli-ui-deno) | ~76MB |
| Node.js + vanilla JS | [boilerplate-cli-ui-node](https://github.com/javimosch/boilerplate-cli-ui-node) | ~123MB |
| Python + React CDN | [boilerplate-cli-ui-python](https://github.com/javimosch/boilerplate-cli-ui-python) | ~10MB |
| Rust + vanilla JS | [boilerplate-cli-ui-rust](https://github.com/javimosch/boilerplate-cli-ui-rust) | ~1.1MB |
| .NET 8 + Vue 3 | [boilerplate-cli-ui-dotnet](https://github.com/javimosch/boilerplate-cli-ui-dotnet) | ~89MB |
| C++ + Vue 3 | [boilerplate-cli-ui-cpp](https://github.com/javimosch/boilerplate-cli-ui-cpp) | ~493KB |
| Nim + Vue 3 | [boilerplate-cli-ui-nim](https://github.com/javimosch/boilerplate-cli-ui-nim) | ~364KB |
| Zig + Vue 3 | [boilerplate-cli-ui-zig](https://github.com/javimosch/boilerplate-cli-ui-zig) | ~190KB |
| Dart + Vue 3 | [boilerplate-cli-ui-dart](https://github.com/javimosch/boilerplate-cli-ui-dart) | ~6.4MB |
|| V + Vue 3 | [boilerplate-cli-ui-v](https://github.com/javimosch/boilerplate-cli-ui-v) | ~1.2MB |

## Philosophy

**CLI-Native, Web-Enabled**

Modern applications should be CLI-native by default, with web interfaces as an enhancement rather than a requirement. This approach provides:

- **Scriptable**: Perfect for automation and CI/CD
- **Portable**: Single binary works everywhere
- **Manageable**: Daemon mode for background services
- **Accessible**: Web UI when visual interaction is preferred
- **Composable**: Designed to work as SuperCLI plugins

## SuperCLI Integration

This boilerplate is specifically designed to create plugins for [SuperCLI](https://github.com/javimosch/supercli):

- **Plugin Structure**: Follows SuperCLI plugin conventions
- **Binary Size**: Optimized to ~5MB for fast distribution
- **CLI Commands**: Start/stop/status for daemon management
- **HTTP Interface**: Optional web UI for plugin configuration
- **Process Management**: Background daemon mode for long-running plugins

## Features

- **CLI-First Design**: Primary interface is command-line
- **HTTP Server**: Built-in web server for visual interface
- **Daemon Mode**: Background process with start/stop/status
- **Single Binary**: Compiles to one executable, no dependencies
- **Web UI**: Clean, modern interface for visual interaction
- **API Endpoints**: JSON API for programmatic access
- **Portable**: Works on any system with the binary

## Architecture

```
CLI (boilerplate-cli-ui-go)
├── main.go - CLI commands and flag parsing
├── server.go - HTTP server and web UI
└── daemon.go - Process management
```

## Usage

### Build

```bash
chmod +x build.sh
./build.sh
```

### CLI Commands

**Start HTTP server (foreground):**
```bash
./boilerplate-cli-ui-go-optimized start
```

**Start HTTP server on custom port:**
```bash
./boilerplate-cli-ui-go-optimized start -port 3000
```

**Start as daemon (background):**
```bash
./boilerplate-cli-ui-go-optimized start -daemon
```

**Start daemon on custom port:**
```bash
./boilerplate-cli-ui-go-optimized start -port 3000 -daemon
```

**Stop daemon:**
```bash
./boilerplate-cli-ui-go-optimized stop
```

**Check daemon status:**
```bash
./boilerplate-cli-ui-go-optimized status
```

**Show version:**
```bash
./boilerplate-cli-ui-go-optimized version
```

## Web Interface

When the server is running, access the UI at:
- `http://localhost:8080` (default port)
- `http://localhost:3000` (if started with -port 3000)

### API Endpoints

- `GET /` - Web UI
- `GET /api/status` - Server status (JSON)
- `GET /api/health` - Health check (JSON)

## Daemon Management

The daemon mode allows the HTTP server to run in the background:

- **PID File**: `/tmp/boilerplate-cli-ui-go.pid`
- **Log File**: `/tmp/boilerplate-cli-ui-go.log`
- **Process Control**: SIGTERM for graceful shutdown

## Binary Size

- **Default**: ~7.3MB
- **Optimized**: ~5.0MB (with `-ldflags "-s -w"`)

The optimized size is ideal for SuperCLI plugin distribution.

## SuperCLI Plugin Development

To use this as a SuperCLI plugin:

1. **Customize the CLI commands** for your specific use case
2. **Extend the HTTP server** with your plugin's web UI
3. **Add plugin-specific API endpoints** for configuration
4. **Build the optimized binary** for distribution
5. **Package as SuperCLI plugin** following plugin conventions

### Example Plugin Structure

```go
// Replace the greet command with your plugin's commands
case "my-plugin":
    handleMyPlugin()
case "start":
    handleStart()  // Keep daemon management
```

## Requirements

- Go 1.21+

## Examples

### Development Workflow
```bash
# Build the binary
./build.sh

# Start server in foreground for development
./boilerplate-cli-ui-go-optimized start

# In another terminal, test the API
curl http://localhost:8080/api/status

# Stop with Ctrl+C
```

### Production Workflow
```bash
# Start as daemon
./boilerplate-cli-ui-go-optimized start -daemon

# Check status
./boilerplate-cli-ui-go-optimized status

# View logs
tail -f /tmp/boilerplate-cli-ui-go.log

# Stop when done
./boilerplate-cli-ui-go-optimized stop
```

## Use Cases

- **SuperCLI Plugins**: UI-enabled plugins for SuperCLI
- **CLI Tools**: Add web interface to existing CLI tools
- **Microservices**: Small HTTP services with CLI management
- **Admin Panels**: Simple admin interfaces for system tools
- **Development**: Quick prototyping of CLI + web applications
- **Monitoring**: Status dashboards for long-running processes

## Modern App Philosophy

**CLI + Web = Perfect Combination**

This boilerplate embodies the modern application philosophy:

1. **CLI First**: Build for automation and scripting
2. **Web Enhanced**: Add visual interfaces when needed
3. **Single Binary**: Easy distribution and installation
4. **Daemon Capable**: Background services when required
5. **Plugin Ready**: Designed for ecosystem integration

**Why This Matters:**

- **DevOps Friendly**: Perfect for CI/CD pipelines
- **User Friendly**: Web UI for visual users
- **Portable**: Single binary, no runtime dependencies
- **Maintainable**: Simple architecture, easy to extend
- **Extensible**: Plugin system for ecosystem growth

## Future Enhancements

- [ ] Add configuration file support
- [ ] Add authentication for web UI
- [ ] Add HTTPS support
- [ ] Add systemd service file generation
- [ ] Add more API endpoints
- [ ] Add database integration
- [ ] Add metrics/monitoring
- [ ] Add SuperCLI plugin packaging script

## Related Projects

- [SuperCLI](https://github.com/javimosch/supercli) - Universal CLI framework
- [supercli-clis](https://github.com/jarancibia/supercli-clis) - Collection of SuperCLI plugins
- [boilerplate-cli](https://github.com/javimosch/supercli-cli-boilerplates) - Binary size benchmarks

## Documentation

### Agent-First Design

This boilerplate follows agent-first CLI design principles to ensure AI agents can effectively use and extend the tool. Key features:

- **JSON-by-default**: All commands output JSON by default
- **Semantic exit codes**: Structured error codes (80-119) for programmatic decision-making
- **Structured errors**: Error objects with recovery hints and suggestions
- **Output separation**: stdout for data, stderr for logs
- **Schema discovery**: `--schema` flag for JSON schema validation

### Documentation Files

- **AGENTS.md**: Comprehensive guide for AI agents on extending and maintaining the boilerplate
- **docs/AGENTS_FRIENDLY_TOOLS.md**: Reference document on agent-friendly CLI design principles
- **.agents/skills/**: Agent skill files for specific tasks:
  - `boilerplate-go-usage.md`: Usage guide for agents
  - `boilerplate-go-dev.md`: Development guide for agents
  - `boilerplate-go-smoke-tests.md`: Smoke testing procedures

### Getting Started for Agents

Agents should start by reading AGENTS.md to understand:
- Project structure and coding rules
- Adding new commands and features
- Error handling patterns
- Configuration management
- Testing guidelines

## License

This boilerplate is provided as-is for educational and development purposes.