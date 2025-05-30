package main

import (
	"log"
	"os"

	"github.com/dzahariev/solei/api/controller"
	"github.com/joho/godotenv"
)

var server = controller.Server{}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not loaded due to:", err)
	}

	// Auth configuration
	authURL := os.Getenv("AUTH_URL")
	authRealm := os.Getenv("AUTH_REALM")
	authClientID := os.Getenv("AUTH_CLIENT_ID")
	authClientSecret := os.Getenv("AUTH_CLIENT_SECRET")

	// Role to Permissions map in YAML format
	roleToPermissionsYaml := os.Getenv("PERMISSIONS_YAML")

	// DB configuration
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbName := os.Getenv("POSTGRES_DB")

	server.Initialize(dbUser, dbPassword, dbPort, dbHost, dbName, authURL, authRealm, authClientID, authClientSecret, roleToPermissionsYaml)
	server.Run(":8800")
}
