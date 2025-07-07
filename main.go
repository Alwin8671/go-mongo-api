package main

import (
	"github.com/Alwin8671/go-mongo-api/config"
	"github.com/Alwin8671/go-mongo-api/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase();
	router := gin.Default()
	router.GET("/users", controller.GetUsers)
	router.POST("/users", controller.CreateUser)
	router.Run("localhost:8080")
}
