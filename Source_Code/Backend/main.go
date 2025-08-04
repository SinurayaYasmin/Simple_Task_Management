package main

import (
	"SimpleTaskManager/config"
	"SimpleTaskManager/routes"
	"time"

	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"Message": "Pong"})
	})

	config.Init()
	routes.UserRoute(r)
	routes.TaskRoute(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default local
	}
	r.Run(":" + port)

}
