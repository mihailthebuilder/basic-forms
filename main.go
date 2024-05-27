package main

import (
	"database/sql"
	"flag"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func main() {
	if isRunningLocally() {
		loadEnvironmentVariablesFromDotEnvFile()
	}

	runApplication()
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

func runApplication() {
	reset := flag.Bool("reset", false, "resets the database")
	flag.Parse()

	r := gin.Default()

	serverRecoversFromAnyPanicAndWrites500(r)
	allowAllOriginsForCORS(r)

	initDb(reset)
	defer db.Close()

	setUpRoutes(r)

	r.Run("localhost:8080")
}

func serverRecoversFromAnyPanicAndWrites500(engine *gin.Engine) {
	engine.Use(gin.Recovery())
}

func allowAllOriginsForCORS(engine *gin.Engine) {
	engine.Use(cors.Default())
}

var db *sql.DB

func initDb(reset *bool) {
	if *reset {
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
			origin_url text not null
		);

		create index if not exists idx_submission_user_id on submission(user_id, origin_url);
	`

	_, err = db.Exec(sql)
	if err != nil {
		log.Fatalf("error initializing table: %s", err)
	}
}
