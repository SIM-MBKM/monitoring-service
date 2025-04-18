package routes

import (
	"monitoring-service/controller"
	"monitoring-service/middleware"
	"monitoring-service/service"

	"github.com/gin-gonic/gin"
)

func TranscriptRoutes(router *gin.Engine, transcriptController controller.TranscriptController, userManagementService service.UserManagementService) {
	authMiddleware := middleware.AuthorizationRole(userManagementService, []string{"ADMIN", "DOSEN PEMBIMBING", "MAHASISWA"})
	adminMiddleware := middleware.AuthorizationRole(userManagementService, []string{"ADMIN"})

	transcriptRoutes := router.Group("/monitoring-service/api/v1/transcripts")
	{
		transcriptRoutes.GET("/", adminMiddleware, transcriptController.Index)
		transcriptRoutes.GET("/registrations/:id/transcripts", transcriptController.FindByRegistrationID)
		transcriptRoutes.GET("/:id", authMiddleware, transcriptController.Show)

		authorized := transcriptRoutes.Group("/")
		authorized.Use(authMiddleware)
		{
			authorized.POST("/", transcriptController.Create)
			authorized.PUT("/:id", transcriptController.Update)
			authorized.DELETE("/:id", transcriptController.Destroy)
		}
	}
}
