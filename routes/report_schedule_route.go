package routes

import (
	"monitoring-service/controller"
	"monitoring-service/middleware"
	"monitoring-service/service"

	"github.com/gin-gonic/gin"
)

func ReportScheduleRoutes(router *gin.Engine, reportScheduleController controller.ReportScheduleController, userManagementService service.UserManagementService) {
	authMiddleware := middleware.AuthorizationRole(userManagementService, []string{"ADMIN", "DOSEN", "MAHASISWA"})

	reportScheduleRoutes := router.Group("/monitoring-service/api/v1/report-schedules")
	{
		reportScheduleRoutes.GET("/", reportScheduleController.Index)
		reportScheduleRoutes.GET("/:id", reportScheduleController.Show)
		reportScheduleRoutes.GET("/registrations/:id/report-schedules", reportScheduleController.FindByRegistrationID)

		authorized := reportScheduleRoutes.Group("/")
		authorized.Use(authMiddleware)
		{
			authorized.POST("/", reportScheduleController.Create)
			authorized.PUT("/:id", reportScheduleController.Update)
			authorized.DELETE("/:id", reportScheduleController.Destroy)
		}
	}
}
