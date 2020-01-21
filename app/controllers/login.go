package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/nuoinguyen/gin-gonic/app/auth"
	"github.com/nuoinguyen/gin-gonic/app/models"
	"github.com/nuoinguyen/gin-gonic/app/repositories"
	"github.com/nuoinguyen/gin-gonic/app/utilities"
	"golang.org/x/crypto/bcrypt"
)

// Login is a function login account
func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		repositories.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		repositories.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		repositories.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := utilities.FormatError(err.Error())
		repositories.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	repositories.JSON(w, http.StatusOK, token)
}

// SignIn is func login account of an user
func (server *Server) SignIn(email, password string) (string, error) {

	var err error

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.ID)
}
