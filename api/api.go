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

// GetDomain returns details of a domain
func (s *Server) GetDomain(domain string) (Domain, error) {
	var d Domain
	err := s.db.Get(&d, "SELECT domain, rank FROM current WHERE domain=?", domain)
	return d, err

}

// GetDomainHandler returns the current details for a given domain
func (s *Server) GetDomainHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		s.log.WithField("func", "GetDomainHandler")
		s.log.Infoln(r.Method, r.URL.Path, r.RemoteAddr)

		vars := mux.Vars(r)
		domain := vars["domain"]

		domainRes, err := s.GetDomain(domain)
		if err != nil {
			s.log.Errorln(err)
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

// GetRank returns the details of the domain at the given rank
func (s *Server) GetRank(rank int) (Rank, error) {
	var r Rank
	err := s.db.Get(&r, "SELECT domain, rank FROM current WHERE rank=?", rank)
	return r, err
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

		rankRes, err := s.GetRank(rank)
		if err != nil {
			s.log.Errorln(err)
			httpResponse(w, &errorResponse{Error: "Rank not found"}, http.StatusNotFound)
			return
		}
		httpResponse(w, &rankRes, http.StatusOK)
	})
}
