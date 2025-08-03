package routes

import (
	"SimpleTaskManager/controllers"

	"github.com/gin-gonic/gin"
)

func TaskRoute(r *gin.Engine) {
	taskGroup := r.Group("/task")
	{
		taskGroup.POST("/:id/createTask", controllers.CreateTask)
		taskGroup.GET("/:id", controllers.GetTask)
		taskGroup.GET("/getAllTask", controllers.GetAllTask)
		taskGroup.PUT("/updateTask", controllers.UpdateTask)
		taskGroup.PUT("/finishedTask/:id", controllers.FinishedTask)
		taskGroup.DELETE("/deleteTask/:id", controllers.Delete)

		taskGroup.POST("/assigned", controllers.ChooseAssignee)
	}

}
