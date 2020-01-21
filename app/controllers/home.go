package controllers

import (
	"net/http"

	"github.com/nuoinguyen/gin-gonic/app/repositories"
)

// Server variable call controller Server
// var Server = routes.Server{}

// Home is default
func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	repositories.JSON(w, http.StatusOK, "Welcome To This Awesome API")

}
