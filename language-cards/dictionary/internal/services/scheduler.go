package services 


import (
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3" 
	"time"
)

func ResetSeenDaily(db *sql.DB) error {
	_, err := db.Exec("UPDATE flashcard SET seen = false")
	return err
}

func StartDailyResetTask(db *sql.DB) {
	go func() {
		for {
			now:= time.Now()
			nextMidnight := now.AddDate(0, 0, 1).Truncate(24 * time.Hour)
			durationUntilMidnight := time.Until(nextMidnight)

			time.Sleep(durationUntilMidnight)

			if err := ResetSeenDaily(db); err != nil {
				log.Println("Error resetting seen status:", err)
			}
		}
	}()
}