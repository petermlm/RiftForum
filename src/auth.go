package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var secret []byte

type Claims struct {
	jwt.StandardClaims
	Id       uint
	Username string
	UserType UserTypes
}

func InitAuth() {
	if Config.DebugMode {
		secret = []byte("secret")
	} else {
		var err error
		secret, err = ioutil.ReadFile(Config.SecretFilename)
		if err != nil {
			panic("Could not read file with secret.")
		}
	}

	log.Println("Authentication initialized")
}

func CreateToken(form_username string, form_password string) (string, error) {
	user, user_err := GetUser(form_username)

	if user_err != nil {
		return "", errors.New("Login credentials are invalid")
	}

	if user.Banned {
		return "", errors.New("")
	}

	if !VerifyUserPass(user, form_password) {
		return "", errors.New("Login credentials are invalid")
	}

	claims := Claims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
		user.Id,
		user.Username,
		user.UserType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token_string, _ := token.SignedString(secret)
	return token_string, nil
}

func VerifyToken(token_string string) *Claims {
	token, err := jwt.ParseWithClaims(token_string, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if err != nil || !token.Valid {
		return nil
	}

	return token.Claims.(*Claims)
}

func VerifyUserPass(user *User, password string) bool {
	incoming := []byte(password)
	existing := []byte(user.PasswordHash)
	err := bcrypt.CompareHashAndPassword(existing, incoming)
	return err == nil
}
