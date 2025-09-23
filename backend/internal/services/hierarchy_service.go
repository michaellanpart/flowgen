package services

import (
	"fmt"
	"github.com/michaellanpart/flowgen/backend/internal/models"
)

// HierarchyService handles diagram hierarchy operations
type HierarchyService struct {
	diagramService *DiagramService
}

// NewHierarchyService creates a new hierarchy service
func NewHierarchyService() *HierarchyService {
	return &HierarchyService{
		diagramService: NewDiagramService(),
	}
}

// GetChildren returns child diagrams for a given parent
func (s *HierarchyService) GetChildren(parentID string) ([]models.FlowDiagram, error) {
	parent, err := s.diagramService.GetByID(parentID)
	if err != nil {
		return nil, err
	}

	children := []models.FlowDiagram{}

	for _, childID := range parent.Children {
		child, err := s.diagramService.GetByID(childID)
		if err != nil {
			// Log error but continue with other children
			fmt.Printf("Error getting child diagram %s: %v\n", childID, err)
			continue
		}
		children = append(children, *child)
	}

	return children, nil
}

// GetParent returns the parent diagram for a given child
func (s *HierarchyService) GetParent(childID string) (*models.FlowDiagram, error) {
	child, err := s.diagramService.GetByID(childID)
	if err != nil {
		return nil, err
	}

	if child.Parent == nil {
		return nil, fmt.Errorf("diagram has no parent")
	}

	return s.diagramService.GetByID(*child.Parent)
}

// LinkDiagrams creates a hierarchical relationship between diagrams
func (s *HierarchyService) LinkDiagrams(parentID, childID, nodeID string) error {
	// Get parent diagram
	parent, err := s.diagramService.GetByID(parentID)
	if err != nil {
		return fmt.Errorf("failed to get parent diagram: %w", err)
	}

	// Get child diagram
	child, err := s.diagramService.GetByID(childID)
	if err != nil {
		return fmt.Errorf("failed to get child diagram: %w", err)
	}

	// Update parent to include child
	childExists := false
	for _, existingChildID := range parent.Children {
		if existingChildID == childID {
			childExists = true
			break
		}
	}

	if !childExists {
		parent.Children = append(parent.Children, childID)
	}

	// If nodeID is specified, update the node to include drill-down reference
	if nodeID != "" {
		nodeFound := false
		for i, node := range parent.Nodes {
			if node.ID == nodeID {
				parent.Nodes[i].DrillDown = &childID
				nodeFound = true
				break
			}
		}

		if !nodeFound {
			return fmt.Errorf("node %s not found in parent diagram", nodeID)
		}
	}

	// Update child to reference parent
	child.Parent = &parentID

	// Save both diagrams
	if _, err := s.diagramService.Update(parent); err != nil {
		return fmt.Errorf("failed to update parent diagram: %w", err)
	}

	if _, err := s.diagramService.Update(child); err != nil {
		return fmt.Errorf("failed to update child diagram: %w", err)
	}

	return nil
}

// UnlinkDiagrams removes a hierarchical relationship
func (s *HierarchyService) UnlinkDiagrams(parentID, childID string) error {
	// Get parent diagram
	parent, err := s.diagramService.GetByID(parentID)
	if err != nil {
		return fmt.Errorf("failed to get parent diagram: %w", err)
	}

	// Get child diagram
	child, err := s.diagramService.GetByID(childID)
	if err != nil {
		return fmt.Errorf("failed to get child diagram: %w", err)
	}

	// Remove child from parent's children list
	newChildren := []string{}
	for _, existingChildID := range parent.Children {
		if existingChildID != childID {
			newChildren = append(newChildren, existingChildID)
		}
	}
	parent.Children = newChildren

	// Remove drill-down references from nodes
	for i, node := range parent.Nodes {
		if node.DrillDown != nil && *node.DrillDown == childID {
			parent.Nodes[i].DrillDown = nil
		}
	}

	// Remove parent reference from child
	child.Parent = nil

	// Save both diagrams
	if _, err := s.diagramService.Update(parent); err != nil {
		return fmt.Errorf("failed to update parent diagram: %w", err)
	}

	if _, err := s.diagramService.Update(child); err != nil {
		return fmt.Errorf("failed to update child diagram: %w", err)
	}

	return nil
}

// GetHierarchyTree returns the complete hierarchy tree starting from a root diagram
func (s *HierarchyService) GetHierarchyTree(rootID string) (*HierarchyNode, error) {
	return s.buildHierarchyNode(rootID, make(map[string]bool))
}

// HierarchyNode represents a node in the hierarchy tree
type HierarchyNode struct {
	Diagram  models.FlowDiagram `json:"diagram"`
	Children []*HierarchyNode   `json:"children"`
}

func (s *HierarchyService) buildHierarchyNode(diagramID string, visited map[string]bool) (*HierarchyNode, error) {
	// Prevent infinite loops
	if visited[diagramID] {
		return nil, fmt.Errorf("circular reference detected in hierarchy: %s", diagramID)
	}
	visited[diagramID] = true

	diagram, err := s.diagramService.GetByID(diagramID)
	if err != nil {
		return nil, err
	}

	node := &HierarchyNode{
		Diagram:  *diagram,
		Children: []*HierarchyNode{},
	}

	// Recursively build child nodes
	for _, childID := range diagram.Children {
		childNode, err := s.buildHierarchyNode(childID, visited)
		if err != nil {
			// Log error but continue with other children
			fmt.Printf("Error building hierarchy for child %s: %v\n", childID, err)
			continue
		}
		node.Children = append(node.Children, childNode)
	}

	// Remove from visited to allow the same diagram in different branches
	delete(visited, diagramID)

	return node, nil
}
