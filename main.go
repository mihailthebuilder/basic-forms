package main

import (
	"database/sql"
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

	initDb(shouldResetDatabase())
	defer db.Close()

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

var db *sql.DB

func initDb(reset bool) {
	if reset {
		os.Remove("./sqlite.db")
	}

	var err error
	db, err = sql.Open("sqlite", "./sqlite.db")
	if err != nil {
		log.Fatalf("error starting database: %s", err)
	}

	sql := `
		create table if not exists submission (
			user_id text not null, 
			content text not null, 
			origin text not null
		);

		create index if not exists idx_submission_user_id on submission(user_id, origin);
	`

	_, err = db.Exec(sql)
	if err != nil {
		log.Fatalf("error initializing table: %s", err)
	}
}
