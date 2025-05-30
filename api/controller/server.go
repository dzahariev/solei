package controller

import (
	"fmt"
	"log"
	"net/http"

	"gopkg.in/yaml.v3"

	"github.com/dzahariev/solei/api/model"
	"github.com/dzahariev/solei/api/security"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Server represent current API server
type Server struct {
	DB                *gorm.DB
	Router            *mux.Router
	AuthClient        *security.AuthClient
	ResourceFactory   *ResourceFactory
	RoleToPermissions map[string][]string
}

// DBInitialize is used to init a DB cnnection
func (server *Server) DBInitialize(dbUser, dbPassword, dbPort, dbHost, dbName string) {
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, dbUser, dbName, dbPassword)
	var err error
	server.DB, err = gorm.Open(postgres.Open(DBURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Cannot connect to database with error: %v", err)
	}
	log.Printf("We are connected to the database")

	// server.DB.AutoMigrate(&model.User{}, &model.Event{}, &model.Session{}, &model.Subscription{}, &model.Comment{})
}

// AuthInitialize is used to register routes
func (server *Server) AuthInitialize(authURL, authRealm, authClientID, authClientSecret string) {
	server.AuthClient = &security.AuthClient{}
	server.AuthClient.Initialize(authURL, authRealm, authClientID, authClientSecret)
}

// ResourcesInitialize is used to register resources
func (server *Server) ResourcesInitialize() {
	server.ResourceFactory = NewResourceFactory()
	// Register all resources here
	server.ResourceFactory.Register(&model.Category{})
	server.ResourceFactory.Register(&model.Meal{})
	server.ResourceFactory.Register(&model.Order{})
	server.ResourceFactory.Register(&model.OrderItem{})
	server.ResourceFactory.Register(&model.Address{})
	server.ResourceFactory.Register(&model.User{})
}

// RoutesInitialize is used to register routes
func (server *Server) RoutesInitialize() {
	server.Router = mux.NewRouter()

	// Unsecured Home Route
	server.Router.HandleFunc("/api/", server.Public(server.ContentTypeJSON(server.Home))).Methods(http.MethodGet)
	// Register all resource routes
	for _, resource := range server.ResourceFactory.Resources() {
		server.Router.HandleFunc(fmt.Sprintf("/api/%s", resource.Name), server.Protected(server.ContentTypeJSON(server.Create(resource.Name)), resource, WRITE)).Methods(http.MethodPost)
		server.Router.HandleFunc(fmt.Sprintf("/api/%s", resource.Name), server.Protected(server.ContentTypeJSON(server.GetAll(resource.Name)), resource, READ)).Methods(http.MethodGet)
		server.Router.HandleFunc(fmt.Sprintf("/api/%s/{id}", resource.Name), server.Protected(server.ContentTypeJSON(server.Get(resource.Name)), resource, READ)).Methods(http.MethodGet)
		server.Router.HandleFunc(fmt.Sprintf("/api/%s/{id}", resource.Name), server.Protected(server.ContentTypeJSON(server.Update(resource.Name)), resource, WRITE)).Methods(http.MethodPut)
		server.Router.HandleFunc(fmt.Sprintf("/api/%s/{id}", resource.Name), server.Protected(server.ContentTypeJSON(server.Delete(resource.Name)), resource, WRITE)).Methods(http.MethodDelete)
	}
	// Static Route
	server.Router.PathPrefix("/").Handler(server.Static())
}

// Initialize is used to init a DB cnnection and register routes
func (server *Server) Initialize(dbUser, dbPassword, dbPort, dbHost, dbName, authURL, authRealm, authClientID, authClientSecret, roleToPermissionsYaml string) {
	server.DBInitialize(dbUser, dbPassword, dbPort, dbHost, dbName)
	server.AuthInitialize(authURL, authRealm, authClientID, authClientSecret)
	server.RolesInitialize(roleToPermissionsYaml)
	server.ResourcesInitialize()
	server.RoutesInitialize()
}

func (server *Server) RolesInitialize(roleToPermissionsYaml string) {
	err := yaml.Unmarshal([]byte(roleToPermissionsYaml), &server.RoleToPermissions)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot parse role to permissions yaml: %w", err))
	}
	log.Println("Roles and permissions initialized")
	for role, permissions := range server.RoleToPermissions {
		log.Printf("%s: %v\n", role, permissions)
	}

}

// Run starts the http server
func (server *Server) Run(addr string) {
	log.Println("Listening to port 8800")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

// Run starts the http server
func (server *Server) ScopedDB(request *http.Request) *gorm.DB {
	dbScopes := NewDBScopes(request)
	return server.DB.Scopes(dbScopes.Owned(), dbScopes.Paginate())
}
