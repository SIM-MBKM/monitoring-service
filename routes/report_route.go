package routes

import (
	"monitoring-service/controller"
	"monitoring-service/middleware"
	"monitoring-service/service"

	"github.com/gin-gonic/gin"
)

func ReportRoutes(router *gin.Engine, reportController controller.ReportController, userManagementService service.UserManagementService) {
	authMiddleware := middleware.AuthorizationRole(userManagementService, []string{"ADMIN", "DOSEN PEMBIMBING", "MAHASISWA"})
	adminMiddleware := middleware.AuthorizationRole(userManagementService, []string{"ADMIN"})
	advisorMiddleware := middleware.AuthorizationRole(userManagementService, []string{"DOSEN PEMBIMBING"})

	reportRoutes := router.Group("/monitoring-service/api/v1/reports")
	{
		reportRoutes.GET("", adminMiddleware, reportController.Index)
		reportRoutes.GET("/report-schedules/:id/reports", reportController.FindByReportScheduleID)
		reportRoutes.POST("/approval/:id", advisorMiddleware, reportController.Approval)

		authorized := reportRoutes.Group("")
		authorized.Use(authMiddleware)
		{
			reportRoutes.GET("/:id", reportController.Show)
			authorized.POST("", reportController.Create)
			authorized.PUT("/:id", reportController.Update)
			authorized.DELETE("/:id", reportController.Destroy)
		}
	}
}
