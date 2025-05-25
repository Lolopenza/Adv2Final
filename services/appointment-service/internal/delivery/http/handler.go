package http

import (
	"net/http"
	"time"

	"appointment-service/internal/domain"
	"appointment-service/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	useCase usecase.UseCase
}

func NewHandler(useCase usecase.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		api.POST("/slots", h.CreateSlots)
		api.GET("/slots", h.GetAvailableSlots)
		api.POST("/appointments", h.BookAppointment)
	}
}

func (h *Handler) CreateSlots(c *gin.Context) {
	var req domain.CreateSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.useCase.CreateSlots(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Slots created successfully"})
}

func (h *Handler) GetAvailableSlots(c *gin.Context) {
	businessID, err := uuid.Parse(c.Query("business_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid business_id"})
		return
	}

	date, err := time.Parse("2006-01-02", c.Query("date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	slots, err := h.useCase.GetAvailableSlots(c.Request.Context(), businessID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, slots)
}

func (h *Handler) BookAppointment(c *gin.Context) {
	var req domain.BookAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appointment, err := h.useCase.BookAppointment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, appointment)
}
