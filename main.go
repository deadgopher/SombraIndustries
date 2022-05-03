package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"germ/controller"
	"io/ioutil"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func generateKeys() {
	filename := "key"
	bitSize := 4096

	// Generate RSA key.
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		panic(err)
	}
	// Extract public component.
	pub := key.Public()

	// Encode private key to PKCS#1 ASN.1 PEM.
	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	// Encode public key to PKCS#1 ASN.1 PEM.
	pubPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pub.(*rsa.PublicKey)),
		},
	)

	// Write private key to file.
	if err := ioutil.WriteFile(filename+".rsa", keyPEM, 0700); err != nil {
		panic(err)
	}

	// Write public key to file.
	if err := ioutil.WriteFile(filename+".rsa.pub", pubPEM, 0755); err != nil {
		panic(err)
	}
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://admin:admin@playtime.3hqsl.gcp.mongodb.net/testDB?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer client.Disconnect(ctx)
	fmt.Println("Database connected")
	db := client.Database("testDB")

	// create controllers
	fmt.Println("going to try to create the controller now...")
	con := controller.New(db, ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("controller created")
	// init routes
	server := gin.Default()

	fmt.Println("server created")
	server.Static("/assets", "./assets")
	server.StaticFile("/favicon.ico", "./resources/zombie.ico")

	fmt.Println("static files set")
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://localhost:8080"}
	corsConfig.AllowHeaders = []string{"Authorization", "Content-Type"}
	server.Use(cors.New(corsConfig))

	fmt.Println("CORS set")
	api := server.Group("/api")
	{
		con.RegisterPilots(api.Group("/pilots"))
		con.RegisterPosts(api.Group("/posts"))
		con.RegisterComments(api.Group("/comments"))
		con.RegisterEve(api.Group("/eve"))
	}

	con.RegisterViews(server.Group("/"))

	fmt.Println("routes registered")
	server.RunTLS(":8080", "./localhost.pem", "./localhost-key.pem")

}
