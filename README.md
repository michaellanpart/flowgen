# FlowGen

A YAML-based flow diagramming library and toolset with hierarchical drill-down capabilities and enterprise integration support.

## Overview

FlowGen is a modern alternative to traditional diagramming tools like Mermaid.js, designed specifically for enterprise software architecture and product modeling. Instead of custom markup languages, FlowGen uses structured YAML files to define flow diagrams, enabling better version control, validation, and integration with existing toolchains.

## Architecture

- **`frontend/`** - Static web UI (Bootstrap) served by the backend. Contains all Node tooling (`package.json`, lockfile, `node_modules/`).
- **`backend/`** - Golang REST API for diagram management and enterprise integrations; also serves static frontend assets.
- **`backend/schemas/`** - YAML schemas and validation rules (backend-owned)
- **`examples/`** - Sample diagrams and use cases
- **`docs/`** - Documentation and API references

## Key Features

- 🎯 YAML-first approach for better developer experience
- 🔍 Hierarchical drill-down from high-level to detailed views
- 🔗 Enterprise integrations (Jira, etc.) for change tracking
- 🏗️ Clean, extensible architecture following design patterns
- 📱 Framework-agnostic rendering engine
- 🔧 TypeScript support with full type safety

## Getting Started

Backend (Go):

```bash
# From repository root
cd backend
go run ./cmd
```

Frontend (Node toolchain lives under `frontend/`):

```bash
cd frontend
npm install
# If you add any build tooling later, run it from here, e.g.:
# npm run build
```

## License

MIT License