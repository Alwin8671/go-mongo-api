package controller

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Alwin8671/go-mongo-api/config"
	"github.com/Alwin8671/go-mongo-api/models"
	"github.com/Alwin8671/go-mongo-api/responses"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func DeleteUser(c *gin.Context) {
	var id = c.Param("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := config.UserCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Faild to delete user"})
	}
	c.IndentedJSON(http.StatusOK, gin.H{"deleted": res.DeletedCount})
}

func UpdateUser(c *gin.Context) {
	var id = c.Param("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid ID format"})
		return
	}

	var updateData models.User

	if err := c.BindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]any{"data": err.Error()}})
		return
	}

	update := bson.M{"name": updateData.Name, "age": updateData.Age, "email": updateData.Email}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := config.UserCollection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error updating user data"})
		return
	}

	if res.MatchedCount == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	var updatedUser models.User
	if res.MatchedCount == 1 {
		err := config.UserCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&updatedUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
	}

	c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]any{"data": updatedUser}})
}
