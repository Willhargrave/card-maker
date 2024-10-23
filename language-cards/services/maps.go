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

// Tables

//cors



//user-login




//functions


//main

//search







//daily-reset


//errors

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
