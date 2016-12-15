package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Domain represents a domain in the database
type Domain struct {
	Domain string `json:"domain"`
	Rank   int    `json:"rank"`
}

// GetDomainHandler returns the current details for a given domain
func (s *Server) GetDomainHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		s.log.WithField("func", "GetDomainHandler")
		s.log.Infoln(r.Method, r.URL.Path, r.RemoteAddr)

		vars := mux.Vars(r)
		domain := vars["domain"]

		var domainRes Domain
		if err = s.db.Get(&domainRes, "SELECT domain, rank FROM current WHERE domain=?", domain); err != nil {
			httpResponse(w, &errorResponse{Error: "Domain not found"}, http.StatusNotFound)
			return
		}
		httpResponse(w, &domainRes, http.StatusOK)
	})
}

// Rank represents a given rank in the database
type Rank struct {
	Domain string `json:"domain"`
	Rank   int    `json:"rank"`
}

// GetRankHandler returns the domain at the given rank
func (s *Server) GetRankHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		s.log.WithField("func", "GetRankHandler")
		s.log.Infoln(r.Method, r.URL.Path, r.RemoteAddr)

		vars := mux.Vars(r)
		rank, _ := strconv.Atoi(vars["rank"]) // Gorilla Mux guarantees this is an integer
		if rank < 0 || rank > 1000000 {
			httpResponse(w, &errorResponse{Error: "Rank must be between 1 and 100000"}, http.StatusBadRequest)
			return
		}

		var rankRes Rank
		if err = s.db.Get(&rankRes, "SELECT domain, rank FROM current WHERE rank=?", rank); err != nil {
			s.log.Errorln(err)
			httpResponse(w, &errorResponse{Error: "Rank not found"}, http.StatusNotFound)
			return
		}
		httpResponse(w, &rankRes, http.StatusOK)
	})
}
