package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sergeimurashev/hospital-system-api/hospital-service/internal/domain"
	"github.com/sergeimurashev/hospital-system-api/hospital-service/internal/service"
)

type Handler struct {
	hospitalService service.HospitalService
}

func NewHandler(hospitalService service.HospitalService) *Handler {
	return &Handler{
		hospitalService: hospitalService,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	api := router.Group("/api")
	{
		hospitals := api.Group("/Hospitals")
		{
			hospitals.GET("", h.authMiddleware(), h.listHospitals)
			hospitals.GET("/:id", h.authMiddleware(), h.getHospital)
			hospitals.GET("/:id/Rooms", h.authMiddleware(), h.getHospitalRooms)
			hospitals.POST("", h.adminMiddleware(), h.createHospital)
			hospitals.PUT("/:id", h.adminMiddleware(), h.updateHospital)
			hospitals.DELETE("/:id", h.adminMiddleware(), h.deleteHospital)
		}
	}
}

func (h *Handler) listHospitals(c *gin.Context) {
	from, _ := strconv.Atoi(c.DefaultQuery("from", "0"))
	count, _ := strconv.Atoi(c.DefaultQuery("count", "10"))

	hospitals, err := h.hospitalService.ListHospitals(from, count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hospitals)
}

func (h *Handler) getHospital(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	hospital, err := h.hospitalService.GetHospitalByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hospital)
}

func (h *Handler) getHospitalRooms(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	rooms, err := h.hospitalService.GetHospitalRooms(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func (h *Handler) createHospital(c *gin.Context) {
	var req domain.CreateHospitalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.hospitalService.CreateHospital(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) updateHospital(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req domain.UpdateHospitalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.hospitalService.UpdateHospital(uint(id), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) deleteHospital(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.hospitalService.DeleteHospital(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		if err := h.hospitalService.ValidateToken(token); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (h *Handler) adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.authMiddleware()(c)
		if c.IsAborted() {
			return
		}

		// TODO: Implement admin role check
		c.Next()
	}
}
