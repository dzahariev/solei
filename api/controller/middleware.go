package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Define a custom type for context keys
type contextKey string

const (
	READ  = "read"
	WRITE = "write"

	GLOBAL = "global"

	CURRENT_USER_ID contextKey = "currentUserID"
	GLOBAL_SCOPE    contextKey = "globalScope"
)

// Wrapper for static resources
func (server *Server) Static() http.Handler {
	return http.FileServer(http.Dir("./public"))
}

// Wrapper for public resources
func (server *Server) Public(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
}

// Wrapper for protected and Global resources
func (server *Server) Protected(next http.HandlerFunc, resource Resource, permission string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse token
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) < 7 {
			server.ERROR(w, http.StatusUnauthorized, fmt.Errorf("unauthorized, missing bearer authorization header"))
			return
		}
		authType := strings.ToLower(authHeader[:6])
		if authType != "bearer" {
			server.ERROR(w, http.StatusUnauthorized, fmt.Errorf("unauthorized, invalid bearer authorization header"))
			return
		}

		// Verify token is valid
		tokenString := authHeader[7:]
		tokenString = strings.TrimSpace(tokenString)
		err := server.AuthClient.RetrospectToken(r.Context(), tokenString)
		if err != nil {
			server.ERROR(w, http.StatusUnauthorized, err)
			return
		}

		// Create user if not exists
		userFromInfo, err := server.AuthClient.GetUserFromToken(r.Context(), tokenString)
		if err != nil {
			server.ERROR(w, http.StatusUnauthorized, err)
			return
		}

		loadedUser, _ := server.DBLoadUser(string(userFromInfo.ID.String())) // we ignore the error as it is expected if user do not exists
		if loadedUser == nil {
			err := server.DBSaveUser(userFromInfo)
			if err != nil {
				server.ERROR(w, http.StatusUnauthorized, err)
				return
			}
		}

		loadedUser, err = server.DBLoadUser(string(userFromInfo.ID.String()))
		if err != nil {
			server.ERROR(w, http.StatusUnauthorized, err)
			return
		}

		// Create new context with current user
		newCtx := context.WithValue(r.Context(), CURRENT_USER_ID, loadedUser.ID.String())

		// Get roles from token
		roles, err := server.AuthClient.GetRolesFromToken(r.Context(), tokenString)
		if err != nil {
			server.ERROR(w, http.StatusUnauthorized, err)
			return
		}

		var permissions []string
		for _, role := range roles {
			permissions = append(permissions, server.RoleToPermissions[role]...)
		}

		// Check for global scope for this resource and if exists add it to the enriched context
		if resource.IsGlobal {
			newCtx = context.WithValue(newCtx, GLOBAL_SCOPE, GLOBAL)
		}

		// Replace request context
		rWithUpdatedContext := r.WithContext(newCtx)

		// Check permissions
		if server.havePermission(resource.Name, permission, permissions) {
			next(w, rWithUpdatedContext)
		} else {
			// lack of permissions
			server.ERROR(w, http.StatusUnauthorized, fmt.Errorf("unauthorized, no permission for %s.%s", resource.Name, permission))
			return
		}
	}
}

// ContentTypeJSON set the content type to JSON
func (server *Server) ContentTypeJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

// JSON returns data as JSON stream
func (server *Server) JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

// ERROR returns error as JSON representation
func (server *Server) ERROR(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		server.JSON(w, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	server.JSON(w, http.StatusBadRequest, nil)
}

func (server *Server) havePermission(resource, permission string, permissions []string) bool {
	for _, currentPermission := range permissions {
		resourcePermission := fmt.Sprintf("%s.%s", resource, permission)
		if strings.EqualFold(currentPermission, resourcePermission) {
			return true
		}
	}
	return false
}
