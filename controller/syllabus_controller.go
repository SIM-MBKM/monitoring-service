package controller

import (
	"monitoring-service/dto"
	"monitoring-service/helper"
	"monitoring-service/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SyllabusController struct {
	syllabusService service.SyllabusService
}

func NewSyllabusController(syllabusService service.SyllabusService) *SyllabusController {
	return &SyllabusController{
		syllabusService: syllabusService,
	}
}

func (c *SyllabusController) FindByAdvisorEmail(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Token is required",
		})
		return
	}

	// Parse pagination parameters
	pagReq := helper.Pagination(ctx)

	// Parse filter from request body
	var filter dto.SyllabusAdvisorFilterRequest
	if err := ctx.ShouldBindJSON(&filter); err != nil {
		// If parsing fails, proceed with empty filter (not a critical error)
		filter = dto.SyllabusAdvisorFilterRequest{}
	}

	syllabuses, metaData, err := c.syllabusService.FindByAdvisorEmailAndGroupByUserNRP(ctx, token, pagReq, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:             dto.STATUS_SUCCESS,
		Data:               syllabuses,
		Message:            "Syllabuses fetched successfully",
		PaginationResponse: &metaData,
	})
}

// Index handles GET /api/v1/syllabuses
func (c *SyllabusController) Index(ctx *gin.Context) {
	syllabuses, err := c.syllabusService.Index(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Data:    syllabuses,
		Message: "Syllabuses fetched successfully",
	})
}

// Create handles POST /api/v1/syllabuses
func (c *SyllabusController) Create(ctx *gin.Context) {
	var syllabusRequest dto.SyllabusRequest
	if err := ctx.ShouldBind(&syllabusRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	file, err := ctx.FormFile("file")
	if file == nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "File is required",
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

	syllabus, err := c.syllabusService.Create(ctx, syllabusRequest, file, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Data:    syllabus,
		Message: "Syllabus created successfully",
	})
}

// Update handles PUT /api/v1/syllabuses/:id
func (c *SyllabusController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "ID is required",
		})
		return
	}

	var syllabusRequest dto.SyllabusRequest
	if err := ctx.ShouldBindJSON(&syllabusRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	err := c.syllabusService.Update(ctx, id, syllabusRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Message: "Syllabus updated successfully",
	})
}

// Show handles GET /api/v1/syllabuses/:id
func (c *SyllabusController) Show(ctx *gin.Context) {
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

	syllabus, err := c.syllabusService.FindByID(ctx, id, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Data:    syllabus,
		Message: "Syllabus fetched successfully",
	})
}

// Destroy handles DELETE /api/v1/syllabuses/:id
func (c *SyllabusController) Destroy(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "ID is required",
		})
		return
	}

	err := c.syllabusService.Destroy(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Message: "Syllabus deleted successfully",
	})
}

// FindByRegistrationID handles GET /api/v1/registrations/:id/syllabuses
func (c *SyllabusController) FindByRegistrationID(ctx *gin.Context) {
	registrationID := ctx.Param("id")
	if registrationID == "" {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Registration ID is required",
		})
		return
	}

	syllabuses, err := c.syllabusService.FindByRegistrationID(ctx, registrationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Data:    syllabuses,
		Message: "Syllabus fetched successfully",
	})
}

// FindAllByRegistrationID handles GET /api/v1/syllabuses/registrations/:id
func (c *SyllabusController) FindAllByRegistrationID(ctx *gin.Context) {
	registrationID := ctx.Param("id")
	if registrationID == "" {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Registration ID is required",
		})
		return
	}

	syllabuses, err := c.syllabusService.FindAllByRegistrationID(ctx, registrationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Data:    syllabuses,
		Message: "All syllabuses for registration fetched successfully",
	})
}

func (c *SyllabusController) FindByUserNRPAndGroupByRegistrationID(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Token is required",
		})
		return
	}

	syllabuses, err := c.syllabusService.FindByUserNRPAndGroupByRegistrationID(ctx, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Data:    syllabuses,
		Message: "Student syllabuses fetched successfully",
	})
}
