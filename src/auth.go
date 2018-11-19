package main

import (
    "errors"
    "fmt"
    "time"

    "golang.org/x/crypto/bcrypt"
    "github.com/dgrijalva/jwt-go"
)

func verify_user_pass(form_username string, form_password string) bool {
    user, user_err := GetUser(form_username)

    if user_err != nil {
        return false
    }

    incoming := []byte(form_password)
    existing := []byte(user.PasswordHash)
    err := bcrypt.CompareHashAndPassword(existing, incoming)

    return err == nil
}

func CreateToken(form_username string, form_password string) (string, error) {
    if !verify_user_pass(form_username, form_password) {
        return "", errors.New("Login credentials are invalid")
    }

    /* Create the token */
    token := jwt.New(jwt.SigningMethodHS256)

    /* Create a map to store our claims */
    claims := token.Claims.(jwt.MapClaims)

    /* Set token claims */
    claims["admin"] = true
    claims["name"] = form_username
    claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

    /* Sign the token with our secret */
    my_signing_key := []byte("secret")
    token_string, _ := token.SignedString(my_signing_key)

    /* Finally, return the token string */
    return token_string, nil
}

func VerifyToken(token_string string) bool {
    token, err := jwt.Parse(token_string, func(token *jwt.Token) (interface{}, error) {
        // Don't forget to validate the alg is what you expect:
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }

        // hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
        return []byte("secret"), nil
    })

    if err != nil {
        return false
    }

    // if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
    // }
    return token.Valid
}
