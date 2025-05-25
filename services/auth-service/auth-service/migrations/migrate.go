package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	clientOptions := options.Client().ApplyURI("mongodb+srv://user:123@cluster0.961p8.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")

	userCollection := client.Database("AuthService").Collection("Users")

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = userCollection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatalf("Error creating index: %v", err)
	}
	log.Println("Successfully created index for email field")

	adminUser := bson.M{
		"username":   "admin",
		"email":      "admin@example.com",
		"password":   "8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918", // sha256 хеш пароля "admin"
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}

	var existingAdmin bson.M
	err = userCollection.FindOne(context.TODO(), bson.M{"email": "admin@example.com"}).Decode(&existingAdmin)
	if err == mongo.ErrNoDocuments {
		_, err = userCollection.InsertOne(context.TODO(), adminUser)
		if err != nil {
			log.Fatalf("Error creating admin user: %v", err)
		}
		log.Println("Created admin user")
	} else if err != nil {
		log.Fatalf("Error checking admin user: %v", err)
	} else {
		log.Println("Admin user already exists")
	}

	log.Println("Migration completed successfully!")
}
