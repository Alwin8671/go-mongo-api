package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/joho/godotenv"
)

type User struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Age   int                `json:"age" bson:"age"`
	Email string             `json:"email" bson:"email"`
}

var userCollection *mongo.Collection

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	URI := os.Getenv("MONGO_URI")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil {
		log.Fatal("MongoDB connection error:", err)
	}

	userCollection = client.Database("userdb").Collection("users")
	log.Println("Connected to mongodb")
}

func getUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var users []User
	for cursor.Next(ctx) {
		var user User
		if err = cursor.Decode(&user); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	c.IndentedJSON(http.StatusOK, users)
}

func createUser(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		return
	}

	res, err := userCollection.InsertOne(context.TODO(), newUser)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to insert user"})
		log.Println(err)
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"InsertedId": res.InsertedID})
}

func main() {
	ConnectDatabase()
	router := gin.Default()
	router.GET("/users", getUsers)
	router.POST("/users", createUser)
	router.Run("localhost:8080")
}
