package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	REPO_PATH = "./images"
)

type Image struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	OwnerId   uint      `gorm:"not null"`
	Path      string    `gorm:"not null"`
	IsPrivate bool      `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}

func saveImage() {

}

func downloadImageByUrl(url string) {

}

// FormData
// 	- userId
// 	- isPrivate
//	- file
func (s *Server) AddOneImage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		panic(err)
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Fprintf(w, "401 bad request: not a valid file")
		w.WriteHeader(401)
		return
	}
	defer file.Close()

	imageName := r.FormValue("imageName")
	if imageName == "" {
		fmt.Fprintf(w, "401 bad request: require imageName")
		w.WriteHeader(401)
		return
	}

	isPrivate := r.FormValue("isPrivate")
	isPrivateB, err := strconv.ParseBool(isPrivate)
	if err != nil {
		fmt.Fprintf(w, "401 bad request: isPrivate needs to be a bool value")
		w.WriteHeader(401)
		return
	}

	userId := r.FormValue("userId")
	userIdInt, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		fmt.Fprintf(w, "401 bad request: require the userId")
		w.WriteHeader(401)
		return
	}

	log.Printf("[AddOneImage] Uploaded File: %+v userId: %+v isPrivate: %+v\n", header.Filename, userId, isPrivateB)

	tempFile, err := ioutil.TempFile(REPO_PATH, "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	n, err := tempFile.Write(fileBytes)
	if err != nil {
		fmt.Println(err)
	}

	image := Image{
		Name:      imageName,
		OwnerId:   uint(userIdInt),
		Path:      REPO_PATH + tempFile.Name(),
		IsPrivate: isPrivateB,
		CreatedAt: time.Now(),
	}
	result := s.db.Create(&image)
	if result.Error != nil {
		log.Println(result.Error)
		fmt.Fprintf(w, "500 Internal Server Error")
		w.WriteHeader(500)
		return
	}

	log.Printf("[AddOneImage] File stored: %+v size %+v\n", tempFile.Name(), n)

	w.WriteHeader(200)
	fmt.Fprintf(w, "Successfully Uploaded File\n")
	return
}

func AddImages(w http.ResponseWriter, r *http.Request) {

}
