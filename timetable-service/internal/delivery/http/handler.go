package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/hospital-system-api/timetable-service/internal/domain"
	"github.com/yourusername/hospital-system-api/timetable-service/internal/service"
	"github.com/yourusername/hospital-system-api/timetable-service/pkg/auth"
)

type Handler struct {
	timetableService service.TimetableService
	authClient       auth.Client
}

func NewHandler(timetableService service.TimetableService, authClient auth.Client) *Handler {
	return &Handler{
		timetableService: timetableService,
		authClient:       authClient,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	api := router.Group("/api/v1")
	{
		timetables := api.Group("/timetables")
		{
			timetables.GET("", h.listTimetables)
			timetables.GET("/:id", h.getTimetable)
			timetables.GET("/:id/appointments", h.getTimetableAppointments)
			timetables.POST("", h.authMiddleware(), h.adminMiddleware(), h.createTimetable)
			timetables.PUT("/:id", h.authMiddleware(), h.adminMiddleware(), h.updateTimetable)
			timetables.DELETE("/:id", h.authMiddleware(), h.adminMiddleware(), h.deleteTimetable)
		}

		appointments := api.Group("/appointments")
		{
			appointments.POST("/:timetableID", h.authMiddleware(), h.createAppointment)
			appointments.DELETE("/:id", h.authMiddleware(), h.deleteAppointment)
		}
	}
}

func (h *Handler) listTimetables(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	timetables, err := h.timetableService.ListTimetables(offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timetables)
}

func (h *Handler) getTimetable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	timetable, err := h.timetableService.GetTimetable(uint(id))
	if err != nil {
		if err == service.ErrTimetableNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timetable)
}

func (h *Handler) getTimetableAppointments(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	appointments, err := h.timetableService.GetAppointments(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, appointments)
}

func (h *Handler) createTimetable(c *gin.Context) {
	var req domain.CreateTimetableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timetable := &domain.Timetable{
		HospitalID: req.HospitalID,
		DoctorID:   req.DoctorID,
		From:       req.From,
		To:         req.To,
		Room:       req.Room,
	}

	if err := h.timetableService.CreateTimetable(timetable); err != nil {
		if err == service.ErrInvalidTimeRange {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, timetable)
}

func (h *Handler) updateTimetable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req domain.UpdateTimetableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timetable := &domain.Timetable{
		ID:         uint(id),
		HospitalID: req.HospitalID,
		DoctorID:   req.DoctorID,
		From:       req.From,
		To:         req.To,
		Room:       req.Room,
	}

	if err := h.timetableService.UpdateTimetable(timetable); err != nil {
		if err == service.ErrTimetableNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err == service.ErrInvalidTimeRange {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timetable)
}

func (h *Handler) deleteTimetable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.timetableService.DeleteTimetable(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) createAppointment(c *gin.Context) {
	timetableID, err := strconv.ParseUint(c.Param("timetableID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid timetable id"})
		return
	}

	var req domain.CreateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	if err := h.timetableService.CreateAppointment(uint(timetableID), userID, req.Time); err != nil {
		if err == service.ErrTimetableNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err == service.ErrInvalidTimeRange || err == service.ErrTimeSlotTaken {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) deleteAppointment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.timetableService.DeleteAppointment(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization token required"})
			c.Abort()
			return
		}

		if err := h.authClient.ValidateToken(token); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// TODO: Extract user ID from token and set it in context
		c.Set("user_id", uint(1)) // Temporary hardcoded user ID
		c.Next()
	}
}

func (h *Handler) adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Check if user has admin role
		c.Next()
	}
}
