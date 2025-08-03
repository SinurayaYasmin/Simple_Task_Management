package routes

import (
	"SimpleTaskManager/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoute(r *gin.Engine) {
	userGroup := r.Group("/users")
	{
		userGroup.POST("/signUp", controllers.CreateAccount)
		userGroup.POST("/signIn", controllers.SignIn)
		userGroup.GET("/:id", controllers.GetUser)
		userGroup.GET("/getAllUser", controllers.GetAllUser)
	}

}
