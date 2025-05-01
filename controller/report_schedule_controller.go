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

	reportSchedules, metaData, err := c.reportScheduleService.FindByAdvisorEmail(ctx, token, pagReq)
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

	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Token is required",
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
