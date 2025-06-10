package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sergeimurashev/hospital-system-api/document-service/internal/domain"
	"github.com/sergeimurashev/hospital-system-api/document-service/internal/service"
	"github.com/sergeimurashev/hospital-system-api/document-service/pkg/auth"
)

type Handler struct {
	documentService service.DocumentService
	authClient      auth.Client
}

func NewHandler(documentService service.DocumentService, authClient auth.Client) *Handler {
	return &Handler{
		documentService: documentService,
		authClient:      authClient,
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	api := router.Group("/api/v1")
	{
		history := api.Group("/history")
		{
			history.GET("/account/:id", h.authMiddleware(), h.getPatientDocuments)
			history.GET("/:id", h.authMiddleware(), h.getDocument)
			history.POST("", h.authMiddleware(), h.createDocument)
			history.PUT("/:id", h.authMiddleware(), h.updateDocument)
		}

		search := api.Group("/search")
		{
			search.GET("", h.authMiddleware(), h.searchDocuments)
		}
	}
}

func (h *Handler) getPatientDocuments(c *gin.Context) {
	patientID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient id"})
		return
	}

	documents, err := h.documentService.GetPatientDocuments(uint(patientID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, documents)
}

func (h *Handler) getDocument(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	document, err := h.documentService.GetDocument(uint(id))
	if err != nil {
		if err == service.ErrDocumentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *Handler) createDocument(c *gin.Context) {
	var req domain.CreateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	document := &domain.Document{
		Date:       req.Date,
		PatientID:  req.PatientID,
		HospitalID: req.HospitalID,
		DoctorID:   req.DoctorID,
		Room:       req.Room,
		Data:       req.Data,
	}

	if err := h.documentService.CreateDocument(document); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, document)
}

func (h *Handler) updateDocument(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req domain.UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	document := &domain.Document{
		ID:         uint(id),
		Date:       req.Date,
		PatientID:  req.PatientID,
		HospitalID: req.HospitalID,
		DoctorID:   req.DoctorID,
		Room:       req.Room,
		Data:       req.Data,
	}

	if err := h.documentService.UpdateDocument(document); err != nil {
		if err == service.ErrDocumentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *Handler) searchDocuments(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search query is required"})
		return
	}

	documents, err := h.documentService.SearchDocuments(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, documents)
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

		c.Set("user_id", uint(1))
		c.Next()
	}
}
