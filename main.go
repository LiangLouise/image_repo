package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Server struct {
	db *gorm.DB

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

	router.HandleFunc("/signup", s.SignUp).Methods("POST")

	router.HandleFunc("/image", s.AddOneImage).Methods("POST")
	router.HandleFunc("/image", s.DeleteImage).Methods("DELETE")

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
	err := s.db.AutoMigrate(&Image{})
	if err != nil {
		log.Fatal(err)
	}

	err = s.db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}

}

func main() {

	log.Println("Start to host at port 8080")

	server := Server{}

	server.InitializeFilesys()

	db, err := gorm.Open(sqlite.Open("./db/image_repo.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	server.db = db

	// defer db.Close()

	server.InitializeDB()
	server.RequestRouter()
}
