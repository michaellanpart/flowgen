package services

import (
	"errors"
	"fmt"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/michaellanpart/flowgen/backend/internal/config"
	"github.com/michaellanpart/flowgen/backend/internal/models"
	"gopkg.in/yaml.v3"
)

var (
	ErrDiagramNotFound = errors.New("diagram not found")
	ErrInvalidDiagram  = errors.New("invalid diagram")
)

// DiagramService handles diagram operations
type DiagramService struct {
	cfg *config.Config
}

// NewDiagramService creates a new diagram service
func NewDiagramService() *DiagramService {
	return &DiagramService{
		cfg: config.Load(),
	}
}

// ListAll returns all diagrams
func (s *DiagramService) ListAll() ([]models.FlowDiagram, error) {
	diagrams := []models.FlowDiagram{}

	err := filepath.Walk(s.cfg.DiagramsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
			diagram, err := s.loadDiagramFromFile(path)
			if err != nil {
				// Log error but continue with other files
				fmt.Printf("Error loading diagram from %s: %v\n", path, err)
				return nil
			}
			diagrams = append(diagrams, *diagram)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan diagrams directory: %w", err)
	}

	return diagrams, nil
}

// GetByID returns a diagram by ID
func (s *DiagramService) GetByID(id string) (*models.FlowDiagram, error) {
	diagrams, err := s.ListAll()
	if err != nil {
		return nil, err
	}

	for _, diagram := range diagrams {
		if diagram.ID == id {
			return &diagram, nil
		}
	}

	return nil, ErrDiagramNotFound
}

// Create creates a new diagram
func (s *DiagramService) Create(diagram *models.FlowDiagram) (*models.FlowDiagram, error) {
	// Set timestamps
	now := time.Now()
	diagram.Created = now
	diagram.Updated = now

	// Validate diagram
	if err := s.validateDiagram(diagram); err != nil {
		return nil, err
	}

	// Generate file path
	filename := fmt.Sprintf("%s.yaml", diagram.ID)
	filePath := filepath.Join(s.cfg.DiagramsPath, filename)
	diagram.FilePath = filePath

	// Ensure directory exists
	if err := os.MkdirAll(s.cfg.DiagramsPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create diagrams directory: %w", err)
	}

	// Save to file
	if err := s.saveDiagramToFile(diagram, filePath); err != nil {
		return nil, err
	}

	return diagram, nil
}

// Update updates an existing diagram
func (s *DiagramService) Update(diagram *models.FlowDiagram) (*models.FlowDiagram, error) {
	// Check if diagram exists
	existing, err := s.GetByID(diagram.ID)
	if err != nil {
		return nil, err
	}

	// Preserve creation time and file path
	diagram.Created = existing.Created
	diagram.Updated = time.Now()
	diagram.FilePath = existing.FilePath

	// Validate diagram
	if err := s.validateDiagram(diagram); err != nil {
		return nil, err
	}

	// Save to file
	if err := s.saveDiagramToFile(diagram, diagram.FilePath); err != nil {
		return nil, err
	}

	return diagram, nil
}

// Delete deletes a diagram
func (s *DiagramService) Delete(id string) error {
	diagram, err := s.GetByID(id)
	if err != nil {
		return err
	}

	// Remove file
	if err := os.Remove(diagram.FilePath); err != nil {
		return fmt.Errorf("failed to delete diagram file: %w", err)
	}

	return nil
}

