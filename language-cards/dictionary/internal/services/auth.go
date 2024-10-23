package services

import (
	"database/sql"
	"dictionary/internal/config"
	"dictionary/internal/models"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func GetUserIdFromRequest(r *http.Request) (int, error) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        return 0, errors.New("authorization header is missing")
    }

    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        return 0, errors.New("authorization header format must be 'Bearer {token}'")
    }

    tknStr := parts[1]
    claims := &models.Claims{}

    tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
        return config.JwtKey, nil
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

func UsernameExists(username string, db *sql.DB) (bool, error) {
    var exists bool
    err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exists)
    if err != nil {
        log.Println("Error checking if username exists:", err)
        return false, err 
    }
    return exists, nil
}

func CheckPasswordHash(password, hashedPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}

// custom errors
func (e DictionaryErr) Error() string {
	return string(e)
}
const (
	ErrNotFound   = DictionaryErr("could not find the word you were looking for")
	ErrWordExists = DictionaryErr("cannot add word because it already exists")
	ErrWordDoesNotExist = DictionaryErr("cannot update word because it does not exist")
)

type DictionaryErr string