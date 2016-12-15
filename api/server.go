package main

import (
	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
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

// Close closes down the server
func (s *Server) Close() {
	s.db.Close()
}
