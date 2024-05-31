package main

import (
	"fmt"
	"log"
	"os"

	"basic-forms/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func main() {
	host := "0.0.0.0"

	if isRunningLocally() {
		loadEnvironmentVariablesFromDotEnvFile()
		host = "localhost"
	}

	runApplication(host)
}

func isRunningLocally() bool {
	return os.Getenv("GIN_MODE") == ""
}

func loadEnvironmentVariablesFromDotEnvFile() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func runApplication(host string) {
	r := gin.Default()

	serverRecoversFromAnyPanicAndWrites500(r)
	allowAllOriginsForCORS(r)

	initDatastore(shouldResetDatabase())

	r.Use(logger.SetLoggerContext)

	setUpRoutes(r)

	r.Run(fmt.Sprintf("%s:8080", host))
}

func shouldResetDatabase() bool {
	return os.Getenv("RESET") == "true"
}

func serverRecoversFromAnyPanicAndWrites500(engine *gin.Engine) {
	engine.Use(gin.Recovery())
}

func allowAllOriginsForCORS(engine *gin.Engine) {
	engine.Use(cors.Default())
}

var datastore Datastore

func initDatastore(reset bool) {
	if reset {
		os.RemoveAll("./users/")
	}

	secret := os.Getenv("SECRET")
	if len(secret) == 0 {
		log.Fatal("error loading secret env var")
	}

	datastore = Datastore{Secret: secret}
}
