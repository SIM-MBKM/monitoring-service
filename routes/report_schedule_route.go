package routes

import (
	"monitoring-service/controller"
	"monitoring-service/middleware"
	"monitoring-service/service"

	"github.com/gin-gonic/gin"
)

func ReportScheduleRoutes(router *gin.Engine, reportScheduleController controller.ReportScheduleController, userManagementService service.UserManagementService) {
	authMiddleware := middleware.AuthorizationRole(userManagementService, []string{"ADMIN", "DOSEN PEMBIMBING", "MAHASISWA"})
	adminMiddleware := middleware.AuthorizationRole(userManagementService, []string{"ADMIN"})
	advisorMiddleware := middleware.AuthorizationRole(userManagementService, []string{"DOSEN PEMBIMBING"})
	studentMiddleware := middleware.AuthorizationRole(userManagementService, []string{"MAHASISWA"})

	reportScheduleRoutes := router.Group("/monitoring-service/api/v1/report-schedules")
	{
		reportScheduleRoutes.GET("/", adminMiddleware, reportScheduleController.Index)
		reportScheduleRoutes.GET("/student", studentMiddleware, reportScheduleController.FindByStudentID)
		reportScheduleRoutes.GET("/advisor", advisorMiddleware, reportScheduleController.FindByAdvisorEmail)
		reportScheduleRoutes.GET("/:id", authMiddleware, reportScheduleController.Show)
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
