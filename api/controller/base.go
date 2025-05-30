package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/dzahariev/solei/api/model"
	"github.com/gofrs/uuid/v5"
	"github.com/gorilla/mux"
)

// GetAll retrieves all objects
func (server *Server) GetAll(resourceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		object, err := server.ResourceFactory.New(resourceName)
		if err != nil {
			server.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		count, err := object.Count(server.ScopedDB(r), object)
		if err != nil {
			server.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		data, err := object.FindAll(server.ScopedDB(r), object)
		if err != nil {
			server.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		list := model.List{
			Count:    count,
			PageSize: GetPageSize(r),
			Page:     GetPage(r),
			Data:     *data,
		}

		server.JSON(w, http.StatusOK, list)
	}
}

// Get loads an object by given ID
func (server *Server) Get(resourceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uid, err := uuid.FromString(vars["id"])
		if err != nil {
			server.ERROR(w, http.StatusBadRequest, err)
			return
		}

		object, err := server.ResourceFactory.New(resourceName)
		if err != nil {
			server.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		err = object.FindByID(server.ScopedDB(r), object, uid)
		if err != nil {
			server.ERROR(w, http.StatusNotFound, err)
			return
		}
		server.JSON(w, http.StatusOK, object)
	}
}

// Create is caled to create an object
func (server *Server) Create(resourceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			server.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		object, err := server.ResourceFactory.New(resourceName)
		if err != nil {
			server.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		err = json.Unmarshal(body, object)
		if err != nil {
			server.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		err = object.Validate()
		if err != nil {
			server.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		if !IsGlobal(r) {
			usetUUID, err := uuid.FromString(GetOwnerID(r))
			if err != nil {
				server.ERROR(w, http.StatusInternalServerError, err)
				return
			}
			objectAsLocalObject := object.(model.LocalObject)
			objectAsLocalObject.SetUserID(usetUUID)
		}

		err = object.Save(server.ScopedDB(r), object)

		if err != nil {
			server.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, object.GetID()))
		server.JSON(w, http.StatusCreated, object)
	}
}

// UpdateBook updates existing object
func (server *Server) Update(resourceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uid, err := uuid.FromString(vars["id"])
		if err != nil {
			server.ERROR(w, http.StatusBadRequest, err)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			server.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		object, err := server.ResourceFactory.New(resourceName)
		if err != nil {
			server.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		err = json.Unmarshal(body, &object)
		if err != nil {
			server.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		err = object.Validate()
		if err != nil {
			server.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		recordExisting := reflect.New(reflect.TypeOf(object).Elem()).Interface().(model.Object)
		err = recordExisting.FindByID(server.ScopedDB(r), recordExisting, uid)
		if err != nil {
			server.ERROR(w, http.StatusNotFound, err)
			return
		}

		object.SetID(uid)

		err = object.Update(server.ScopedDB(r), object)
		if err != nil {
			server.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		server.JSON(w, http.StatusOK, object)
	}
}

// Delete deletes an objec t
func (server *Server) Delete(resourceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		uid, err := uuid.FromString(vars["id"])
		if err != nil {
			server.ERROR(w, http.StatusBadRequest, err)
			return
		}

		object, err := server.ResourceFactory.New(resourceName)
		if err != nil {
			server.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		err = object.FindByID(server.ScopedDB(r), object, uid)
		if err != nil {
			server.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		err = object.Delete(server.ScopedDB(r), object)
		if err != nil {
			server.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Entity", fmt.Sprintf("%s", uid))
		server.JSON(w, http.StatusNoContent, "")
	}
}
