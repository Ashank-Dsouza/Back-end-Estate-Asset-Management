package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"os"

	"bitbucket.org/staydigital/truvest-identity-management/api/auth"
	"bitbucket.org/staydigital/truvest-identity-management/api/models"

	"bitbucket.org/staydigital/truvest-identity-management/api/responses"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"encoding/json"
	"io/ioutil"
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

func (server *Server) ExtractUserId(r *http.Request) uuid.UUID {
	tokenString := auth.ExtractToken(r)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		fmt.Println(err, " is the error that was happened!")
		return uuid.Nil
	}
	fmt.Println(token, " is the token that was happened!")

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		u_id, err := uuid.Parse(fmt.Sprint(claims["user_id"]))
		if err != nil {
			return uuid.Nil
		}
		return u_id
	}
	return uuid.Nil
}

func (server *Server) GetTokenId(w http.ResponseWriter, r *http.Request) uuid.UUID {
	authID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return uuid.Nil
	}
	tokenID, err := uuid.Parse(authID)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return uuid.Nil
	}
	return tokenID
}

// ConfirmEmail godoc
// @Summary Confirm a user's email
// @Description User can confirm their email with confirmation token.
// @Tags ConfirmEmail
// @Accept  json
// @Produce  json
// @Param logout body ConfirmEmail true "ConfirmEmail"
// @Success 200 {object} string
// @Security ApiKeyAuth
// @Router /confirmEmail [put]
func (server *Server) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside ConfirmEmail()")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	confirm_email := models.Confirm_email{}
	err = json.Unmarshal(body, &confirm_email)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	updatedEmail, err := confirm_email.ConfirmAUserEmail(server.DB)

	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	responses.JSON(w, http.StatusOK, updatedEmail)

}
