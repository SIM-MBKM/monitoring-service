package routes

import (
	"monitoring-service/controller"
	"monitoring-service/middleware"
	"monitoring-service/service"

	"github.com/gin-gonic/gin"
)

func ReportRoutes(router *gin.Engine, reportController controller.ReportController, userManagementService service.UserManagementService) {
	authMiddleware := middleware.AuthorizationRole(userManagementService, []string{"ADMIN", "DOSEN PEMBIMBING", "MAHASISWA"})

	reportRoutes := router.Group("/monitoring-service/api/v1/reports")
	{
		reportRoutes.GET("/", reportController.Index)
		reportRoutes.GET("/:id", reportController.Show)
		reportRoutes.GET("/report-schedules/:id/reports", reportController.FindByReportScheduleID)

		authorized := reportRoutes.Group("/")
		authorized.Use(authMiddleware)
		{
			authorized.POST("/", reportController.Create)
			authorized.PUT("/:id", reportController.Update)
			authorized.DELETE("/:id", reportController.Destroy)
		}
	}
}
