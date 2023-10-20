package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/iNDicat0r/company/internal/app/utils"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	jwtSigner := "privateKey-secret"
	router := gin.New()
	router.Use(AuthMiddleware(jwtSigner))

	jwtToken, err := utils.GenerateJWT(jwtSigner, "12")
	assert.NoError(t, err)

	router.GET("/v1/auth/introspect", func(c *gin.Context) {
		userID := c.GetString("userID")
		c.JSON(http.StatusOK, gin.H{"message": "Authenticated user ID: " + userID})
	})

	req := httptest.NewRequest("GET", "/v1/auth/introspect", nil)
	req.Header.Set("Authorization", jwtToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"message":"Authenticated user ID: 12"}`, w.Body.String())
}
