package middlewares

import (
	"errors"
	"net/http"

	"github.com/nuoinguyen/gin-gonic/app/auth"
	"github.com/nuoinguyen/gin-gonic/app/repositories"
)

// SetMiddlewareJSON is function format data return is a JSON
func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

// SetMiddlewareAuthentication is function format data return is a JSON
func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			repositories.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next(w, r)
	}
}
