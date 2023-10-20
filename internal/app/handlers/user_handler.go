package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iNDicat0r/company/internal/app/services"
)

// authenticateRequestBody represents the authentication payload.
type authenticateRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// athenticationResponse represnts the authentication jwt token.
type athenticationResponse struct {
	Token string `json:"token"`
}

// UserHandler is responsible for handling routes for user resources.
type UserHandler struct {
	userService services.User
}

// NewUserHandler creates a new user routes handler.
func NewUserHandler(userService services.User) (*UserHandler, error) {
	if userService == nil {
		return nil, errors.New("user service is nil")
	}

	return &UserHandler{userService: userService}, nil
}

// HandleAuthenticate handles the authentication of a user.
func (uh *UserHandler) HandleAuthenticate(c *gin.Context) {
	var reqBody authenticateRequestBody

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwt, err := uh.userService.Authenticate(c, reqBody.Username, reqBody.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := athenticationResponse{
		Token: jwt,
	}
	c.JSON(http.StatusOK, resp)
}

// HandleIntrospect handles the user token introspection.
func (uh *UserHandler) HandleIntrospect(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := uh.userService.GetUser(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
