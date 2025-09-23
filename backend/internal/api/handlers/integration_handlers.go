package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetJiraProjects returns available Jira projects
func GetJiraProjects(c *gin.Context) {
	// TODO: Implement Jira integration
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":   "Jira integration not yet implemented",
		"message": "This endpoint will return available Jira projects",
	})
}

// GetJiraIssue returns a specific Jira issue
func GetJiraIssue(c *gin.Context) {
	issueKey := c.Param("key")
	if issueKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Issue key is required",
		})
		return
	}

	// TODO: Implement Jira integration
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":   "Jira integration not yet implemented",
		"message": "This endpoint will return Jira issue details for: " + issueKey,
	})
}

// CreateJiraIssue creates a new Jira issue
func CreateJiraIssue(c *gin.Context) {
	var issueRequest struct {
		Summary     string `json:"summary" binding:"required"`
		Description string `json:"description"`
		Project     string `json:"project" binding:"required"`
		IssueType   string `json:"issueType" binding:"required"`
		Priority    string `json:"priority"`
	}

	if err := c.ShouldBindJSON(&issueRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid issue request",
			"details": err.Error(),
		})
		return
	}

	// TODO: Implement Jira integration
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":   "Jira integration not yet implemented",
		"message": "This endpoint will create a Jira issue",
		"request": issueRequest,
	})
}
