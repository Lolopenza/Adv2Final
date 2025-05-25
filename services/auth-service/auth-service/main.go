package main

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
	//MongoDB
	clientOptions := options.Client().ApplyURI("mongodb+srv://user:123@cluster0.961p8.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	//Connection check
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")

	InitUsersCollection()

	// Запускаем оба сервера в отдельных горутинах
	var wg sync.WaitGroup
	wg.Add(2)

	// Запуск gRPC сервера
	go func() {
		defer wg.Done()
		log.Println("Starting gRPC server...")
		if err := startGRPCServer(); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// Запуск REST API сервера
	go func() {
		defer wg.Done()
		startRESTServer()
	}()

	wg.Wait()
}

func startRESTServer() {
	//Gin
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	//routes for static files
	router.Static("/assets", "../frontend/dist/assets")
	router.StaticFile("/favicon.ico", "../frontend/dist/favicon.ico")

	//routes
	router.GET("/", func(c *gin.Context) {
		c.File("../frontend/dist/index.html")
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	//API routes
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", RegisterHandler)
			auth.POST("/login", LoginHandler)
			auth.GET("/me", AuthMiddleware(), UserInfoHandler)
			auth.POST("/logout", LogoutHandler)
			auth.PUT("/profile", AuthMiddleware(), UpdateProfileHandler)
			auth.PUT("/password", AuthMiddleware(), ChangePasswordHandler)
		}
	}

	//server
	log.Println("REST API Server running on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start REST server: %v", err)
	}
}
