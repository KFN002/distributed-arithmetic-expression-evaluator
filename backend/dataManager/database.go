package dataManager

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func test_db() {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	_, err = db.Exec(``)
	fmt.Println("working with database")
}
