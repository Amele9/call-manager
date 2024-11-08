package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/Amele9/call-manager/internal/models"
)

// AddCall adds a call to the database
func (s *GinServer) AddCall(c *gin.Context) {
	var call models.CallInfo

	if err := c.ShouldBindJSON(&call); err != nil {
		var validationError validator.ValidationErrors
		if errors.As(err, &validationError) {
			for _, fieldError := range validationError {
				switch fieldError.Field() {
				case "ClientName", "Description":
					if fieldError.Tag() == "required" {
						c.JSON(http.StatusBadRequest, gin.H{"error": fieldError.Error()})

						return
					}
				case "PhoneNumber":
					if fieldError.Tag() == "required" {
						c.JSON(http.StatusBadRequest, gin.H{"error": fieldError.Error()})

						return
					}

					if fieldError.Tag() == "e164" {
						c.JSON(http.StatusBadRequest, gin.H{"error": fieldError.Error()})

						return
					}
				}
			}
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	ID, err := s.Database.CreateCall(&call)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": ID})
}

// GetCalls returns all calls from the database
func (s *GinServer) GetCalls(c *gin.Context) {
	calls, err := s.Database.GetCalls()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"calls": calls})
}

// GetCallInfo returns call information from the database
func (s *GinServer) GetCallInfo(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	callInfo, err := s.Database.GetCallInfo(ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"callInfo": callInfo})
}

// UpdateCallStatus updates a call status in the database
func (s *GinServer) UpdateCallStatus(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	err = s.Database.UpdateCallStatus(ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Call updated"})
}

// DeleteCall deletes a call from the database
func (s *GinServer) DeleteCall(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	err = s.Database.DeleteCall(ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Call deleted"})
}
