package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	IMG_DIR_PATH = "./images"
	DB_DIR_PATH  = "./db"
	DB_PATH      = DB_DIR_PATH + "/image_repo.db"
)

type Server struct {
	db     *gorm.DB
	router *mux.Router

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
	router.Path("/image/{id:[0-9]+}").Queries("userid", "{userid:[0-9]+}").HandlerFunc(s.GetOneImage).Methods("GET")

	router.Path("/search").
		Queries("userid", "{userid:[0-9]+}").
		Queries("text", "{text:[a-bA-B\\s]+}").
		Queries("page", "{page:[0-9]+}").
		HandlerFunc(s.SearchImages).Methods("GET")

	s.router = router

}

func (s *Server) InitializeFileSys(reset bool) {
	if _, err := os.Stat(IMG_DIR_PATH); os.IsNotExist(err) {
		log.Println("Initialize Directory for storing the images")
		err = os.Mkdir(IMG_DIR_PATH, os.ModeDir)
		if err != nil {
			log.Fatal(err)
		}
	} else if reset {
		os.RemoveAll(IMG_DIR_PATH)
		err = os.Mkdir(IMG_DIR_PATH, os.ModeDir)
		if err != nil {
			log.Fatal(err)
		}
	}

	// sql
	if _, err := os.Stat(DB_DIR_PATH); os.IsNotExist(err) {
		log.Println("Initialize Directory for storing the sql dbs")
		err := os.Mkdir(DB_DIR_PATH, os.ModeDir)
		if err != nil {
			log.Fatal(err)
		}
	} else if reset {
		os.RemoveAll(DB_DIR_PATH)
		err := os.Mkdir(DB_DIR_PATH, os.ModeDir)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s *Server) InitializeDB() {
	newSchema := false
	if _, err := os.Stat(DB_PATH); os.IsNotExist(err) {
		newSchema = true
	}
	log.Print("Connected to DB...")
	db, err := gorm.Open(sqlite.Open(DB_PATH), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	s.db = db

	if newSchema {
		log.Print("New DB initialize schemas...")
		err = s.db.AutoMigrate(&Image{})
		if err != nil {
			log.Fatal(err)
		}

		err = s.db.AutoMigrate(&User{})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Start to host at port 8080")

	server := Server{}

	server.InitializeFileSys(false)
	server.InitializeDB()
	server.RequestRouter()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: server.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-done
	log.Print("Server Stopped")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")

}
