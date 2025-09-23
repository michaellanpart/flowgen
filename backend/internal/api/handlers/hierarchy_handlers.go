package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/michaellanpart/flowgen/backend/internal/services"
)

// GetChildDiagrams returns child diagrams for hierarchical drill-down
func GetChildDiagrams(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Diagram ID is required",
		})
		return
	}

	hierarchyService := services.NewHierarchyService()

	children, err := hierarchyService.GetChildren(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get child diagrams",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"parent":   id,
		"children": children,
		"count":    len(children),
	})
}

// GetParentDiagram returns the parent diagram
func GetParentDiagram(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Diagram ID is required",
		})
		return
	}

	hierarchyService := services.NewHierarchyService()

	parent, err := hierarchyService.GetParent(id)
	if err != nil {
		if err == services.ErrDiagramNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Parent diagram not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get parent diagram",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"child":  id,
		"parent": parent,
	})
}

// LinkDiagrams creates a hierarchical relationship between diagrams
func LinkDiagrams(c *gin.Context) {
	parentID := c.Param("id")
	if parentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Parent diagram ID is required",
		})
		return
	}

	var linkRequest struct {
		ChildID string `json:"childId" binding:"required"`
		NodeID  string `json:"nodeId"` // Optional: specific node for drill-down
	}

	if err := c.ShouldBindJSON(&linkRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid link request",
			"details": err.Error(),
		})
		return
	}

	hierarchyService := services.NewHierarchyService()

	err := hierarchyService.LinkDiagrams(parentID, linkRequest.ChildID, linkRequest.NodeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to link diagrams",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Diagrams linked successfully",
		"parent":  parentID,
		"child":   linkRequest.ChildID,
		"node":    linkRequest.NodeID,
	})
}
