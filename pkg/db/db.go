package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	install := false
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("stat db file: %w", err)
		}
		install = true
	}

	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}

	if _, err := DB.Exec(schema); err != nil {
		return fmt.Errorf("ensure schema: %w", err)
	}

	if install {
		log.Printf("database created at %s", dbFile)
	}

	return nil
}
