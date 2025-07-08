package main

import (
	"github.com/Alwin8671/go-mongo-api/config"
	"github.com/Alwin8671/go-mongo-api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase();
	router := gin.Default()
	routes.UserRoutes(router);
	router.Run("localhost:8080")
}
