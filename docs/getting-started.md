# FlowGen Getting Started Guide

Welcome to FlowGen, a YAML-based flow diagramming library designed for enterprise software architecture and product modeling.

## Quick Start

### 1. Installation

```bash
# Install the core library
npm install @flowgen/core @flowgen/renderer

# Or install the full suite
npm install @flowgen/core @flowgen/renderer
```

### 2. Basic Usage

```typescript
import { FlowParser } from '@flowgen/core';
import { SVGRenderer } from '@flowgen/renderer';

// Parse a YAML diagram
const parser = new FlowParser();
const result = parser.parse(yamlContent);

if (result.success && result.data) {
  // Render the diagram
  const container = document.getElementById('diagram-container');
  const renderer = new SVGRenderer(container);
  
  await renderer.render(result.data);
}
```

### 3. Backend API

Start the backend server:

```bash
cd backend
go run main.go
```

The API will be available at `http://localhost:8080`

## Core Concepts

### Diagrams
FlowGen diagrams are defined in YAML format with the following structure:

```yaml
id: "my_diagram"
name: "My Process Flow"
version: "1.0.0"
nodes: [...]
edges: [...]
```

### Nodes
Nodes represent steps, decisions, or states in your process:

```yaml
nodes:
  - id: "start"
    name: "Start Process"
    type: "start"
    position: { x: 100, y: 50 }
    dimensions: { width: 80, height: 40 }
```

**Node Types:**
- `start` - Process starting point
- `end` - Process ending point
- `process` - General processing step
- `decision` - Decision point with multiple outcomes
- `subprocess` - Reference to another process (supports drill-down)
- `data` - Data storage or retrieval
- `external` - External system interaction
- `custom` - Custom node type

### Edges
Edges define connections between nodes:

```yaml
edges:
  - id: "edge1"
    name: "Next Step"
    type: "sequence"
    from: "start"
    to: "process1"
```

**Edge Types:**
- `sequence` - Sequential flow
- `conditional` - Conditional flow with conditions
- `data_flow` - Data passing between nodes
- `association` - Loose association
- `composition` - Strong composition relationship
- `aggregation` - Aggregation relationship

### Hierarchical Drill-Down
Create detailed views by linking diagrams:

```yaml
# Parent diagram node
- id: "complex_process"
  name: "Complex Process"
  type: "subprocess"
  drillDown: "detailed_process_id"

# Child diagram
id: "detailed_process_id"
parent: "parent_diagram_id"
```

### Enterprise Integrations

#### Jira Integration
Link nodes to Jira issues for change tracking:

```yaml
nodes:
  - id: "approval_step"
    name: "Approval Required"
    type: "process"
    integrations:
      jira:
        projectKey: "PROJ"
        issueKey: "PROJ-123"
```

## API Reference

### Core Library

#### FlowParser
```typescript
const parser = new FlowParser(options);
const result = parser.parse(yamlContent);
```

#### FlowValidator
```typescript
const validator = new FlowValidator(options);
const result = validator.validate(diagram);
```

### Renderer

#### SVGRenderer
```typescript
const renderer = new SVGRenderer(container, options);
await renderer.render(diagram);
```

**Options:**
- `width` - Canvas width
- `height` - Canvas height  
- `theme` - Visual theme
- `enableInteractivity` - Enable mouse/touch interactions
- `enableAnimations` - Enable animations

### Backend API

#### Diagram Operations
- `GET /api/v1/diagrams` - List all diagrams
- `POST /api/v1/diagrams` - Create new diagram
- `GET /api/v1/diagrams/:id` - Get specific diagram
- `PUT /api/v1/diagrams/:id` - Update diagram
- `DELETE /api/v1/diagrams/:id` - Delete diagram
- `POST /api/v1/diagrams/:id/validate` - Validate diagram

#### Hierarchy Operations
- `GET /api/v1/hierarchy/:id/children` - Get child diagrams
- `GET /api/v1/hierarchy/:id/parent` - Get parent diagram
- `POST /api/v1/hierarchy/:id/link` - Link diagrams

#### Search
- `GET /api/v1/search/diagrams?q=query&tags=tag1,tag2` - Search diagrams
- `GET /api/v1/search/nodes?q=query&type=process` - Search nodes

## Example Diagrams

Check the `examples/` directory for sample diagrams:

- `simple-user-flow.yaml` - Basic user registration flow
- `account-creation-detail.yaml` - Detailed account creation process
- `payment-processing.yaml` - Enterprise payment processing

## Schema Validation

FlowGen includes comprehensive YAML schema validation:

- **Syntax validation** - YAML structure and types
- **Semantic validation** - Business logic rules
- **Reference validation** - Node and edge relationships
- **Best practices** - Recommendations for better diagrams

## Theming

Customize the visual appearance:

```typescript
import { DefaultTheme } from '@flowgen/renderer';

const customTheme = {
  ...DefaultTheme,
  colors: {
    ...DefaultTheme.colors,
    primary: '#your-color',
  },
};

renderer.setTheme(customTheme);
```

## Performance

FlowGen is optimized for large diagrams:

- **Lazy rendering** - Only visible elements are rendered
- **Virtual scrolling** - Handle thousands of nodes
- **Incremental updates** - Only changed elements are re-rendered
- **Memory management** - Automatic cleanup of unused resources

## Development

### Building from Source

```bash
# Install dependencies
npm install

# Build all packages
npm run build

# Run tests
npm test

# Start development server
npm run dev
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes following our coding standards
4. Add tests for new functionality
5. Submit a pull request

## Support

- **Documentation**: [https://flowgen.dev/docs](https://flowgen.dev/docs)
- **Issues**: [GitHub Issues](https://github.com/michaellanpart/flowgen/issues)
- **Discussions**: [GitHub Discussions](https://github.com/michaellanpart/flowgen/discussions)

## License

MIT License - see LICENSE file for details.