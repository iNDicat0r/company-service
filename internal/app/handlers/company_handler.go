package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iNDicat0r/company/common"
	"github.com/iNDicat0r/company/internal/app/services"
)

// CompanyHandler is responsible for handling routes for company resources.
type CompanyHandler struct {
	CompanyService services.CompanyGetCreateUpdateDeleter
}

// NewCompanyHandler creates a new company handler.
func NewCompanyHandler(companyService services.CompanyGetCreateUpdateDeleter) (*CompanyHandler, error) {
	if companyService == nil {
		return nil, errors.New("company service is nil")
	}
	return &CompanyHandler{CompanyService: companyService}, nil
}

// HandleGetCompany get a company handler.
func (h *CompanyHandler) HandleGetCompany(c *gin.Context) {
	id, err := uuid.Parse(c.Param("companyID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comp, err := h.CompanyService.Get(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comp)
}

type createCompanyRequestPayload struct {
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	EmployeesAmount int         `json:"employees_amount"`
	Registered      bool        `json:"registered"`
	Type            common.Type `json:"type"`
}

// HandleCreateCompany handles creating a company.
func (h *CompanyHandler) HandleCreateCompany(c *gin.Context) {
	var reqBody createCompanyRequestPayload
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload := services.CreateUpdateCompanyPayload{
		Name:            reqBody.Name,
		Description:     reqBody.Description,
		EmployeesAmount: reqBody.EmployeesAmount,
		Registered:      reqBody.Registered,
		Type:            reqBody.Type,
	}

	comp, err := h.CompanyService.Create(c, userID, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, comp)
}

// HandleUpdateCompany handles updating a company.
func (h *CompanyHandler) HandleUpdateCompany(c *gin.Context) {
	id, err := uuid.Parse(c.Param("companyID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var reqBody createCompanyRequestPayload
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload := services.CreateUpdateCompanyPayload{
		Name:            reqBody.Name,
		Description:     reqBody.Description,
		EmployeesAmount: reqBody.EmployeesAmount,
		Registered:      reqBody.Registered,
		Type:            reqBody.Type,
	}

	comp, err := h.CompanyService.Update(c, id, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comp)
}

// HandleDeleteCompany handles deleting a company.
func (h *CompanyHandler) HandleDeleteCompany(c *gin.Context) {
	id, err := uuid.Parse(c.Param("companyID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.CompanyService.Delete(c, userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}
