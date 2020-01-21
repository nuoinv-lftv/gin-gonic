package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nuoinguyen/gin-gonic/app/auth"
	"github.com/nuoinguyen/gin-gonic/app/models"
	"github.com/nuoinguyen/gin-gonic/app/repositories"
	"github.com/nuoinguyen/gin-gonic/app/utilities"
)

// Server variable call controller Server
// var Server = routes.Server{}

// CreateUser is a function create user account
func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		repositories.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		repositories.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user.Prepare()
	err = user.Validate("")
	if err != nil {
		repositories.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	userCreated, err := user.SaveUser(server.DB)

	if err != nil {

		formattedError := utilities.FormatError(err.Error())

		repositories.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	repositories.JSON(w, http.StatusCreated, userCreated)
}

// GetUsers is func get list user account
func (server *Server) GetUsers(w http.ResponseWriter, r *http.Request) {

	user := models.User{}

	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		repositories.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	repositories.JSON(w, http.StatusOK, users)
}

// GetUser is function get info an user
func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		repositories.ERROR(w, http.StatusBadRequest, err)
		return
	}
	user := models.User{}
	userGotten, err := user.FindUserByID(server.DB, uint32(uid))
	if err != nil {
		repositories.ERROR(w, http.StatusBadRequest, err)
		return
	}
	repositories.JSON(w, http.StatusOK, userGotten)
}

// UpdateUser is function Update info user account
func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		repositories.ERROR(w, http.StatusBadRequest, err)
		return
	}
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
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		repositories.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != uint32(uid) {
		repositories.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		repositories.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	updatedUser, err := user.UpdateAUser(server.DB, uint32(uid))
	if err != nil {
		formattedError := utilities.FormatError(err.Error())
		repositories.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	repositories.JSON(w, http.StatusOK, updatedUser)
}

// DeleteUser is func remove user account in system
func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	user := models.User{}

	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		repositories.ERROR(w, http.StatusBadRequest, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		repositories.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != 0 && tokenID != uint32(uid) {
		repositories.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	_, err = user.DeleteAUser(server.DB, uint32(uid))
	if err != nil {
		repositories.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	repositories.JSON(w, http.StatusNoContent, "")
}
