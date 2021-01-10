package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Image struct {
	Name      string
	OwnerId   string
	Path      string
	IsPrivate bool
}

type Server struct {
	db *sql.DB

	Users  []string
	Images []Image
}

func (s *Server) homePage(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) InitializeFilesys() {
	// upload
	if _, err := os.Stat("./images"); os.IsNotExist(err) {
		log.Println("Initialize Directory for storing the images")
		err = os.Mkdir("./images", os.ModeDir)
		if err != nil {
			log.Fatal(err)
		}
	}

	// sql
	if _, err := os.Stat("./db"); os.IsNotExist(err) {
		log.Println("Initialize Directory for storing the sql dbs")
	} else {
		os.RemoveAll("./db")
	}

	err := os.Mkdir("./db", os.ModeDir)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) InitializeDB() {
	_, err := s.db.Exec(CreateProfileTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = s.db.Exec(CreatImageTable)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {

	log.Println("Start to host at port 8080")

	server := Server{}

	server.InitializeFilesys()

	db, err := sql.Open("sqlite3", "./db/image_repo.db")
	if err != nil {
		log.Fatal(err)
	}
	server.db = db

	defer db.Close()

	server.InitializeDB()

	server.RequestRouter()
}
