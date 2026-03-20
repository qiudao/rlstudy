package env

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type StepRequest struct {
	Action int `json:"action"`
}

type StepResponse struct {
	Reward  float64 `json:"reward"`
	Optimal bool    `json:"optimal"`
}

type InfoResponse struct {
	Arms int `json:"arms"`
}

type Server struct {
	mux    *http.ServeMux
	bandit *Bandit
}

func NewServer(arms int, seed int64) *Server {
	s := &Server{
		mux:    http.NewServeMux(),
		bandit: NewBandit(arms, seed),
	}
	s.mux.HandleFunc("GET /info", s.handleInfo)
	s.mux.HandleFunc("POST /reset", s.handleReset)
	s.mux.HandleFunc("POST /step", s.handleStep)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleInfo(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(InfoResponse{Arms: s.bandit.Arms()})
}

func (s *Server) handleReset(w http.ResponseWriter, r *http.Request) {
	s.bandit.Reset()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status":"ok"}`)
}

func (s *Server) handleStep(w http.ResponseWriter, r *http.Request) {
	var req StepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Action < 0 || req.Action >= s.bandit.Arms() {
		http.Error(w, "invalid action", http.StatusBadRequest)
		return
	}
	reward := s.bandit.Step(req.Action)
	optimal := req.Action == s.bandit.OptimalAction()
	json.NewEncoder(w).Encode(StepResponse{Reward: reward, Optimal: optimal})
}
