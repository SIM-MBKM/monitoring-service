package controller

import (
	"fmt"
	"log"
	"monitoring-service/dto"
	"monitoring-service/helper"
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

func (c *ReportController) Approval(ctx *gin.Context) {
	id := ctx.Param("id")

	var reportApprovalRequest dto.ReportApprovalRequest
	if err := ctx.ShouldBindJSON(&reportApprovalRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	// If ID is provided in the URL and IDs array is empty, use the ID from the URL
	if id != "" && len(reportApprovalRequest.IDs) == 0 {
		reportApprovalRequest.IDs = []string{id}
	} else if id != "" {
		// If both ID in URL and IDs array are provided, ensure the ID from URL is included
		idFound := false
		for _, existingID := range reportApprovalRequest.IDs {
			if existingID == id {
				idFound = true
				break
			}
		}

		if !idFound {
			reportApprovalRequest.IDs = append(reportApprovalRequest.IDs, id)
		}
	}

	// Make sure we have at least one ID
	if len(reportApprovalRequest.IDs) == 0 {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "At least one report ID is required",
		})
		return
	}

	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Token is required",
		})
		return
	}

	err := c.reportService.Approval(ctx, token, reportApprovalRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	successMessage := "Report approved successfully"
	if len(reportApprovalRequest.IDs) > 1 {
		successMessage = "Reports approved successfully"
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Message: successMessage,
	})
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
		Status:  dto.STATUS_SUCCESS,
		Data:    reports,
		Message: "Reports fetched successfully",
	})
}

// Create handles POST /api/v1/reports
func (c *ReportController) Create(ctx *gin.Context) {
	var reportRequest dto.ReportRequest
	if ctx.Request.ContentLength > helper.MaxFileSize+helper.MaxContentLength {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Request too large",
		})
		return
	}

	if err := ctx.ShouldBind(&reportRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	if err := helper.ValidateReportRequest(reportRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	reportRequest.Title = helper.SanitizeString(reportRequest.Title)
	reportRequest.Content = helper.SanitizeString(reportRequest.Content)

	file, err := ctx.FormFile("file")
	// if file is not required, set file to nil
	if err != nil && err != http.ErrMissingFile {
		log.Printf("File upload error: %v", err)
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "File upload error",
		})
		return
	}

	if file != nil {
		if err := helper.ValidateFileUpload(file); err != nil {
			ctx.JSON(http.StatusBadRequest, dto.Response{
				Status:  dto.STATUS_ERROR,
				Message: fmt.Sprintf("File validation failed: %s", err.Error()),
			})
			return
		}

		// 🔒 SECURITY FIX 6: Validate MIME type by reading file content
		fileContent, err := file.Open()
		if err != nil {
			log.Printf("Error opening file: %v", err)
			ctx.JSON(http.StatusBadRequest, dto.Response{
				Status:  dto.STATUS_ERROR,
				Message: "Unable to process file",
			})
			return
		}
		defer fileContent.Close()

		if err := helper.ValidateMimeType(fileContent); err != nil {
			ctx.JSON(http.StatusBadRequest, dto.Response{
				Status:  dto.STATUS_ERROR,
				Message: err.Error(),
			})
			return
		}
	}

	if file == nil {
		file = nil
		if reportRequest.Content == "" {
			ctx.JSON(http.StatusBadRequest, dto.Response{
				Status:  dto.STATUS_ERROR,
				Message: "Content is required",
			})
			return
		}
	}

	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Token is required",
		})
		return
	}

	report, err := c.reportService.Create(ctx, reportRequest, file, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Data:    report,
		Message: "Report created successfully",
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
		Status:  dto.STATUS_SUCCESS,
		Message: "Report updated successfully",
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

	if !helper.ValidateUUID(id) {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Invalid ID format",
		})
		return
	}

	// 🔒 SECURITY FIX 2: Sanitize ID input (additional protection)
	sanitizedId := helper.SanitizeString(id)
	if sanitizedId != id {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Invalid characters in ID",
		})
		return
	}

	// get token from header
	token := ctx.GetHeader("Authorization")

	if !helper.IsValidTokenFormat(token) {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Invalid authorization format",
		})
		return
	}

	if token == "" {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Token is required",
		})
		return
	}

	report, err := c.reportService.FindByID(ctx, id, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Data:    report,
		Message: "Report fetched successfully",
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
		Status:  dto.STATUS_SUCCESS,
		Message: "Report deleted successfully",
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
		Status:  dto.STATUS_SUCCESS,
		Data:    reports,
		Message: "Reports fetched successfully",
	})
}
