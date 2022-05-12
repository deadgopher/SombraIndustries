package main

import (
	"database/sql"
	"fmt"
	"germ/controller"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

// func generateKeys() {
// 	filename := "key"
// 	bitSize := 4096

// 	// Generate RSA key.
// 	key, err := rsa.GenerateKey(rand.Reader, bitSize)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// Extract public component.
// 	pub := key.Public()

// 	// Encode private key to PKCS#1 ASN.1 PEM.
// 	keyPEM := pem.EncodeToMemory(
// 		&pem.Block{
// 			Type:  "RSA PRIVATE KEY",
// 			Bytes: x509.MarshalPKCS1PrivateKey(key),
// 		},
// 	)

// 	// Encode public key to PKCS#1 ASN.1 PEM.
// 	pubPEM := pem.EncodeToMemory(
// 		&pem.Block{
// 			Type:  "RSA PUBLIC KEY",
// 			Bytes: x509.MarshalPKCS1PublicKey(pub.(*rsa.PublicKey)),
// 		},
// 	)

// 	// Write private key to file.
// 	if err := ioutil.WriteFile(filename+".rsa", keyPEM, 0700); err != nil {
// 		panic(err)
// 	}

// 	// Write public key to file.
// 	if err := ioutil.WriteFile(filename+".rsa.pub", pubPEM, 0755); err != nil {
// 		panic(err)
// 	}
// }

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASSWORD"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "eveHQ",
	}
	// Get a database handle.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("db connected")
	defer db.Close()

	control := controller.New(db)
	server := gin.Default()

	server.Use(static.Serve("/static", static.LocalFile("./static", true)))
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"https://localhost:8080", "http://localhost:3000"}
	corsConfig.AllowHeaders = []string{"Content-Type"}
	server.Use(cors.New(corsConfig))
	control.Register(server.Group("api"))

	server.RunTLS(":8080", "./localhost.pem", "./localhost-key.pem")

}