// Validate validates a diagram
func (s *DiagramService) Validate(diagram *models.FlowDiagram) (*models.ValidationResult, error) {
	result := &models.ValidationResult{
		Valid:    true,
		Errors:   []models.ValidationError{},
		Warnings: []models.ValidationError{},
	}

	// Basic validation
	if diagram.ID == "" {
		result.Errors = append(result.Errors, models.ValidationError{
			Path:    "id",
			Message: "Diagram ID is required",
			Code:    "MISSING_ID",
		})
	}

	if diagram.Name == "" {
		result.Errors = append(result.Errors, models.ValidationError{
			Path:    "name",
			Message: "Diagram name is required",
			Code:    "MISSING_NAME",
		})
	}

	if diagram.Version == "" {
		result.Errors = append(result.Errors, models.ValidationError{
			Path:    "version",
			Message: "Diagram version is required",
			Code:    "MISSING_VERSION",
		})
	}

	// Validate nodes
	nodeIDs := make(map[string]bool)
	for i, node := range diagram.Nodes {
		if node.ID == "" {
			result.Errors = append(result.Errors, models.ValidationError{
				Path:    fmt.Sprintf("nodes[%d].id", i),
				Message: "Node ID is required",
				Code:    "MISSING_NODE_ID",
			})
		} else if nodeIDs[node.ID] {
			result.Errors = append(result.Errors, models.ValidationError{
				Path:    fmt.Sprintf("nodes[%d].id", i),
				Message: fmt.Sprintf("Duplicate node ID: %s", node.ID),
				Code:    "DUPLICATE_NODE_ID",
			})
		} else {
			nodeIDs[node.ID] = true
		}

		if node.Name == "" {
			result.Errors = append(result.Errors, models.ValidationError{
				Path:    fmt.Sprintf("nodes[%d].name", i),
				Message: "Node name is required",
				Code:    "MISSING_NODE_NAME",
			})
		}
	}

	// Validate edges
	edgeIDs := make(map[string]bool)
	for i, edge := range diagram.Edges {
		if edge.ID == "" {
			result.Errors = append(result.Errors, models.ValidationError{
				Path:    fmt.Sprintf("edges[%d].id", i),
				Message: "Edge ID is required",
				Code:    "MISSING_EDGE_ID",
			})
		} else if edgeIDs[edge.ID] {
			result.Errors = append(result.Errors, models.ValidationError{
				Path:    fmt.Sprintf("edges[%d].id", i),
				Message: fmt.Sprintf("Duplicate edge ID: %s", edge.ID),
				Code:    "DUPLICATE_EDGE_ID",
			})
		} else {
			edgeIDs[edge.ID] = true
		}

		if !nodeIDs[edge.From] {
			result.Errors = append(result.Errors, models.ValidationError{
				Path:    fmt.Sprintf("edges[%d].from", i),
				Message: fmt.Sprintf("Edge references non-existent from node: %s", edge.From),
				Code:    "INVALID_FROM_NODE",
			})
		}

		if !nodeIDs[edge.To] {
			result.Errors = append(result.Errors, models.ValidationError{
				Path:    fmt.Sprintf("edges[%d].to", i),
				Message: fmt.Sprintf("Edge references non-existent to node: %s", edge.To),
				Code:    "INVALID_TO_NODE",
			})
		}
	}

	result.Valid = len(result.Errors) == 0
	return result, nil
}

// Search searches for diagrams
func (s *DiagramService) Search(query string, tags []string) ([]models.SearchResult, error) {
	diagrams, err := s.ListAll()
	if err != nil {
		return nil, err
	}

	results := []models.SearchResult{}
	query = strings.ToLower(query)

	for _, diagram := range diagrams {
		score := 0.0
		matchType := ""

		// Search in name
		if strings.Contains(strings.ToLower(diagram.Name), query) {
			score += 1.0
			matchType = "name"
		}

		// Search in description
		if diagram.Description != nil && strings.Contains(strings.ToLower(*diagram.Description), query) {
			score += 0.8
			if matchType == "" {
				matchType = "description"
			}
		}

		// Search in tags
		for _, tag := range diagram.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				score += 0.6
				if matchType == "" {
					matchType = "tag"
				}
			}
		}

		// Filter by tags if specified
		if len(tags) > 0 {
			hasAllTags := true
			for _, requiredTag := range tags {
				found := false
				for _, diagramTag := range diagram.Tags {
					if strings.EqualFold(diagramTag, requiredTag) {
						found = true
						break
					}
				}
				if !found {
					hasAllTags = false
					break
				}
			}
			if !hasAllTags {
				continue
			}
		}

		if score > 0 || len(tags) > 0 {
			results = append(results, models.SearchResult{
				Diagram:   diagram,
				Score:     score,
				MatchType: matchType,
			})
		}
	}

	return results, nil
}

// SearchNodes searches for nodes across all diagrams
func (s *DiagramService) SearchNodes(query string, nodeType string) ([]models.NodeSearchResult, error) {
	diagrams, err := s.ListAll()
	if err != nil {
		return nil, err
	}

	results := []models.NodeSearchResult{}
	query = strings.ToLower(query)

	for _, diagram := range diagrams {
		for _, node := range diagram.Nodes {
			score := 0.0
			matchType := ""

			// Filter by node type if specified
			if nodeType != "" && string(node.Type) != nodeType {
				continue
			}

			// Search in node name
			if strings.Contains(strings.ToLower(node.Name), query) {
				score += 1.0
				matchType = "name"
			}

			// Search in node description
			if node.Description != nil && strings.Contains(strings.ToLower(*node.Description), query) {
				score += 0.8
				if matchType == "" {
					matchType = "description"
				}
			}

			if score > 0 || nodeType != "" {
				results = append(results, models.NodeSearchResult{
					Node:      node,
					DiagramID: diagram.ID,
					Diagram:   diagram,
					Score:     score,
					MatchType: matchType,
				})
			}
		}
	}

	return results, nil
}

// Private helper methods

func (s *DiagramService) loadDiagramFromFile(filePath string) (*models.FlowDiagram, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var diagram models.FlowDiagram
	if err := yaml.Unmarshal(data, &diagram); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	diagram.FilePath = filePath
	return &diagram, nil
}

