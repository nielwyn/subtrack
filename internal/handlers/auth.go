package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nielwyn/inventory-system/internal/models"
	"github.com/nielwyn/inventory-system/internal/service"
	apierrors "github.com/nielwyn/inventory-system/pkg/errors"
	"github.com/nielwyn/inventory-system/pkg/logger"
	"github.com/nielwyn/inventory-system/pkg/response"
	"github.com/nielwyn/inventory-system/pkg/validator"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apierrors.CodeInvalidInput, validator.FormatValidationError(err))
		return
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		if errors.Is(err, service.ErrUsernameExists) {
			response.Error(c, http.StatusConflict, apierrors.CodeUsernameExists, err.Error())
			return
		}
		if errors.Is(err, service.ErrEmailExists) {
			response.Error(c, http.StatusConflict, apierrors.CodeEmailExists, err.Error())
			return
		}
		logger.Error("Registration failed", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, apierrors.CodeInternal, "Registration failed")
		return
	}

	response.Success(c, http.StatusCreated, "User registered successfully", gin.H{
		"user": user,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apierrors.CodeInvalidInput, validator.FormatValidationError(err))
		return
	}

	loginResponse, err := h.authService.Login(&req)
	if err != nil {
		logger.Error("Login failed", zap.Error(err))
		response.Error(c, http.StatusUnauthorized, apierrors.CodeUnauthorized, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Login successful", loginResponse)
}
