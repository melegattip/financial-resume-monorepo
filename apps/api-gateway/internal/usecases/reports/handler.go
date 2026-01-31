package reports

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/melegattip/financial-resume-engine/internal/core/logs"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/logger"
	"gorm.io/gorm"
)

type ReportHandler struct {
	db *gorm.DB
}

func NewReportHandler(db *gorm.DB) *ReportHandler {
	return &ReportHandler{db: db}
}

func (h *ReportHandler) HandleGenerateReport(c *gin.Context) {
	// Validar header X-Caller-ID (obtenemos del contexto seteo por middleware)
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "x-caller-id header is required"})
		return
	}

	var request GenerateReportRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", request.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
		return
	}

	// Ajustar endDate para incluir todo el día (hasta 23:59:59.999999999)
	endDate = endDate.Add(24*time.Hour - time.Nanosecond)

	// Validar que la fecha final no sea anterior a la inicial
	if endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date cannot be before start date"})
		return
	}

	service := NewGenerateFinancialReport(NewReportRepository(h.db))
	report, err := service.Execute(startDate, endDate, userID)
	if err != nil {
		logger.Error(c.Request.Context(), err, logs.ErrorGeneratingReport.GetMessage(), logs.Tags{
			"user_id":    userID,
			"start_date": request.StartDate,
			"end_date":   request.EndDate,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating report"})
		return
	}

	c.JSON(http.StatusOK, report)
}
