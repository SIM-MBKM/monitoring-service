package controller

import (
	"log"
	"monitoring-service/dto"
	"monitoring-service/helper"
	"monitoring-service/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReportScheduleController struct {
	reportScheduleService service.ReportScheduleService
}

func NewReportScheduleController(reportScheduleService service.ReportScheduleService) *ReportScheduleController {
	return &ReportScheduleController{
		reportScheduleService: reportScheduleService,
	}
}

func (c *ReportScheduleController) FindByAdvisorEmail(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Token is required",
		})
		return
	}

	pagReq := helper.Pagination(ctx)

	var reportScheduleRequest dto.ReportScheduleAdvisorRequest
	if err := ctx.ShouldBindJSON(&reportScheduleRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	reportSchedules, metaData, err := c.reportScheduleService.FindByAdvisorEmail(ctx, token, pagReq, reportScheduleRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:             dto.STATUS_SUCCESS,
		Data:               reportSchedules,
		Message:            "Report schedule found successfully",
		PaginationResponse: &metaData,
	})
}

func (c *ReportScheduleController) FindByStudentID(ctx *gin.Context) {
	// get token from header
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Token is required",
		})
		return
	}

	reportSchedules, err := c.reportScheduleService.FindByUserNRPAndGroupByRegistrationID(ctx, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Data:    reportSchedules,
		Message: "Report schedule found successfully",
	})
}

// Index handles GET /api/v1/report-schedules
func (c *ReportScheduleController) Index(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Token is required",
		})
		return
	}

	reportSchedules, err := c.reportScheduleService.Index(ctx, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Data:    reportSchedules,
		Message: "Report schedule found successfully",
	})
}

// Create handles POST /api/v1/report-schedules
func (c *ReportScheduleController) Create(ctx *gin.Context) {
	var reportScheduleRequest dto.ReportScheduleRequest
	if err := ctx.ShouldBindJSON(&reportScheduleRequest); err != nil {
		log.Println("ERROR: ", err)
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
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

	if !helper.IsValidTokenFormat(token) {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Invalid token format",
		})
		return
	}

	// Sanitize the request data
	reportScheduleRequest.AcademicAdvisorEmail = helper.SanitizeString(reportScheduleRequest.AcademicAdvisorEmail)
	reportScheduleRequest.RegistrationID = helper.SanitizeString(reportScheduleRequest.RegistrationID)
	reportScheduleRequest.AcademicAdvisorID = helper.SanitizeString(reportScheduleRequest.AcademicAdvisorID)
	reportScheduleRequest.UserNRP = helper.SanitizeString(reportScheduleRequest.UserNRP)
	reportScheduleRequest.UserID = helper.SanitizeString(reportScheduleRequest.UserID)
	reportScheduleRequest.StartDate = helper.SanitizeString(reportScheduleRequest.StartDate)
	reportScheduleRequest.EndDate = helper.SanitizeString(reportScheduleRequest.EndDate)
	reportScheduleRequest.ReportType = helper.SanitizeString(reportScheduleRequest.ReportType)

	if !helper.ValidateUUID(reportScheduleRequest.RegistrationID) {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Invalid Registration ID format",
		})
		return
	}
	if !helper.ValidateUUID(reportScheduleRequest.AcademicAdvisorID) {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Invalid Academic Advisor ID format",
		})
		return
	}

	reportSchedule, err := c.reportScheduleService.Create(ctx, reportScheduleRequest, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Data:    reportSchedule,
		Message: "Report schedule created successfully",
	})
}

// Update handles PUT /api/v1/report-schedules/:id
func (c *ReportScheduleController) Update(ctx *gin.Context) {
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

	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Token is required",
		})
		return
	}

	if !helper.IsValidTokenFormat(token) {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Invalid token format",
		})
		return
	}

	var reportScheduleRequest dto.ReportScheduleRequest
	if err := ctx.ShouldBindJSON(&reportScheduleRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	// Sanitize the request data
	reportScheduleRequest.AcademicAdvisorEmail = helper.SanitizeString(reportScheduleRequest.AcademicAdvisorEmail)
	reportScheduleRequest.RegistrationID = helper.SanitizeString(reportScheduleRequest.RegistrationID)
	reportScheduleRequest.AcademicAdvisorID = helper.SanitizeString(reportScheduleRequest.AcademicAdvisorID)
	reportScheduleRequest.UserNRP = helper.SanitizeString(reportScheduleRequest.UserNRP)
	reportScheduleRequest.UserID = helper.SanitizeString(reportScheduleRequest.UserID)
	reportScheduleRequest.StartDate = helper.SanitizeString(reportScheduleRequest.StartDate)
	reportScheduleRequest.EndDate = helper.SanitizeString(reportScheduleRequest.EndDate)
	reportScheduleRequest.ReportType = helper.SanitizeString(reportScheduleRequest.ReportType)

	if !helper.ValidateUUID(reportScheduleRequest.RegistrationID) {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Invalid Registration ID format",
		})
		return
	}
	if !helper.ValidateUUID(reportScheduleRequest.AcademicAdvisorID) {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Invalid Academic Advisor ID format",
		})
		return
	}

	err := c.reportScheduleService.Update(ctx, id, reportScheduleRequest, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Message: "Report schedule updated successfully",
	})
}

// Show handles GET /api/v1/report-schedules/:id
func (c *ReportScheduleController) Show(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "ID is required",
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

	reportSchedule, err := c.reportScheduleService.FindByID(ctx, id, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Message: "Report schedule found successfully",
		Data:    reportSchedule,
	})
}

// Destroy handles DELETE /api/v1/report-schedules/:id
func (c *ReportScheduleController) Destroy(ctx *gin.Context) {
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

	sanitizedId := helper.SanitizeString(id)
	if sanitizedId != id {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Invalid characters in ID",
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

	err := c.reportScheduleService.Destroy(ctx, id, token)
	if err != nil {
		if err.Error() == "record not found" {
			ctx.JSON(http.StatusNotFound, dto.Response{
				Status:  dto.STATUS_ERROR,
				Message: "Report schedule not found",
			})
			return
		}

		if err.Error() == "user role not allowed" {
			ctx.JSON(http.StatusForbidden, dto.Response{
				Status:  dto.STATUS_ERROR,
				Message: "Access denied",
			})
			return
		}
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Message: "Report schedule deleted successfully",
	})
}

// FindByRegistrationID handles GET /api/v1/registrations/:id/report-schedules
func (c *ReportScheduleController) FindByRegistrationID(ctx *gin.Context) {
	registrationID := ctx.Param("id")
	if registrationID == "" {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Registration ID is required",
		})
		return
	}

	reportSchedules, err := c.reportScheduleService.FindByRegistrationID(ctx, registrationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Message: "Report schedule found successfully",
		Data:    reportSchedules,
	})
}
