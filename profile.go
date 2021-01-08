package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type SignUpRequest struct {
	UserId string `json:"UserId"`
}

func (s *Server) SignUp(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var signUp SignUpRequest
	err := json.Unmarshal(reqBody, &signUp)
	log.Printf( "[SignUp] request: %v\n", string(reqBody))

	if err != nil {
		fmt.Fprintf(w, "401: bad request")
		w.WriteHeader(401)
		return
	}

	s.Users = append(s.Users, signUp.UserId)
	log.Printf("[SignUp]: New User %v\n", signUp.UserId)

	w.WriteHeader(200)
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) LogOut(w http.ResponseWriter, r *http.Request) {

}
