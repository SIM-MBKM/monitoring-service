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

	if !helper.IsValidTokenFormat(token) {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Invalid authorization format",
		})
		return
	}

	_, _, err := helper.ValidatePaginationParams(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
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

	filter.UserNRP = helper.SanitizeString(filter.UserNRP)

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

	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Authorization required",
		})
		return
	}

	if !helper.IsValidTokenFormat(token) {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Invalid authorization format",
		})
		return
	}

	_, _, err := helper.ValidatePaginationParams(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

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

	if ctx.Request.ContentLength > helper.MaxFileSize+helper.MaxContentLength { // +1KB for form data
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Request too large",
		})
		return
	}

	syllabusRequest.Title = helper.SanitizeString(syllabusRequest.Title)

	file, err := ctx.FormFile("file")

	if err := helper.ValidateFileUpload(file); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: fmt.Sprintf("File validation failed: %s", err.Error()),
		})
		return
	}

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
			Message: "Invalid authorization format",
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
			Message: "Invalid authorization format",
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
			Message: "Invalid authorization format",
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
	if !helper.ValidateUUID(registrationID) {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Invalid Registration ID format",
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
			Message: "Invalid authorization format",
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

	if !helper.ValidateUUID(registrationID) {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Invalid Registration ID format",
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
