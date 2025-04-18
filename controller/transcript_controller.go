package controller

import (
	"monitoring-service/dto"
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

// Index handles GET /api/v1/transcripts
func (c *TranscriptController) Index(ctx *gin.Context) {
	transcripts, err := c.transcriptService.Index(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
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

	transcript, err := c.transcriptService.Create(ctx, transcriptRequest, file, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
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

	var transcriptRequest dto.TranscriptRequest
	if err := ctx.ShouldBindJSON(&transcriptRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Status:  dto.STATUS_ERROR,
			Message: err.Error(),
		})
		return
	}

	err := c.transcriptService.Update(ctx, id, transcriptRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
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

	transcript, err := c.transcriptService.FindByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
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

	err := c.transcriptService.Destroy(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
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
		ctx.JSON(http.StatusInternalServerError, dto.Response{
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
