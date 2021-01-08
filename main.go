package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

type Image struct {
	Name      string
	OwnerId   string
	Path      string
	IsPrivate bool
}

type Server struct {
	Users []string
	Images []Image
}

func(s *Server) homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func (s *Server) RequestRouter() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", s.homePage)

	router.HandleFunc("/signup", s.SignUp).Methods("PUT")

	router.HandleFunc("/upload", s.AddOneImage).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func (s *Server) InitialUploadDir() {
	if _, err := os.Stat("./images"); os.IsNotExist(err) {
		log.Println("Initialize Directory for storing the images")
		err = os.Mkdir("./images", os.ModeDir)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {

	log.Println("Start to host at port 8080")

	server := Server{}
	server.InitialUploadDir()
	server.RequestRouter()
}