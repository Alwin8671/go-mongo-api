package controller

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Alwin8671/go-mongo-api/config"
	"github.com/Alwin8671/go-mongo-api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := config.UserCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var users []models.User
	for cursor.Next(ctx) {
		var user models.User
		if err = cursor.Decode(&user); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	c.IndentedJSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
	var newUser models.User
	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	res, err := config.UserCollection.InsertOne(context.TODO(), newUser)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to insert user"})
		log.Println(err)
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"InsertedId": res.InsertedID})
}
