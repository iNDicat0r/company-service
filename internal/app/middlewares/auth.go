package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iNDicat0r/company/internal/app/utils"
)

// AuthMiddleware is an authenticator middleware.
func AuthMiddleware(jwtPrivateKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		userID, err := utils.ParseJWT(jwtPrivateKey, authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			return
		}

		// set the logged in userID in the context
		// so it can be used later on
		c.Set("userID", userID)

		c.Next()
	}
}
