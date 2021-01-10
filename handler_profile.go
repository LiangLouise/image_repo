package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type SignUpRequest struct {
	Name string `json:"name"`
}

type SignUpReply struct {
	ID uint `json:"id"`
}

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}

func (s *Server) SignUp(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var signUp SignUpRequest
	err := json.Unmarshal(reqBody, &signUp)
	log.Printf("[SignUp] request: %v\n", string(reqBody))

	if err != nil {
		w.WriteHeader(401)
		fmt.Fprintf(w, "401: bad request")
		return
	}
	newUser := User{
		Name:      signUp.Name,
		CreatedAt: time.Now(),
	}

	result := s.db.Where("name = ?", signUp.Name).First(&newUser)

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println(result.Error)
		w.WriteHeader(500)
		fmt.Fprintf(w, "500 Internal Server Error")
		return
	} else {
		log.Printf("[SignUp]name: %v used already. ID %v\n", signUp.Name, newUser.ID)
		w.WriteHeader(403)
		fmt.Fprintf(w, "401:name: %v used already", signUp.Name)
		return
	}

	result = s.db.Create(&newUser)
	if result.Error != nil {
		log.Println(result.Error)
		w.WriteHeader(500)
		fmt.Fprintf(w, "500 Internal Server Error")
		return
	}
	log.Printf("[SignUp]: New User id: %v name: %v\n", newUser.ID, signUp.Name)

	reply := SignUpReply{ID: newUser.ID}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reply)
}
