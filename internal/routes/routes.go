package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, db *sql.DB) {
	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OTAS API - Group System Standalone",
		})
	})

}
