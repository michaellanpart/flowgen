package models

import "time"

// FlowEntity represents the base entity with common properties
type FlowEntity struct {
	ID          string                 `json:"id" yaml:"id"`
	Name        string                 `json:"name" yaml:"name"`
	Description *string                `json:"description,omitempty" yaml:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Tags        []string               `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// Position represents X,Y coordinates
type Position struct {
	X float64 `json:"x" yaml:"x"`
	Y float64 `json:"y" yaml:"y"`
}

// Dimensions represents width and height
type Dimensions struct {
	Width  float64 `json:"width" yaml:"width"`
	Height float64 `json:"height" yaml:"height"`
}

// Style represents visual styling options
type Style struct {
	Fill            *string  `json:"fill,omitempty" yaml:"fill,omitempty"`
	Stroke          *string  `json:"stroke,omitempty" yaml:"stroke,omitempty"`
	StrokeWidth     *float64 `json:"strokeWidth,omitempty" yaml:"strokeWidth,omitempty"`
	StrokeDasharray *string  `json:"strokeDasharray,omitempty" yaml:"strokeDasharray,omitempty"`
	Opacity         *float64 `json:"opacity,omitempty" yaml:"opacity,omitempty"`
	FontSize        *float64 `json:"fontSize,omitempty" yaml:"fontSize,omitempty"`
	FontFamily      *string  `json:"fontFamily,omitempty" yaml:"fontFamily,omitempty"`
	FontWeight      *string  `json:"fontWeight,omitempty" yaml:"fontWeight,omitempty"`
	TextColor       *string  `json:"textColor,omitempty" yaml:"textColor,omitempty"`
}

// NodeType represents different types of nodes
type NodeType string

const (
	NodeTypeProcess    NodeType = "process"
	NodeTypeDecision   NodeType = "decision"
	NodeTypeStart      NodeType = "start"
	NodeTypeEnd        NodeType = "end"
	NodeTypeSubprocess NodeType = "subprocess"
	NodeTypeData       NodeType = "data"
	NodeTypeExternal   NodeType = "external"
	NodeTypeCustom     NodeType = "custom"
)

// ConnectionType represents different types of connections
type ConnectionType string

const (
	ConnectionTypeSequence    ConnectionType = "sequence"
	ConnectionTypeConditional ConnectionType = "conditional"
	ConnectionTypeDataFlow    ConnectionType = "data_flow"
	ConnectionTypeAssociation ConnectionType = "association"
	ConnectionTypeComposition ConnectionType = "composition"
	ConnectionTypeAggregation ConnectionType = "aggregation"
)

// JiraIntegration represents Jira integration data
type JiraIntegration struct {
	IssueKey   *string `json:"issueKey,omitempty" yaml:"issueKey,omitempty"`
	ProjectKey *string `json:"projectKey,omitempty" yaml:"projectKey,omitempty"`
}

// Integrations represents external system integrations
type Integrations struct {
	Jira   *JiraIntegration       `json:"jira,omitempty" yaml:"jira,omitempty"`
	Custom map[string]interface{} `json:"custom,omitempty" yaml:"custom,omitempty"`
}

// FlowNode represents a node in the flow diagram
type FlowNode struct {
	FlowEntity   `yaml:",inline"`
	Type         NodeType      `json:"type" yaml:"type"`
	Position     Position      `json:"position" yaml:"position"`
	Dimensions   *Dimensions   `json:"dimensions,omitempty" yaml:"dimensions,omitempty"`
	Style        *Style        `json:"style,omitempty" yaml:"style,omitempty"`
	DrillDown    *string       `json:"drillDown,omitempty" yaml:"drillDown,omitempty"`
	Integrations *Integrations `json:"integrations,omitempty" yaml:"integrations,omitempty"`
}

// FlowEdge represents an edge/connection in the flow diagram
type FlowEdge struct {
	FlowEntity `yaml:",inline"`
	Type       ConnectionType `json:"type" yaml:"type"`
	From       string         `json:"from" yaml:"from"`
	To         string         `json:"to" yaml:"to"`
	Condition  *string        `json:"condition,omitempty" yaml:"condition,omitempty"`
	Style      *Style         `json:"style,omitempty" yaml:"style,omitempty"`
	Waypoints  []Position     `json:"waypoints,omitempty" yaml:"waypoints,omitempty"`
}

// LayoutDirection represents diagram layout direction
type LayoutDirection string

const (
	LayoutDirectionTopBottom LayoutDirection = "top-bottom"
	LayoutDirectionBottomTop LayoutDirection = "bottom-top"
	LayoutDirectionLeftRight LayoutDirection = "left-right"
	LayoutDirectionRightLeft LayoutDirection = "right-left"
)

// LayoutSpacing represents spacing configuration
type LayoutSpacing struct {
	Node *float64 `json:"node,omitempty" yaml:"node,omitempty"`
	Rank *float64 `json:"rank,omitempty" yaml:"rank,omitempty"`
}

// Layout represents diagram layout configuration
type Layout struct {
	Direction *LayoutDirection `json:"direction,omitempty" yaml:"direction,omitempty"`
	Spacing   *LayoutSpacing   `json:"spacing,omitempty" yaml:"spacing,omitempty"`
}

// FlowDiagram represents a complete flow diagram
type FlowDiagram struct {
	FlowEntity `yaml:",inline"`
	Version    string     `json:"version" yaml:"version"`
	Nodes      []FlowNode `json:"nodes" yaml:"nodes"`
	Edges      []FlowEdge `json:"edges" yaml:"edges"`
	Layout     *Layout    `json:"layout,omitempty" yaml:"layout,omitempty"`
	Parent     *string    `json:"parent,omitempty" yaml:"parent,omitempty"`
	Children   []string   `json:"children,omitempty" yaml:"children,omitempty"`
	Created    time.Time  `json:"created" yaml:"created"`
	Updated    time.Time  `json:"updated" yaml:"updated"`
	FilePath   string     `json:"filePath,omitempty" yaml:"-"` // Internal use only
}

// ValidationError represents a validation error
type ValidationError struct {
	Path    string      `json:"path"`
	Message string      `json:"message"`
	Code    string      `json:"code"`
	Value   interface{} `json:"value,omitempty"`
}

// ValidationResult represents the result of diagram validation
type ValidationResult struct {
	Valid    bool              `json:"valid"`
	Errors   []ValidationError `json:"errors"`
	Warnings []ValidationError `json:"warnings,omitempty"`
}

// SearchResult represents a search result item
type SearchResult struct {
	Diagram   FlowDiagram `json:"diagram"`
	Score     float64     `json:"score"`
	MatchType string      `json:"matchType"` // "name", "description", "tag", "node", etc.
}

// NodeSearchResult represents a node search result
type NodeSearchResult struct {
	Node      FlowNode    `json:"node"`
	DiagramID string      `json:"diagramId"`
	Diagram   FlowDiagram `json:"diagram"`
	Score     float64     `json:"score"`
	MatchType string      `json:"matchType"`
}