func (s *DiagramService) saveDiagramToFile(diagram *models.FlowDiagram, filePath string) error {
	data, err := s.marshalDiagramYAML(diagram)
	if err != nil {
		return fmt.Errorf("failed to marshal diagram to YAML: %w", err)
	}

	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (s *DiagramService) validateDiagram(diagram *models.FlowDiagram) error {
	result, err := s.Validate(diagram)
	if err != nil {
		return err
	}

	if !result.Valid {
		return fmt.Errorf("diagram validation failed: %d errors", len(result.Errors))
	}

	return nil
}

// LoadYAMLByID returns the raw YAML content for a diagram ID
func (s *DiagramService) LoadYAMLByID(id string) (string, error) {
	// Scan for file named <id>.yaml or <id>.yml in diagrams path
	candidates := []string{
		filepath.Join(s.cfg.DiagramsPath, id+".yaml"),
		filepath.Join(s.cfg.DiagramsPath, id+".yml"),
	}
	var found string
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			found = p
			break
		}
	}
	if found == "" {
		return "", ErrDiagramNotFound
	}
	b, err := os.ReadFile(found)
	if err != nil {
		return "", fmt.Errorf("failed to read yaml: %w", err)
	}
	return string(b), nil
}

// SaveYAMLByID writes YAML content to the diagram file, validating it first
func (s *DiagramService) SaveYAMLByID(id, yamlText string) error {
	// Parse YAML to ensure validity and that ID matches
	var diagram models.FlowDiagram
	if err := yaml.Unmarshal([]byte(yamlText), &diagram); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}
	if diagram.ID == "" {
		// If no ID in YAML, set from path
		diagram.ID = id
	} else if diagram.ID != id {
		return fmt.Errorf("diagram id mismatch: yaml has '%s', path has '%s'", diagram.ID, id)
	}

	// Validate semantic model
	if err := s.validateDiagram(&diagram); err != nil {
		return err
	}

	// Determine file path (prefer .yaml)
	if err := os.MkdirAll(s.cfg.DiagramsPath, 0o755); err != nil {
		return fmt.Errorf("failed to ensure diagrams dir: %w", err)
	}
	filePath := filepath.Join(s.cfg.DiagramsPath, id+".yaml")

	// Marshal back to canonical YAML to keep formatting consistent
	out, err := s.marshalDiagramYAML(&diagram)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}
	if err := os.WriteFile(filePath, out, 0o644); err != nil {
		return fmt.Errorf("failed to write YAML: %w", err)
	}
	return nil
}

// marshalDiagramYAML marshals the diagram to YAML and normalizes key styles for consistency.
// Historically we quoted keys like 'x' and 'y' to avoid YAML 1.1 plain-scalar ambiguity.
// We now prefer plain (unquoted) keys and explicitly tag them as strings to avoid misresolution.
func (s *DiagramService) marshalDiagramYAML(diagram *models.FlowDiagram) ([]byte, error) {
	// First marshal to bytes, then load into a yaml.Node tree to adjust styles
	raw, err := yaml.Marshal(diagram)
	if err != nil {
		return nil, err
	}
	var root yaml.Node
	if err := yaml.Unmarshal(raw, &root); err != nil {
		return nil, err
	}
	// Walk and normalize key styles
	normalizeMapKeyStyles(&root)
	// Encode with a stable indent
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(&root); err != nil {
		return nil, err
	}
	_ = enc.Close()
	return buf.Bytes(), nil
}

// normalizeMapKeyStyles walks the YAML AST and enforces preferred styles for certain keys.
// For keys named 'x' and 'y', we force them to be plain scalars (no quotes) and explicitly
// tag them as strings (!!str) to avoid any YAML 1.1 ambiguity while keeping a clean style.
func normalizeMapKeyStyles(n *yaml.Node) {
	if n == nil {
		return
	}
	switch n.Kind {
	case yaml.DocumentNode:
		for _, c := range n.Content {
			normalizeMapKeyStyles(c)
		}
	case yaml.SequenceNode:
		for _, c := range n.Content {
			normalizeMapKeyStyles(c)
		}
	case yaml.MappingNode:
		// Content is [k1, v1, k2, v2, ...]
		for i := 0; i+1 < len(n.Content); i += 2 {
			k := n.Content[i]
			v := n.Content[i+1]
			if k != nil && k.Kind == yaml.ScalarNode {
				if k.Value == "x" || k.Value == "y" {
					// Prefer plain style keys; ensure string tag for safety
					k.Tag = "!!str"
					k.Style = 0 // PlainStyle
				}
			}
			normalizeMapKeyStyles(v)
		}
	case yaml.ScalarNode, yaml.AliasNode:
		// nothing
	}
}
