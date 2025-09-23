package handlers

import (
	"net/http"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/michaellanpart/flowgen/backend/internal/models"
	"github.com/michaellanpart/flowgen/backend/internal/services"
)

// ListDiagrams returns all available diagrams
func ListDiagrams(c *gin.Context) {
	diagramService := services.NewDiagramService()

	diagrams, err := diagramService.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list diagrams",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"diagrams": diagrams,
		"count":    len(diagrams),
	})
}

// GetDiagram returns a specific diagram by ID
func GetDiagram(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Diagram ID is required",
		})
		return
	}

	diagramService := services.NewDiagramService()

	diagram, err := diagramService.GetByID(id)
	if err != nil {
		if err == services.ErrDiagramNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Diagram not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get diagram",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, diagram)
}

// CreateDiagram creates a new diagram
func CreateDiagram(c *gin.Context) {
	var diagram models.FlowDiagram

	if err := c.ShouldBindJSON(&diagram); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid diagram data",
			"details": err.Error(),
		})
		return
	}

	diagramService := services.NewDiagramService()

	createdDiagram, err := diagramService.Create(&diagram)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create diagram",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, createdDiagram)
}

// UpdateDiagram updates an existing diagram
func UpdateDiagram(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Diagram ID is required",
		})
		return
	}

	var diagram models.FlowDiagram

	if err := c.ShouldBindJSON(&diagram); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid diagram data",
			"details": err.Error(),
		})
		return
	}

	// Ensure the ID matches
	diagram.ID = id

	diagramService := services.NewDiagramService()

	updatedDiagram, err := diagramService.Update(&diagram)
	if err != nil {
		if err == services.ErrDiagramNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Diagram not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update diagram",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updatedDiagram)
}

// DeleteDiagram deletes a diagram
func DeleteDiagram(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Diagram ID is required",
		})
		return
	}

	diagramService := services.NewDiagramService()

	err := diagramService.Delete(id)
	if err != nil {
		if err == services.ErrDiagramNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Diagram not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete diagram",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Diagram deleted successfully",
	})
}

// ValidateDiagram validates a diagram against the schema
func ValidateDiagram(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Diagram ID is required",
		})
		return
	}

	diagramService := services.NewDiagramService()

	diagram, err := diagramService.GetByID(id)
	if err != nil {
		if err == services.ErrDiagramNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Diagram not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get diagram",
			"details": err.Error(),
		})
		return
	}

	validationResult, err := diagramService.Validate(diagram)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to validate diagram",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, validationResult)
}

// GetDiagramYAML returns the raw YAML of a diagram by ID
func GetDiagramYAML(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.String(http.StatusBadRequest, "diagram id is required")
		return
	}

	svc := services.NewDiagramService()
	yamlContent, err := svc.LoadYAMLByID(id)
	if err != nil {
		if err == services.ErrDiagramNotFound {
			c.String(http.StatusNotFound, "diagram not found")
			return
		}
		c.String(http.StatusInternalServerError, "failed to load yaml: %v", err)
		return
	}

	c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(yamlContent))
}

// UpdateDiagramYAML updates the raw YAML for a diagram ID
func UpdateDiagramYAML(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.String(http.StatusBadRequest, "diagram id is required")
		return
	}

	// Read raw text body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "failed to read request body: %v", err)
		return
	}

	yamlText := string(body)
	svc := services.NewDiagramService()

	// Save and validate
	if err := svc.SaveYAMLByID(id, yamlText); err != nil {
		if err == services.ErrDiagramNotFound {
			c.String(http.StatusNotFound, "diagram not found")
			return
		}
		c.String(http.StatusBadRequest, "invalid yaml: %v", err)
		return
	}

	c.String(http.StatusOK, "ok")
}

// SearchDiagrams searches for diagrams based on query parameters
func SearchDiagrams(c *gin.Context) {
	query := c.Query("q")
	tags := c.QueryArray("tags")

	diagramService := services.NewDiagramService()

	results, err := diagramService.Search(query, tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to search diagrams",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results": results,
		"count":   len(results),
		"query":   query,
		"tags":    tags,
	})
}

// SearchNodes searches for nodes across all diagrams
func SearchNodes(c *gin.Context) {
	query := c.Query("q")
	nodeType := c.Query("type")

	diagramService := services.NewDiagramService()

	results, err := diagramService.SearchNodes(query, nodeType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to search nodes",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results": results,
		"count":   len(results),
		"query":   query,
		"type":    nodeType,
	})
}
