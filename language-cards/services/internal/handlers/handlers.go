package handlers 

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
func showWordsHandler {
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