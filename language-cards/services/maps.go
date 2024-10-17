package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3" 
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	"errors"
    "strings"
	"strconv"
)

// type Dictionary map[string]string


var jwtKey = []byte("your_secret_key")

// Struct
type Claims struct {
    Username string `json:"username"`
	UserID   int    `json:"UserID"`
    jwt.StandardClaims
}
type WordDefinition struct {
    Word       string `json:"word"`
    Definition string `json:"definition"`
	Seen       bool   `json:"seen"`
}

type Set struct {
    SetId   int    `json:"setId"`
    SetName string `json:"setName"`
}

// Tables
func createUserTable(db *sql.DB) {
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

func createFlashcardTable(db *sql.DB) {
	
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

func createFlashcardSetsTable(db *sql.DB) {

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
//cors
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Allow any origin
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE") // Allow specific methods
    (*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization") // Allow specific headers
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}


//user-login

func registerHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    enableCors(&w)

    if r.Method != "POST" {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    var request struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        log.Println("Error decoding request:", err)
        http.Error(w, "Error reading request", http.StatusBadRequest)
        return 
    }

    if request.Username == "" || request.Password == "" {
        http.Error(w, "Username and password cannot be empty", http.StatusBadRequest)
        return
    }

    hashedPassword, err := HashPassword(request.Password)
    if err != nil {
        log.Println("Error hashing password:", err)
        http.Error(w, "Error processing the password", http.StatusInternalServerError)
        return
    }

    exists, err := usernameExists(request.Username, db)
    if err != nil {
        // Handle the error
        http.Error(w, "Error checking username", http.StatusInternalServerError)
        return
    }
    if exists {
        http.Error(w, "Username already exists", http.StatusBadRequest)
        return
    }

    stmt, err := db.Prepare("INSERT INTO Users(username, password) VALUES(?, ?)")
    if err != nil {
        log.Println("Error preparing SQL statement:", err)
        http.Error(w, "Error processing request", http.StatusInternalServerError)
        return
    }
    defer stmt.Close()

    _, err = stmt.Exec(request.Username, hashedPassword)
    if err != nil {
        log.Println("Error executing SQL statement:", err)
        http.Error(w, "Error processing request", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode("Registered successfully")
}



func loginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    enableCors(&w)

    if r.Method != "POST" {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    var request struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Error reading request", http.StatusBadRequest)
        return
    }

    var userID int
    var storedHash string
    err := db.QueryRow("SELECT UserID, password FROM users WHERE username = ?", request.Username).Scan(&userID, &storedHash)
    if err != nil {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

	log.Printf("Fetched User: Username: %s, UserID: %d\n", request.Username, userID)


    if CheckPasswordHash(request.Password, storedHash) {
        expirationTime := time.Now().Add(30 * time.Minute)

        claims := &Claims{
            UserID: userID, // Now fetching from DB
            Username: request.Username,
            StandardClaims: jwt.StandardClaims{
                ExpiresAt: expirationTime.Unix(),
            },
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

        tokenString, err := token.SignedString(jwtKey)
        if err != nil {
            http.Error(w, "Error creating the token", http.StatusInternalServerError)
            return
        }

        http.SetCookie(w, &http.Cookie{
            Name:    "token",
            Value:   tokenString,
            Expires: expirationTime,
        })

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
    } else {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
    }
}


//functions

func getUserIdFromRequest(r *http.Request) (int, error) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        return 0, errors.New("authorization header is missing")
    }

    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        return 0, errors.New("authorization header format must be 'Bearer {token}'")
    }

    tknStr := parts[1]
    claims := &Claims{}

    // Parse the token
    tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil {
        if err == jwt.ErrSignatureInvalid {
            return 0, err
        }
        return 0, err
    }
    if !tkn.Valid {
        return 0, errors.New("invalid token")
    }

    return claims.UserID, nil
}


func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func usernameExists(username string, db *sql.DB) (bool, error) {
    var exists bool
    err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exists)
    if err != nil {
        log.Println("Error checking if username exists:", err)
        return false, err // Or handle the error differently if needed
    }
    return exists, nil
}

func CheckPasswordHash(password, hashedPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}
//main
func main() {
    // Initialize the SQL database
    db, err := sql.Open("sqlite3", "./dictionary.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
	createUserTable(db)
	createFlashcardTable(db)
	createFlashcardSetsTable(db)


    // Register the root handler
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Welcome")
    })

    // Register the add-word handler
    http.HandleFunc("/add-word", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			enableCors(&w)
			w.WriteHeader(http.StatusOK)
			return
		}
        addWordHandler(w, r, db)
    })
	// Register the update-word handler
    http.HandleFunc("/update-word", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			enableCors(&w)
			w.WriteHeader(http.StatusOK)
			return
		}
        updateWordHandler(w, r, db)
    })
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			enableCors(&w)
			w.WriteHeader(http.StatusOK)
			return
		}
		searchHandler(w, r, db)
	})
	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			enableCors(&w)
			w.WriteHeader(http.StatusOK)
			return
		}
		deleteWordHandler(w, r, db)
	})
	http.HandleFunc("/show-words", func(w http.ResponseWriter, r *http.Request) {
		showWordsHandler(w, r, db)
	})

	http.HandleFunc("/seen", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			enableCors(&w)
			w.WriteHeader(http.StatusOK)
			return
		}
		updateSeenHandler(w, r, db)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			enableCors(&w)
			w.WriteHeader(http.StatusOK)
			return
		}
		registerHandler(w, r, db)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			enableCors(&w)
			w.WriteHeader(http.StatusOK)
			return
		}
		loginHandler(w, r, db)
	})
    
	http.HandleFunc("/create-set", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			enableCors(&w)
			w.WriteHeader(http.StatusOK)
			return
		}
		createSetHandler(w, r, db)
	})

	http.HandleFunc("/fetch-sets", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			enableCors(&w)
			w.WriteHeader(http.StatusOK)
			return
		}
		fetchUserSetsHandler(w, r, db)
	})
	


	// if err := seedDatabase(db); err != nil {
    //     log.Fatal("Error seeding the database:", err)
    // }


	startDailyResetTask(db)
	
    // Only one call to ListenAndServe
    fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("Error starting server: ", err)
    }
}

