package controllers

import (
	"fmt"
	"net/http"

	"os"

	"bitbucket.org/staydigital/truvest-identity-management/api/auth"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Is Call a self-edit check
func (server *Server) IsSelfEditRequest(r *http.Request) bool {
	vars := mux.Vars(r)
	uid, err := uuid.Parse(vars["id"])
	fmt.Println(uid, " is the user id from the request!")

	if err != nil {
		fmt.Println(err, " is the error that was happened!")
		return false
	}

	tokenString := auth.ExtractToken(r)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		fmt.Println(err, " is the error that was happened!")
		return false
	}
	fmt.Println(token, " is the token that was happened!")

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		u_id, err := uuid.Parse(fmt.Sprint(claims["user_id"]))
		if err != nil {
			return false
		}
		fmt.Println(u_id, " is the user id from the jwt!")

		if uid == u_id {
			return true
		} else {
			return false
		}
	}
	return false
}
