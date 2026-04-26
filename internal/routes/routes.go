package routes

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"otas/internal/account"
	"otas/internal/transaction"
	"otas/internal/user"
	"otas/pkg/jwt"
)

// AuthMiddleware extracts the user ID from the JWT token and sets it in the Gin context.
// Temporarily handles empty auth headers for local development if needed, but enforces token validation.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// For local dev without auth fully wired, we could default to user 1,
			// but we will enforce security. You must pass a token.
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		claims, err := jwt.Validate(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token. " + err.Error()})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}

func Register(r *gin.Engine, db *sql.DB) {
	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OTAS API - Group System Standalone",
		})
	})

	// Initialize Repositories
	accountRepo := account.NewAccountRepository(db)
	userRepo := user.NewUserRepository(db)
	transactionRepo := transaction.NewTransactionRepository(db)

	// Initialize Services
	transactionService := transaction.NewTransactionService(transactionRepo, accountRepo, userRepo)

	// Initialize Handlers
	transactionHandler := transaction.NewTransactionHandler(transactionService)

	// Protected routes group
	protected := r.Group("/")
	protected.Use(AuthMiddleware())
	{
		// Transactions - triggers saving engine built by Dev 1
		protected.POST("/transactions", transactionHandler.ProcessTransaction)
	}
}