//search


func searchHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    enableCors(&w)

    if r.Method != "GET" {
        http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
        return
    }

    word := r.URL.Query().Get("word")

    definition, err := SearchDB(db, word) 
    if err != nil {
		log.Printf("Error searching for word '%s': %v", word, err)
		if err == ErrWordDoesNotExist {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(definition) // Ensure this is the correct response format
}

func SearchDB(db *sql.DB, word string) (string, error) {
 var exists bool
 	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM flashcard WHERE word = ?)", word).Scan(&exists)
 if err != nil {
	return "", err
}
if !exists {
	return "", ErrWordDoesNotExist
}
var definition string
	err = db.QueryRow("SELECT definition FROM flashcard WHERE word = ?", word).Scan(&definition)
if err != nil {
	log.Printf("Error retrieving definition for word '%s': %v", word, err)
   return "", err
}
return definition, nil
}

//add

func addWordHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    enableCors(&w)
    

    if r.Method != "POST" {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    var request struct {
        Word       string `json:"word"`
        Definition string `json:"definition"`
        SetID      int    `json:"setId"` 
        UserID     int    `json:"UserID"` 
    }

   

    // Decode the JSON body into the struct
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        log.Println("Error decoding request:", err)
        http.Error(w, "Error reading request", http.StatusBadRequest)
        return 
    }
    
    log.Printf("Received request data: %+v\n", request)

    // Check if the word and definition are not empty
    if request.Word == "" || request.Definition == "" || request.SetID == 0 || request.UserID == 0 {
        http.Error(w, "Required fields cannot be empty", http.StatusBadRequest)
        return
    }


    // Call AddWordToDB to add the word to the database
    err := AddWordToDB(db, request.Word, request.Definition, request.SetID, request.UserID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode("Word added successfully")
}


