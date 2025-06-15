package controller

import (
	"fmt"
	"monitoring-service/dto"
	"monitoring-service/helper"
	"monitoring-service/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TranscriptController struct {
	transcriptService service.TranscriptService
}

func NewTranscriptController(transcriptService service.TranscriptService) *TranscriptController {
	return &TranscriptController{
		transcriptService: transcriptService,
	}
}

func (c *TranscriptController) FindByAdvisorEmail(ctx *gin.Context) {
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
	var filter dto.TranscriptAdvisorFilterRequest
	if err := ctx.ShouldBindJSON(&filter); err != nil {
		// If parsing fails, proceed with empty filter (not a critical error)
		filter = dto.TranscriptAdvisorFilterRequest{}
	}

	transcripts, metaData, err := c.transcriptService.FindByAdvisorEmailAndGroupByUserNRP(ctx, token, pagReq, filter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:             dto.STATUS_SUCCESS,
		Data:               transcripts,
		Message:            "Transcripts fetched successfully",
		PaginationResponse: &metaData,
	})
}

// Index handles GET /api/v1/transcripts
func (c *TranscriptController) Index(ctx *gin.Context) {

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

	transcripts, err := c.transcriptService.Index(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Data:    transcripts,
		Message: "Transcripts fetched successfully",
	})
}

// Create handles POST /api/v1/transcripts
func (c *TranscriptController) Create(ctx *gin.Context) {
	var transcriptRequest dto.TranscriptRequest
	if err := ctx.ShouldBind(&transcriptRequest); err != nil {
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

	transcriptRequest.Title = helper.SanitizeString(transcriptRequest.Title)
	transcriptRequest.RegistrationID = helper.SanitizeString(transcriptRequest.RegistrationID)

	if transcriptRequest.Title == "" {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Title is required",
		})
		return
	}
	if transcriptRequest.RegistrationID == "" {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Registration ID is required",
		})
		return
	}

	if !helper.ValidateUUID(transcriptRequest.RegistrationID) {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Invalid Registration ID format",
		})
		return
	}

	file, err := ctx.FormFile("file")
	if err := helper.ValidateFileUpload(file); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: fmt.Sprintf("File validation failed: %s", err.Error()),
		})
		return
	}

	// fileContent, err := file.Open()
	// if err != nil {
	// 	log.Printf("Error opening file: %v", err)
	// 	ctx.JSON(http.StatusBadRequest, dto.Response{
	// 		Status:  dto.STATUS_ERROR,
	// 		Message: "Unable to process file",
	// 	})
	// 	return
	// }
	// defer fileContent.Close()

	// if err := helper.ValidateMimeType(fileContent); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, dto.Response{
	// 		Status:  dto.STATUS_ERROR,
	// 		Message: err.Error(),
	// 	})
	// 	return
	// }

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
	if !helper.IsValidTokenFormat(token) {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Invalid authorization format",
		})
		return
	}

	transcript, err := c.transcriptService.Create(ctx, transcriptRequest, file, token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Data:    transcript,
		Message: "Transcript created successfully",
	})
}

// Update handles PUT /api/v1/transcripts/:id
func (c *TranscriptController) Update(ctx *gin.Context) {
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

	var transcriptRequest dto.TranscriptRequest
	if err := ctx.ShouldBindJSON(&transcriptRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	transcriptRequest.Title = helper.SanitizeString(transcriptRequest.Title)
	transcriptRequest.RegistrationID = helper.SanitizeString(transcriptRequest.RegistrationID)
	if transcriptRequest.Title == "" {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Title is required",
		})
		return
	}
	if transcriptRequest.RegistrationID == "" {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Registration ID is required",
		})
		return
	}
	if !helper.ValidateUUID(transcriptRequest.RegistrationID) {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Invalid Registration ID format",
		})
		return
	}

	err := c.transcriptService.Update(ctx, id, transcriptRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Message: "Transcript updated successfully",
	})
}

// Show handles GET /api/v1/transcripts/:id
func (c *TranscriptController) Show(ctx *gin.Context) {
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

	transcript, err := c.transcriptService.FindByID(ctx, id, token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Data:    transcript,
		Message: "Transcript fetched successfully",
	})
}

// Destroy handles DELETE /api/v1/transcripts/:id
func (c *TranscriptController) Destroy(ctx *gin.Context) {
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

	err := c.transcriptService.Destroy(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Message: "Transcript deleted successfully",
	})
}

// FindByRegistrationID handles GET /api/v1/registrations/:id/transcripts
func (c *TranscriptController) FindByRegistrationID(ctx *gin.Context) {
	registrationID := ctx.Param("id")
	if registrationID == "" {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Registration ID is required",
		})
		return
	}

	transcripts, err := c.transcriptService.FindByRegistrationID(ctx, registrationID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Data:    transcripts,
		Message: "Transcripts fetched successfully",
	})
}

func (c *TranscriptController) FindByUserNRPAndGroupByRegistrationID(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Token is required",
		})
		return
	}

	transcripts, err := c.transcriptService.FindByUserNRPAndGroupByRegistrationID(ctx, token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Data:    transcripts,
		Message: "Student transcripts fetched successfully",
	})
}

// FindAllByRegistrationID handles GET /api/v1/transcripts/registrations/:id
func (c *TranscriptController) FindAllByRegistrationID(ctx *gin.Context) {
	registrationID := ctx.Param("id")
	if registrationID == "" {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: "Registration ID is required",
		})
		return
	}

	transcripts, err := c.transcriptService.FindAllByRegistrationID(ctx, registrationID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Status:  dto.STATUS_SUCCESS,
		Data:    transcripts,
		Message: "All transcripts for registration fetched successfully",
	})
}
