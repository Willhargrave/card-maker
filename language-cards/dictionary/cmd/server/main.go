package main

import (
	"fmt"
	"net/http"
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3" 
	"dictionary/internal/database"
	"dictionary/internal/middleware"
	"dictionary/internal/handlers"
	"dictionary/internal/services"

)

func main() {
    db, err := sql.Open("sqlite3", "./dictionary.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
	database.CreateUserTable(db)
	database.CreateFlashcardTable(db)
	database.CreateFlashcardSetsTable(db)



    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Welcome")
    })

 
    http.HandleFunc("/add-word", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			middleware.EnableCors(&w)
			w.WriteHeader(http.StatusOK)
			return
		}
        handlers.AddWordHandler(w, r, db)
    })

    http.HandleFunc("/update-word", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			middleware.EnableCors(&w)
			w.WriteHeader(http.StatusOK)
			return
		}
        handlers.UpdateWordHandler(w, r, db)
    })
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			middleware.EnableCors(&w)
			w.WriteHeader(http.StatusOK)
			return
		}
		handlers.SearchHandler(w, r, db)
	})
	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			middleware.EnableCors(&w)
			w.WriteHeader(http.StatusOK)
			return
		}
		handlers.DeleteWordHandler(w, r, db)
	})
	http.HandleFunc("/show-words", func(w http.ResponseWriter, r *http.Request) {
		handlers.ShowWordsHandler(w, r, db)
	})

	http.HandleFunc("/seen", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			middleware.EnableCors(&w)
			w.WriteHeader(http.StatusOK)
			return
		}
		handlers.UpdateSeenHandler(w, r, db)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			middleware.EnableCors(&w)
			w.WriteHeader(http.StatusOK)
			return
		}
		handlers.RegisterHandler(w, r, db)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			middleware.EnableCors(&w)
			w.WriteHeader(http.StatusOK)
			return
		}
		handlers.LoginHandler(w, r, db)
	})
    
	http.HandleFunc("/create-set", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			middleware.EnableCors(&w)
			w.WriteHeader(http.StatusOK)
			return
		}
		handlers.CreateSetHandler(w, r, db)
	})

	http.HandleFunc("/fetch-sets", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			middleware.EnableCors(&w)
			w.WriteHeader(http.StatusOK)
			return
		}
		handlers.FetchUserSetsHandler(w, r, db)
	})
	


	// if err := database.SeedDatabase(db); err != nil {
    //     log.Fatal("Error seeding the database:", err)
    // }

	services.StartDailyResetTask(db)
	
    fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("Error starting server: ", err)
    }
}