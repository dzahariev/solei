package controller

import (
	"net/http"
)

// Home is an API root route controller
func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	server.JSON(w, http.StatusOK, "Welcome!")
}
