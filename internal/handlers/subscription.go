package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nielwyn/inventory-system/internal/models"
	"github.com/nielwyn/inventory-system/internal/service"
	apierrors "github.com/nielwyn/inventory-system/pkg/errors"
	"github.com/nielwyn/inventory-system/pkg/logger"
	"github.com/nielwyn/inventory-system/pkg/response"
	"github.com/nielwyn/inventory-system/pkg/validator"
	"go.uber.org/zap"
)

type SubscriptionHandler struct {
	subscriptionService service.SubscriptionService
}

func NewSubscriptionHandler(subscriptionService service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{subscriptionService: subscriptionService}
}

func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var req models.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apierrors.CodeInvalidInput, validator.FormatValidationError(err))
		return
	}

	sub, err := h.subscriptionService.CreateSubscription(&req)
	if err != nil {
		logger.Error("Failed to create subscription", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, apierrors.CodeInternal, "Failed to create subscription")
		return
	}

	response.Success(c, http.StatusCreated, "Subscription created successfully", sub)
}

func (h *SubscriptionHandler) GetAllSubscriptions(c *gin.Context) {
	var query models.SubscriptionQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, apierrors.CodeInvalidInput, validator.FormatValidationError(err))
		return
	}

	result, err := h.subscriptionService.GetAllSubscriptions(query)
	if err != nil {
		logger.Error("Failed to retrieve subscriptions", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, apierrors.CodeInternal, "Failed to retrieve subscriptions")
		return
	}

	response.Success(c, http.StatusOK, "Subscriptions retrieved successfully", result)
}

func (h *SubscriptionHandler) GetSubscriptionByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apierrors.CodeInvalidInput, "Invalid subscription ID")
		return
	}

	sub, err := h.subscriptionService.GetSubscriptionByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrSubscriptionNotFound) {
			response.Error(c, http.StatusNotFound, apierrors.CodeNotFound, err.Error())
			return
		}
		logger.Error("Failed to retrieve subscription", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, apierrors.CodeInternal, "Failed to retrieve subscription")
		return
	}

	response.Success(c, http.StatusOK, "Subscription retrieved successfully", sub)
}

func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apierrors.CodeInvalidInput, "Invalid subscription ID")
		return
	}

	var req models.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, apierrors.CodeInvalidInput, validator.FormatValidationError(err))
		return
	}

	sub, err := h.subscriptionService.UpdateSubscription(uint(id), &req)
	if err != nil {
		if errors.Is(err, service.ErrSubscriptionNotFound) {
			response.Error(c, http.StatusNotFound, apierrors.CodeNotFound, err.Error())
			return
		}
		logger.Error("Failed to update subscription", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, apierrors.CodeInternal, "Failed to update subscription")
		return
	}

	response.Success(c, http.StatusOK, "Subscription updated successfully", sub)
}

func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, apierrors.CodeInvalidInput, "Invalid subscription ID")
		return
	}

	if err := h.subscriptionService.DeleteSubscription(uint(id)); err != nil {
		if errors.Is(err, service.ErrSubscriptionNotFound) {
			response.Error(c, http.StatusNotFound, apierrors.CodeNotFound, err.Error())
			return
		}
		logger.Error("Failed to delete subscription", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, apierrors.CodeInternal, "Failed to delete subscription")
		return
	}

	response.Success(c, http.StatusOK, "Subscription deleted successfully", nil)
}
