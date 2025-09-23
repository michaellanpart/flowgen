package api

import (
	"github.com/gin-gonic/gin"
	"github.com/michaellanpart/flowgen/backend/internal/api/handlers"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		// Diagram routes
		diagrams := api.Group("/diagrams")
		{
			diagrams.GET("", handlers.ListDiagrams)
			diagrams.POST("", handlers.CreateDiagram)
			diagrams.GET("/:id", handlers.GetDiagram)
			diagrams.PUT("/:id", handlers.UpdateDiagram)
			diagrams.DELETE("/:id", handlers.DeleteDiagram)
			diagrams.POST("/:id/validate", handlers.ValidateDiagram)
			// Raw YAML access for Git-friendly workflows
			diagrams.GET("/:id/yaml", handlers.GetDiagramYAML)
			diagrams.PUT("/:id/yaml", handlers.UpdateDiagramYAML)
		}

		// Hierarchy routes for drill-down functionality
		hierarchy := api.Group("/hierarchy")
		{
			hierarchy.GET("/:id/children", handlers.GetChildDiagrams)
			hierarchy.GET("/:id/parent", handlers.GetParentDiagram)
			hierarchy.POST("/:id/link", handlers.LinkDiagrams)
		}

		// Integration routes
		integrations := api.Group("/integrations")
		{
			jira := integrations.Group("/jira")
			{
				jira.GET("/projects", handlers.GetJiraProjects)
				jira.GET("/issues/:key", handlers.GetJiraIssue)
				jira.POST("/issues", handlers.CreateJiraIssue)
			}
		}

		// Search and analytics
		search := api.Group("/search")
		{
			search.GET("/diagrams", handlers.SearchDiagrams)
			search.GET("/nodes", handlers.SearchNodes)
		}
	}
}
