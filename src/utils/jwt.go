package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var JWTSecret = []byte("this is my !secret!key!")

func GenerateJWT(id uint, username string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, _ := token.SignedString(JWTSecret)
	return t
}
