package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// Server represents our API server
type Server struct {
	db  *sqlx.DB
	log *logrus.Logger
}

// NewServer creates a new server :)
func NewServer() (*Server, error) {
	return &Server{log: logrus.New()}, nil
}

// ConnectDB connects our server to the given DB
func (s *Server) ConnectDB(db string) error {
	var err error
	s.db, err = sqlx.Connect("mysql", db)
	return err
}

// Router returns and HTTP router with the handlers for this server
func (s *Server) Router() *mux.Router {
	r := mux.NewRouter()
	r.Handle("/domain/{domain:[a-zA-Z0-9.-]+}", s.GetDomainHandler()).Methods("GET")
	r.Handle("/rank/{rank:[0-9]+}", s.GetRankHandler()).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("build")))
	return r
}

// Close closes down the server
func (s *Server) Close() error {
	return s.db.Close()
}
