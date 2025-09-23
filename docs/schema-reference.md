# FlowGen YAML Schema Reference

This document provides a comprehensive reference for the FlowGen YAML schema used to define flow diagrams.

## Document Structure

```yaml
# Required fields
id: string              # Unique diagram identifier
name: string           # Human-readable diagram name  
version: string        # Semantic version (e.g., "1.0.0")
nodes: array           # Array of node definitions
edges: array           # Array of edge definitions

# Optional fields
description: string    # Diagram description
metadata: object       # Additional metadata
tags: array           # Array of string tags
layout: object        # Layout configuration
parent: string        # Parent diagram ID (for hierarchy)
children: array       # Array of child diagram IDs
```

## Node Definition

```yaml
nodes:
  - id: string                    # Required: Unique node ID
    name: string                  # Required: Display name
    type: enum                    # Required: Node type
    position:                     # Required: Position
      x: number
      y: number
    
    # Optional fields
    description: string           # Node description
    dimensions:                   # Node size
      width: number
      height: number
    style:                        # Visual styling
      fill: string
      stroke: string
      strokeWidth: number
      # ... more style properties
    drillDown: string            # Child diagram ID
    metadata: object             # Additional data
    tags: array                  # String tags
    integrations:                # External integrations
      jira:
        issueKey: string
        projectKey: string
```

### Node Types

| Type | Description | Visual Shape | Use Case |
|------|-------------|--------------|----------|
| `start` | Process starting point | Oval/Circle | Entry points |
| `end` | Process ending point | Oval/Circle | Exit points |
| `process` | General processing step | Rectangle | Business logic |
| `decision` | Decision point | Diamond | Conditional flows |
| `subprocess` | Reference to sub-process | Rectangle with border | Drill-down capability |
| `data` | Data storage/retrieval | Cylinder/Database | Data operations |
| `external` | External system | Rectangle with shadow | Third-party systems |
| `custom` | Custom node type | Configurable | Special cases |

## Edge Definition

```yaml
edges:
  - id: string                    # Required: Unique edge ID
    name: string                  # Required: Display name
    type: enum                    # Required: Edge type
    from: string                  # Required: Source node ID
    to: string                    # Required: Target node ID
    
    # Optional fields
    description: string           # Edge description
    condition: string             # Condition for conditional edges
    style:                        # Visual styling
      stroke: string
      strokeWidth: number
      strokeDasharray: string
    waypoints:                    # Custom routing points
      - x: number
        y: number
    metadata: object              # Additional data
    tags: array                   # String tags
```

### Edge Types

| Type | Description | Visual Style | Use Case |
|------|-------------|--------------|----------|
| `sequence` | Sequential flow | Solid line | Normal progression |
| `conditional` | Conditional flow | Dashed line | Decision outcomes |
| `data_flow` | Data transfer | Dotted line | Data movement |
| `association` | Loose association | Thin line | Relationships |
| `composition` | Strong composition | Thick line | Ownership |
| `aggregation` | Aggregation | Dashed thick | Part-of relationships |

## Layout Configuration

```yaml
layout:
  direction: enum               # Layout direction
  spacing:
    node: number               # Space between nodes
    rank: number               # Space between levels
```

### Layout Directions

- `top-bottom` - Vertical flow from top to bottom
- `bottom-top` - Vertical flow from bottom to top  
- `left-right` - Horizontal flow from left to right
- `right-left` - Horizontal flow from right to left

## Style Properties

### Node Styles
```yaml
style:
  fill: string                  # Fill color (#hex, rgb(), rgba())
  stroke: string                # Border color
  strokeWidth: number           # Border width (0-20)
  strokeDasharray: string       # Border dash pattern
  opacity: number               # Opacity (0-1)
  cornerRadius: number          # Corner rounding
  fontSize: number              # Text size (8-72)
  fontFamily: string            # Font family
  fontWeight: string            # Font weight
  textColor: string             # Text color
```

