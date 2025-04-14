package controller

import (
	"monitoring-service/dto"
	"monitoring-service/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	reportService service.ReportService
}

func NewReportController(reportService service.ReportService) *ReportController {
	return &ReportController{
		reportService: reportService,
	}
}

// Index handles GET /api/v1/reports
func (c *ReportController) Index(ctx *gin.Context) {
	reports, err := c.reportService.Index(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status: dto.STATUS_SUCCESS,
		Data:   reports,
	})
}

// Create handles POST /api/v1/reports
func (c *ReportController) Create(ctx *gin.Context) {
	var reportRequest dto.ReportRequest
	if err := ctx.ShouldBindJSON(&reportRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "File is required",
		})
		return
	}

	report, err := c.reportService.Create(ctx, reportRequest, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Status: dto.STATUS_SUCCESS,
		Data:   report,
	})
}

// Update handles PUT /api/v1/reports/:id
func (c *ReportController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "ID is required",
		})
		return
	}

	var reportRequest dto.ReportRequest
	if err := ctx.ShouldBindJSON(&reportRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	err := c.reportService.Update(ctx, id, reportRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status: dto.STATUS_SUCCESS,
	})
}

// Show handles GET /api/v1/reports/:id
func (c *ReportController) Show(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "ID is required",
		})
		return
	}

	report, err := c.reportService.FindByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status: dto.STATUS_SUCCESS,
		Data:   report,
	})
}

// Destroy handles DELETE /api/v1/reports/:id
func (c *ReportController) Destroy(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "ID is required",
		})
		return
	}

	err := c.reportService.Destroy(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status: dto.STATUS_SUCCESS,
	})
}

// FindByReportScheduleID handles GET /api/v1/report-schedules/:id/reports
func (c *ReportController) FindByReportScheduleID(ctx *gin.Context) {
	reportScheduleID := ctx.Param("id")
	if reportScheduleID == "" {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Report Schedule ID is required",
		})
		return
	}

	reports, err := c.reportService.FindByReportScheduleID(ctx, reportScheduleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status: dto.STATUS_SUCCESS,
		Data:   reports,
	})
}
