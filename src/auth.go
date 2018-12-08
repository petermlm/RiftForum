package main

import (
    "errors"
    "fmt"
    "time"

    "golang.org/x/crypto/bcrypt"
    "github.com/dgrijalva/jwt-go"
)

type Claims struct {
    jwt.StandardClaims
    Username string
    Usertype UserTypes
}

func verify_user_pass(user *User, password string) bool {
    incoming := []byte(password)
    existing := []byte(user.PasswordHash)
    err := bcrypt.CompareHashAndPassword(existing, incoming)
    return err == nil
}

func CreateToken(form_username string, form_password string) (string, error) {
    user, user_err := GetUser(form_username)

    if user_err != nil {
        return "", errors.New("Login credentials are invalid")
    }

    if !verify_user_pass(user, form_password) {
        return "", errors.New("Login credentials are invalid")
    }

    claims := Claims {
        jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
        },
        user.Username,
        user.Usertype,
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    my_signing_key := []byte("secret")
    token_string, _ := token.SignedString(my_signing_key)

    return token_string, nil
}

func VerifyToken(token_string string) *Claims {
    token, err := jwt.ParseWithClaims(token_string, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }

        return []byte("secret"), nil
    })

    if err != nil || !token.Valid {
        return nil
    }

    return token.Claims.(*Claims)
}