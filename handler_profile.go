package main

import (
	"encoding/json"
	"errors"
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
		http.Error(w, "400 bad request: missing required field or invalid json", http.StatusBadRequest)
		return
	}
	newUser := User{
		Name:      signUp.Name,
		CreatedAt: time.Now(),
	}

	result := s.db.Where("name = ?", signUp.Name).First(&newUser)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			result = s.db.Create(&newUser)
			if result.Error != nil {
				log.Println(result.Error)
				http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
				return
			}
			log.Printf("[SignUp]: New User id: %v name: %v\n", newUser.ID, signUp.Name)

			reply := SignUpReply{ID: newUser.ID}

			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(reply)
			return
		} else {
			log.Println(result.Error)
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		log.Printf("[SignUp]name: %v used already. ID %v\n", signUp.Name, newUser.ID)
		http.Error(w, "403: input name used alread", http.StatusForbidden)
		return
	}
}