func AddWordToDB(db *sql.DB, word, definition string, setId, UserID int) error {
	var exists bool
    err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM flashcard WHERE word = ? AND UserID = ?)", word, UserID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return ErrWordExists
	}
    _, err = db.Exec("INSERT INTO flashcard (word, definition, UserID, setId) VALUES (?, ?, ?, ?)", word, definition, UserID, setId)
	return err
}
//update
func updateWordHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	enableCors(&w)

	if r.Method != "PUT" {
		http.Error(w, "Only PUT method is allowed", http.StatusMethodNotAllowed)
		return
	}
	var request struct {
        OriginalWord string `json:"originalWord"`
        Word         string `json:"word"`
        Definition   string `json:"definition"`
        UserID       int    `json:"UserID"` // Add this line
        SetID        int    `json:"setId"`  // Add this line
    }

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Error reading request", http.StatusBadRequest)
        return 
    }
	
	if request.OriginalWord == "" || request.Word == "" || request.Definition == "" || request.UserID == 0 || request.SetID == 0 {
        http.Error(w, "Required fields cannot be empty", http.StatusBadRequest)
        return
    }


	err := UpdateWordToDB(db, request.OriginalWord, request.Word, request.Definition, request.UserID, request.SetID)
	if err != nil {

		if err == ErrWordDoesNotExist {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return 
	}  
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("word updated successfully")
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
//delete
func deleteWordHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	enableCors(&w)

	if r.Method != "DELETE" {
		http.Error(w, "Only DELETE method is allowed", http.StatusMethodNotAllowed)
		return
	}
	 
	var request struct {  
        Word         string `json:"word"`
        UserID       int    `json:"UserID"` // Add this line
        SetID        int    `json:"setId"`  // Add this line
    }

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Error reading request", http.StatusBadRequest)
        return 
    }

    if request.Word == "" || request.UserID == 0 || request.SetID == 0 {
        http.Error(w, "Required fields cannot be empty", http.StatusBadRequest)
        return
    }
    
	err := DeleteWordFromDB(db, request.Word, request.UserID, request.SetID) 
	
	if err != nil {
		if err == ErrWordDoesNotExist {
			http.Error(w, "Word does not exist", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
    
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode("Word deleted succesfully") // Ensure this is the correct response format
}

func DeleteWordFromDB(db *sql.DB, word string, UserID, setId int) error{
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM flashcard WHERE word = ? AND UserID = ? AND setId = ?)", word, UserID, setId).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return ErrWordDoesNotExist
	}

    _, err = db.Exec("DELETE FROM flashcard WHERE word = ? AND UserID = ? AND setId = ?", word, UserID, setId)
	return err
}


//show
func showWordsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    enableCors(&w)

    if r.Method != "GET" {
        http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse query parameters
    queryValues := r.URL.Query()
    UserID, err := strconv.Atoi(queryValues.Get("UserID"))
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    setId, _ := strconv.Atoi(queryValues.Get("setId")) // setId is optional

    words, err := ShowAllWordsFromDB(db, UserID, setId)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(words)
}