### Edge Styles
```yaml
style:
  stroke: string                # Line color
  strokeWidth: number           # Line width (0-20)
  strokeDasharray: string       # Line dash pattern
  opacity: number               # Opacity (0-1)
```

## Metadata Object

The metadata object can contain any additional information:

```yaml
metadata:
  author: string
  department: string
  created: string (ISO date)
  lastModified: string (ISO date)
  version: string
  approvers: array
  status: string
  priority: string
  # Any custom fields
```

## Integration Objects

### Jira Integration
```yaml
integrations:
  jira:
    issueKey: string            # Format: PROJECT-123
    projectKey: string          # Format: PROJECT
```

### Custom Integrations
```yaml
integrations:
  custom:
    system: string
    reference: string
    # Any custom fields
```

## Validation Rules

### Required Fields
- Diagram: `id`, `name`, `version`, `nodes`, `edges`
- Node: `id`, `name`, `type`, `position`
- Edge: `id`, `name`, `type`, `from`, `to`

### Format Constraints
- IDs: Must match pattern `^[a-zA-Z][a-zA-Z0-9_-]*$`
- Versions: Must match semantic versioning `^\\d+\\.\\d+\\.\\d+$`
- Colors: Must be valid CSS colors
- Jira issue keys: Must match `^[A-Z]+-\\d+$`
- Jira project keys: Must match `^[A-Z]+$`

### Business Rules
- All node and edge IDs must be unique within a diagram
- Edge `from` and `to` must reference existing node IDs
- No self-referencing edges (from = to)
- Parent-child relationships cannot form cycles
- DrillDown references must point to existing child diagrams

## Example Schema Usage

### Basic Flow
```yaml
id: "basic_flow"
name: "Basic Process"
version: "1.0.0"
nodes:
  - id: "start"
    name: "Start"
    type: "start"
    position: { x: 100, y: 50 }
  - id: "process"
    name: "Do Work"
    type: "process"
    position: { x: 100, y: 150 }
  - id: "end"
    name: "End"
    type: "end"
    position: { x: 100, y: 250 }
edges:
  - id: "start_to_process"
    name: "Begin"
    type: "sequence"
    from: "start"
    to: "process"
  - id: "process_to_end"
    name: "Complete"
    type: "sequence"
    from: "process"
    to: "end"
```

### Advanced Flow with Styling
```yaml
id: "styled_flow"
name: "Styled Process"
version: "1.0.0"
tags: ["styled", "example"]
metadata:
  author: "Developer"
  created: "2024-01-15"
nodes:
  - id: "start"
    name: "Start"
    type: "start"
    position: { x: 100, y: 50 }
    style:
      fill: "#27ae60"
      stroke: "#229954"
      textColor: "#ffffff"
  - id: "decision"
    name: "Check Status"
    type: "decision"
    position: { x: 75, y: 150 }
    style:
      fill: "#f39c12"
      stroke: "#e67e22"
  - id: "process_yes"
    name: "Process A"
    type: "process"
    position: { x: 25, y: 250 }
    integrations:
      jira:
        projectKey: "PROJ"
        issueKey: "PROJ-123"
  - id: "process_no"
    name: "Process B"
    type: "process"
    position: { x: 175, y: 250 }
edges:
  - id: "start_to_decision"
    name: "Check"
    type: "sequence"
    from: "start"
    to: "decision"
  - id: "decision_yes"
    name: "Yes"
    type: "conditional"
    from: "decision"
    to: "process_yes"
    condition: "status === 'active'"
  - id: "decision_no"
    name: "No"
    type: "conditional"
    from: "decision"
    to: "process_no"
    condition: "status !== 'active'"
    style:
      stroke: "#e74c3c"
      strokeDasharray: "5,5"
layout:
  direction: "top-bottom"
  spacing:
    node: 50
    rank: 100
```

## Schema Evolution

The FlowGen schema supports versioning to ensure backward compatibility:

- **Major versions** (x.0.0) - Breaking changes
- **Minor versions** (0.x.0) - New features, backward compatible
- **Patch versions** (0.0.x) - Bug fixes, backward compatible

Current schema version: **1.0.0**