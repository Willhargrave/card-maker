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



    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Welcome")
    })

 
    http.HandleFunc("/add-word", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			enableCors(&w)
			w.WriteHeader(http.StatusOK)
			return
		}
        addWordHandler(w, r, db)
    })

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