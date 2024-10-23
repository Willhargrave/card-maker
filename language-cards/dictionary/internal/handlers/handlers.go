package handlers

import (
	"database/sql"
	"dictionary/internal/config"
	"dictionary/internal/middleware"
	"dictionary/internal/models"
	"dictionary/internal/repository"
	"dictionary/internal/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
)

func SearchHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    middleware.EnableCors(&w)

    if r.Method != "GET" {
        http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
        return
    }

    word := r.URL.Query().Get("word")

    definition, err := repository.SearchDB(db, word) 
    if err != nil {
		log.Printf("Error searching for word '%s': %v", word, err)
		if err == services.ErrWordDoesNotExist {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(definition)
}

func AddWordHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    middleware.EnableCors(&w)
    

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

    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        log.Println("Error decoding request:", err)
        http.Error(w, "Error reading request", http.StatusBadRequest)
        return 
    }
    
    log.Printf("Received request data: %+v\n", request)

    if request.Word == "" || request.Definition == "" || request.SetID == 0 || request.UserID == 0 {
        http.Error(w, "Required fields cannot be empty", http.StatusBadRequest)
        return
    }
    err := repository.AddWordToDB(db, request.Word, request.Definition, request.SetID, request.UserID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode("Word added successfully")
}

func UpdateWordHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	middleware.EnableCors(&w)

	if r.Method != "PUT" {
		http.Error(w, "Only PUT method is allowed", http.StatusMethodNotAllowed)
		return
	}
	var request struct {
        OriginalWord string `json:"originalWord"`
        Word         string `json:"word"`
        Definition   string `json:"definition"`
        UserID       int    `json:"UserID"` 
        SetID        int    `json:"setId"`  
    }

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Error reading request", http.StatusBadRequest)
        return 
    }
	
	if request.OriginalWord == "" || request.Word == "" || request.Definition == "" || request.UserID == 0 || request.SetID == 0 {
        http.Error(w, "Required fields cannot be empty", http.StatusBadRequest)
        return
    }


	err := repository.UpdateWordToDB(db, request.OriginalWord, request.Word, request.Definition, request.UserID, request.SetID)
	if err != nil {

		if err == services.ErrWordDoesNotExist {
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

func DeleteWordHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	middleware.EnableCors(&w)

	if r.Method != "DELETE" {
		http.Error(w, "Only DELETE method is allowed", http.StatusMethodNotAllowed)
		return
	}
	 
	var request struct {  
        Word         string `json:"word"`
        UserID       int    `json:"UserID"` 
        SetID        int    `json:"setId"`  
    }

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Error reading request", http.StatusBadRequest)
        return 
    }

    if request.Word == "" || request.UserID == 0 || request.SetID == 0 {
        http.Error(w, "Required fields cannot be empty", http.StatusBadRequest)
        return
    }
    
	err := repository.DeleteWordFromDB(db, request.Word, request.UserID, request.SetID) 
	
	if err != nil {
		if err == services.ErrWordDoesNotExist {
			http.Error(w, "Word does not exist", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
    
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode("Word deleted succesfully")
}
func ShowWordsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
middleware.EnableCors(&w)

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

words, err := repository.ShowAllWordsFromDB(db, UserID, setId)
if err != nil {
	http.Error(w, "Internal server error", http.StatusInternalServerError)
	return
}

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(words)
}
func UpdateSeenHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    middleware.EnableCors(&w)

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

    if err := repository.UpdateSeenInDB(db, request.Word, request.UserID, request.SetID); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode("Word marked as seen successfully")
}

func CreateSetHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	middleware.EnableCors(&w)

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

	if err := repository.CreateSetInDB(db, request.UserID, request.SetName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Set created succesfully")

}

func FetchUserSetsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    middleware.EnableCors(&w)

    if r.Method != "GET" {
        http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
        return
    }

    // Log the incoming request for debugging
    log.Println("Received fetch sets request")

    // Extract the user ID from the request
    UserID, err := services.GetUserIdFromRequest(r)
    if err != nil {
        // Log the error if any
        log.Printf("Error getting user ID from request: %v\n", err)
        http.Error(w, "Unauthorized - Error decoding token", http.StatusUnauthorized)
        return
    }

    log.Printf("Fetched User ID: %d from token\n", UserID)


    userSets, err := repository.FetchUserSetsFromDB(db, UserID)
    if err != nil {
        log.Printf("Error fetching user sets from DB: %v\n", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    log.Printf("Fetched sets for User ID %d: %+v\n", UserID, userSets)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(userSets)
}


func RegisterHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    middleware.EnableCors(&w)

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

    hashedPassword, err := services.HashPassword(request.Password)
    if err != nil {
        log.Println("Error hashing password:", err)
        http.Error(w, "Error processing the password", http.StatusInternalServerError)
        return
    }

    exists, err := services.UsernameExists(request.Username, db)
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



func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    middleware.EnableCors(&w)

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


    if services.CheckPasswordHash(request.Password, storedHash) {
        expirationTime := time.Now().Add(30 * time.Minute)

        claims := &models.Claims{
            UserID: userID, // Now fetching from DB
            Username: request.Username,
            StandardClaims: jwt.StandardClaims{
                ExpiresAt: expirationTime.Unix(),
            },
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

        tokenString, err := token.SignedString(config.JwtKey)
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