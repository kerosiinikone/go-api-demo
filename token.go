package main

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

type ItemClaims struct {
	Name string `json:"name"`
	jwt.MapClaims
}

func CreateAndSignToken(name string) (string, error) {
	secret := []byte(os.Getenv("SECRET"))
	claims := &jwt.MapClaims{
		"Name": name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return ss, nil
}

func ParseToken(tokenString string) (*Item, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ItemClaims{}, func(token *jwt.Token) (interface{}, error) {
		secret := []byte(os.Getenv("SECRET"))
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*ItemClaims); ok && token.Valid {
		// Not optimal but will do
		item := &Item{
			Name: claims.Name,
		}
		return item, nil
	} 

	return nil, errors.New("invalid token")
}

