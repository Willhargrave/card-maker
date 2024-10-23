package database


import (
	"fmt"
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3" 
)
func CreateUserTable(db *sql.DB) {
    createTableSQL := `CREATE TABLE IF NOT EXISTS Users (
        "UserID" INTEGER PRIMARY KEY AUTOINCREMENT,
        "username" TEXT UNIQUE NOT NULL,
        "password" TEXT NOT NULL
    );`

    _, err := db.Exec(createTableSQL)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("User table created")
}

func CreateFlashcardTable(db *sql.DB) {
	
	createTableSQL := `CREATE TABLE IF NOT EXISTS flashcard (
		"word" TEXT NOT NULL,
		"definition" TEXT,
		"seen" BOOLEAN NOT NULL DEFAULT false,
		"UserID" INTEGER,
		"setId" INTEGER,
		PRIMARY KEY ("word", "UserID"),
		FOREIGN KEY ("setId") REFERENCES FlashcardSets("setId")
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Flashcard Table created")
}

func CreateFlashcardSetsTable(db *sql.DB) {

	createTableSQL := `CREATE TABLE IF NOT EXISTS FlashcardSets (
		"setId" INTEGER PRIMARY KEY AUTOINCREMENT,
		"UserID" INTEGER,
		"setName" TEXT NOT NULL,
		FOREIGN KEY ("UserID") REFERENCES Users("UserID")
	);`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("FlashcardSets table created")
}