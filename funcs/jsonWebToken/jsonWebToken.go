package jsonWebToken

import (
	"errors"
	"time"
	"log"

	"github.com/dgrijalva/jwt-go"
	"app/funcs/errores"
	"app/funcs/structs"
)

type JWTClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var SecretKey = []byte("Interfell2023")


func GenerateToken(username string) (string, error) {
	
	TToken := time.Duration(structs.DatosInit.JsonTokenTime ) * time.Hour
	claims := JWTClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TToken).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   "auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(SecretKey)
	errores.CheckErr(err)

	return tokenString, nil
}


func ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	errores.CheckErr(err)
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		log.Println(claims.Username)
		return true, nil
	}
	return false, errors.New("Token inválido")
}
