package routes

import (
	"monitoring-service/controller"
	"monitoring-service/middleware"
	"monitoring-service/service"

	"github.com/gin-gonic/gin"
)

func SyllabusRoutes(router *gin.Engine, syllabusController controller.SyllabusController, userManagementService service.UserManagementService) {
	authMiddleware := middleware.AuthorizationRole(userManagementService, []string{"ADMIN", "DOSEN PEMBIMBING", "MAHASISWA"})
	adminMiddleware := middleware.AuthorizationRole(userManagementService, []string{"ADMIN"})
	advisorMiddleware := middleware.AuthorizationRole(userManagementService, []string{"DOSEN PEMBIMBING"})
	studentMiddleware := middleware.AuthorizationRole(userManagementService, []string{"MAHASISWA"})

	syllabusRoutes := router.Group("/monitoring-service/api/v1/syllabuses")
	{
		syllabusRoutes.GET("", adminMiddleware, syllabusController.Index)
		syllabusRoutes.GET("/advisor", advisorMiddleware, syllabusController.FindByAdvisorEmail)
		syllabusRoutes.GET("/student", studentMiddleware, syllabusController.FindByUserNRPAndGroupByRegistrationID)
		syllabusRoutes.GET("/registrations/:id", authMiddleware, syllabusController.FindAllByRegistrationID)
		syllabusRoutes.GET("/registrations/:id/syllabuses", authMiddleware, syllabusController.FindByRegistrationID)
		syllabusRoutes.GET("/:id", authMiddleware, syllabusController.Show)

		authorized := syllabusRoutes.Group("")
		authorized.Use(authMiddleware)
		{
			authorized.POST("", syllabusController.Create)
			authorized.PUT("/:id", syllabusController.Update)
			authorized.DELETE("/:id", syllabusController.Destroy)
		}
	}
}
