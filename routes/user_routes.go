package routes

import (
	"github.com/Alwin8671/go-mongo-api/controller"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine){
	router.GET("/users", controller.GetUsers)
	router.POST("/users", controller.CreateUser)
	router.DELETE("/users/:id", controller.DeleteUser)
	router.PUT("/users/:id", controller.UpdateUser)
}