package database

import (
	"fmt"
	"database/sql"
)
var SeedData = []struct {
    Word       string
    Definition string
}{
    {"こんにちは", "Hello"},
    {"ありがとう", "Thank you"},
    {"さようなら", "Goodbye"},
    // Add more words here...
}

func SeedDatabase(db *sql.DB) error {
    for _, entry := range SeedData {
        _, err := db.Exec("INSERT INTO flashcard (word, definition) VALUES (?, ?)", entry.Word, entry.Definition)
        if err != nil {
            return fmt.Errorf("failed to seed word '%s': %w", entry.Word, err)
        }
    }
    fmt.Println("Database seeded successfully.")
    return nil
}
