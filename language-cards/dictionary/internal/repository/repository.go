package repository 


import (
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3" 
	"errors"
	"dictionary/internal/models"
	"dictionary/internal/services"
   
)

func SearchDB(db *sql.DB, word string) (string, error) {
	var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM flashcard WHERE word = ?)", word).Scan(&exists)
	if err != nil {
	   return "", err
   }
   if !exists {
	   return "", services.ErrWordDoesNotExist
   }
   var definition string
	   err = db.QueryRow("SELECT definition FROM flashcard WHERE word = ?", word).Scan(&definition)
   if err != nil {
	   log.Printf("Error retrieving definition for word '%s': %v", word, err)
	  return "", err
   }
   return definition, nil
   }
   
   
   func AddWordToDB(db *sql.DB, word, definition string, setId, UserID int) error {
	   var exists bool
	   err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM flashcard WHERE word = ? AND UserID = ?)", word, UserID).Scan(&exists)
	   if err != nil {
		   return err
	   }
	   if exists {
		   return services.ErrWordExists
	   }
	   _, err = db.Exec("INSERT INTO flashcard (word, definition, UserID, setId) VALUES (?, ?, ?, ?)", word, definition, UserID, setId)
	   return err
   }

   
   func UpdateWordToDB(db *sql.DB, originalWord, newWord, newDefinition string, UserID, setId int) error {
	   var exists bool
	   err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM flashcard WHERE word = ? AND UserID = ? AND setId = ?)", originalWord, UserID, setId).Scan(&exists)	
	   if err != nil {
		   return err
	   }
	   if !exists {
		   return errors.New("word does not exist for this user and set")
	   }
   
	   _, err = db.Exec("UPDATE flashcard SET word = ?, definition = ? WHERE word = ? AND UserID = ? AND setId = ?", newWord, newDefinition, originalWord, UserID, setId)
	   return err
   }
   
   func DeleteWordFromDB(db *sql.DB, word string, UserID, setId int) error{
	   var exists bool
	   err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM flashcard WHERE word = ? AND UserID = ? AND setId = ?)", word, UserID, setId).Scan(&exists)
	   if err != nil {
		   return err
	   }
	   if !exists {
		   return services.ErrWordDoesNotExist
	   }
   
	   _, err = db.Exec("DELETE FROM flashcard WHERE word = ? AND UserID = ? AND setId = ?", word, UserID, setId)
	   return err
   }

   func ShowAllWordsFromDB(db *sql.DB, UserID, setId int) ([]models.WordDefinition, error) {
	   var rows *sql.Rows
	   var err error
   
	   if setId > 0 {
		   rows, err = db.Query("SELECT word, definition, seen FROM flashcard WHERE UserID = ? AND setId = ?", UserID, setId)
	   } else {
		   rows, err = db.Query("SELECT word, definition, seen FROM flashcard WHERE UserID = ?", UserID)
	   }
   
	   if err != nil {
		   return nil, err
	   }
	   defer rows.Close()
   
	   var words []models.WordDefinition
	   for rows.Next() {
		   var wd models.WordDefinition
		   if err := rows.Scan(&wd.Word, &wd.Definition, &wd.Seen); err != nil {
			   return nil, err
		   }
		   words = append(words, wd)
	   }
	   if err = rows.Err(); err != nil {
		   return nil, err
	   }
	   return words, nil
   }
   
   
   func UpdateSeenInDB(db *sql.DB, word string, userID, setId int) error {
	   var query string
	   if setId > 0 {
		   query = "UPDATE flashcard SET seen = true WHERE word = ? AND userID = ? AND setId = ?"
		   _, err := db.Exec(query, word, userID, setId)
		   return err
	   }
	   query = "UPDATE flashcard SET seen = true WHERE word = ? AND userID = ?"
	   _, err := db.Exec(query, word, userID)
	   return err
   }
   
   
   func CreateSetInDB(db *sql.DB, userID int, setName string) error {
	   stmt, err := db.Prepare("INSERT INTO flashcardSets (userID, setName) VALUES (?, ?)")
   
	   if err != nil {
		   return err
	   }
	   defer stmt.Close()
   
	   _, err = stmt.Exec(userID, setName)
	   return err
   }
   
   func FetchUserSetsFromDB(db *sql.DB, UserID int) ([]models.Set, error) {
	   log.Printf("Fetching sets for UserID: %d\n", UserID)
   
	   rows, err := db.Query("SELECT setId, setName FROM flashcardSets WHERE UserID = ?", UserID)
	   if err != nil {
		   log.Printf("Error executing query: %v\n", err)
		   return nil, err
	   }
	   defer rows.Close()
   
	   var sets []models.Set
	   for rows.Next() {
		   var set models.Set
		   if err := rows.Scan(&set.SetId, &set.SetName); err != nil {
			   log.Printf("Error scanning row: %v\n", err)
			   return nil, err
		   }
		   sets = append(sets, set)
	   }
   
	   if err = rows.Err(); err != nil {
		   log.Printf("Error iterating over rows: %v\n", err)
		   return nil, err
	   }
   
	   log.Printf("Fetched sets: %+v\n", sets)
   
	   return sets, nil
   }