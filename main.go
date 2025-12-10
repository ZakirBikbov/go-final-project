package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"final-project/pkg/db"
	"final-project/pkg/server"
)

func main() {
	webDir := filepath.Join(".", "web")

	dbFile := "scheduler.db"
	if env := os.Getenv("TODO_DBFILE"); env != "" {
		dbFile = env
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}
	defer db.DB.Close()

	defaultPort := 7540
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if parsed, err := strconv.Atoi(envPort); err == nil && parsed > 0 && parsed < 65536 {
			defaultPort = parsed
		}
	}

	if err := server.Start(webDir, defaultPort); err != nil {
		log.Fatal(err)
	}
}
