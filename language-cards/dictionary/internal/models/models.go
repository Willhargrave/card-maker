package models


import (
	"github.com/dgrijalva/jwt-go"

)

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