func ShowAllWordsFromDB(db *sql.DB, UserID, setId int) ([]WordDefinition, error) {
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

    var words []WordDefinition
    for rows.Next() {
        var wd WordDefinition
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


//seen
func updateSeenHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    enableCors(&w)

    if r.Method != "POST" {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    var request struct {
        Word   string `json:"word"`
        UserID int    `json:"userId"` // Add this line
        SetID  int    `json:"setId"`  // Add this line
    }

    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        log.Println("Error decoding request:", err)
        http.Error(w, "Error reading request", http.StatusBadRequest)
        return
    }

    if request.Word == "" || request.UserID == 0 {
        http.Error(w, "Word and user ID cannot be empty", http.StatusBadRequest)
        return
    }

    if err := updateSeenInDB(db, request.Word, request.UserID, request.SetID); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode("Word marked as seen successfully")
}

func updateSeenInDB(db *sql.DB, word string, userID, setId int) error {
    var query string
    if setId > 0 {
        // Update seen status for a specific set
        query = "UPDATE flashcard SET seen = true WHERE word = ? AND userID = ? AND setId = ?"
        _, err := db.Exec(query, word, userID, setId)
        return err
    }
    // Update seen status for all sets of the user
    query = "UPDATE flashcard SET seen = true WHERE word = ? AND userID = ?"
    _, err := db.Exec(query, word, userID)
    return err
}

// sets

//create-set 
func createSetHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	enableCors(&w)

	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	var request struct {
        UserID  int    `json:"userID"`
        SetName string `json:"setName"`
    }

	if err := json.NewDecoder(r.Body).Decode(&request); err!= nil {
		log.Println("Error decoding request:", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}

	if request.SetName == "" {
		http.Error(w, "Name cannot be empty", http.StatusBadRequest)
		return
	}

	if err := createSetInDB(db, request.UserID, request.SetName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Set created succesfully")

}
func createSetInDB(db *sql.DB, userID int, setName string) error {
	stmt, err := db.Prepare("INSERT INTO flashcardSets (userID, setName) VALUES (?, ?)")

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, setName)
	return err
}

// fetch-sets

func fetchUserSetsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    enableCors(&w)

    if r.Method != "GET" {
        http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
        return
    }

    // Log the incoming request for debugging
    log.Println("Received fetch sets request")

    // Extract the user ID from the request
    UserID, err := getUserIdFromRequest(r)
    if err != nil {
        // Log the error if any
        log.Printf("Error getting user ID from request: %v\n", err)
        http.Error(w, "Unauthorized - Error decoding token", http.StatusUnauthorized)
        return
    }

    // Log the extracted UserID for confirmation
    log.Printf("Fetched User ID: %d from token\n", UserID)

    // Proceed with fetching the user's sets using the UserID
    userSets, err := fetchUserSetsFromDB(db, UserID)
    if err != nil {
        // Log the database error
        log.Printf("Error fetching user sets from DB: %v\n", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // Log the fetched sets for confirmation
    log.Printf("Fetched sets for User ID %d: %+v\n", UserID, userSets)

    // Respond with the user sets
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(userSets)
}


func fetchUserSetsFromDB(db *sql.DB, UserID int) ([]Set, error) {
    // Log the UserID for which sets are being fetched
    log.Printf("Fetching sets for UserID: %d\n", UserID)

    rows, err := db.Query("SELECT setId, setName FROM flashcardSets WHERE UserID = ?", UserID)
    if err != nil {
        // Log the error if query execution fails
        log.Printf("Error executing query: %v\n", err)
        return nil, err
    }
    defer rows.Close()

    var sets []Set
    for rows.Next() {
        var set Set
        if err := rows.Scan(&set.SetId, &set.SetName); err != nil {
            // Log the error if row scanning fails
            log.Printf("Error scanning row: %v\n", err)
            return nil, err
        }
        sets = append(sets, set)
    }

    if err = rows.Err(); err != nil {
        // Log any errors encountered during iteration over the rows
        log.Printf("Error iterating over rows: %v\n", err)
        return nil, err
    }

    // Log the fetched sets for confirmation
    log.Printf("Fetched sets: %+v\n", sets)

    return sets, nil
}


//daily-reset
func resetSeenDaily(db *sql.DB) error {
	_, err := db.Exec("UPDATE flashcard SET seen = false")
	return err
}

func startDailyResetTask(db *sql.DB) {
	go func() {
		for {
			now:= time.Now()
			nextMidnight := now.AddDate(0, 0, 1).Truncate(24 * time.Hour)
			durationUntilMidnight := time.Until(nextMidnight)

			time.Sleep(durationUntilMidnight)

			if err := resetSeenDaily(db); err != nil {
				log.Println("Error resetting seen status:", err)
			}
		}
	}()
}

//errors
func (e DictionaryErr) Error() string {
	return string(e)
}
const (
	ErrNotFound   = DictionaryErr("could not find the word you were looking for")
	ErrWordExists = DictionaryErr("cannot add word because it already exists")
	ErrWordDoesNotExist = DictionaryErr("cannot update word because it does not exist")
)

type DictionaryErr string
//seeding

// var seedData = []struct {
//     Word       string
//     Definition string
// }{
//     {"こんにちは", "Hello"},
//     {"ありがとう", "Thank you"},
//     {"さようなら", "Goodbye"},
//     // Add more words here...
// }

// func seedDatabase(db *sql.DB) error {
//     for _, entry := range seedData {
//         _, err := db.Exec("INSERT INTO flashcard (word, definition) VALUES (?, ?)", entry.Word, entry.Definition)
//         if err != nil {
//             return fmt.Errorf("failed to seed word '%s': %w", entry.Word, err)
//         }
//     }
//     fmt.Println("Database seeded successfully.")
//     return nil
// }
