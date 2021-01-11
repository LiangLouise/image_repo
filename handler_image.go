package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

// POST FormData
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
		http.Error(w, "400 bad request: not a valid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imageName := r.FormValue("imageName")
	if imageName == "" {
		http.Error(w, "400 bad request: require imageName", http.StatusBadRequest)
		return
	}

	isPrivate := r.FormValue("isPrivate")
	isPrivateB, err := strconv.ParseBool(isPrivate)
	if err != nil {
		http.Error(w, "400 bad request: isPrivate needs to be a bool value", http.StatusBadRequest)
		return
	}

	var user User
	userId := r.FormValue("userId")
	userIdInt, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		http.Error(w, "400 bad request: require the userId", http.StatusBadRequest)
		return
	}
	user.ID = uint(userIdInt)
	result := s.db.First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("[AddOneImage] user id: %v not found\n", user.ID)
			http.Error(w, "404: input userId not found", http.StatusNotFound)
			return
		} else {
			log.Println(result.Error)
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}
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
		OwnerId:   user.ID,
		Path:      "./" + tempFile.Name(),
		IsPrivate: isPrivateB,
		CreatedAt: time.Now(),
	}
	result = s.db.Create(&image)
	if result.Error != nil {
		log.Println(result.Error)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("[AddOneImage] File stored: %+v size %+v\n", tempFile.Name(), n)

	w.WriteHeader(200)
	fmt.Fprintf(w, "Successfully Uploaded File\n")
	return
}

type DeleteImageArgs struct {
	UserId  uint `json:"userid"`
	ImageId uint `json:"imageid"`
}

// DELETE
// Request Body
// 	- UserId
// 	- ImageId
func (s *Server) DeleteImage(w http.ResponseWriter, r *http.Request) {
	var del DeleteImageArgs

	reqBody, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(reqBody, &del)
	log.Printf("[SignUp] request: %v\n", string(reqBody))
	if err != nil {
		http.Error(w, "400 bad request: missing required field or invalid json", http.StatusBadRequest)
		return
	}

	image := Image{
		ID: del.ImageId,
	}

	result := s.db.First(&image)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("[DeleteImage]image: %v not found\n", del.ImageId)
			http.Error(w, "404: target image not found", http.StatusNotFound)
			return
		} else {
			log.Println(result.Error)
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else if image.OwnerId != del.UserId {
		log.Printf("[DeleteImage]image %v: wrong owner %v \n", del.ImageId, del.UserId)
		http.Error(w, "403: you are not owner", http.StatusForbidden)
		return
	} else {
		go func() {
			err := os.Remove(image.Path)
			if err != nil {
				log.Printf("[DeleteImage]Unable to delete image %v: wrong path %v \n", del.ImageId, err)
			}
		}()
		s.db.Delete(&image)
		w.WriteHeader(200)
		return
	}
}

func (s *Server) GetOneImage(w http.ResponseWriter, r *http.Request) {
	imageId := mux.Vars(r)["id"]
	userId := r.FormValue("userid")

	imageIdInt, err := strconv.ParseUint(imageId, 10, 32)
	if err != nil {
		http.Error(w, "400 bad request: require the userId", http.StatusBadRequest)
		return
	}

	userIdInt, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		http.Error(w, "400 bad request: require the userId", http.StatusBadRequest)
		return
	}
	image := Image{
		ID: uint(imageIdInt),
	}

	result := s.db.First(&image)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("[DeleteImage]image: %v not found\n", image.ID)
			http.Error(w, "404: target image not found", http.StatusNotFound)
			return
		} else {
			log.Println(result.Error)
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else if image.OwnerId != uint(userIdInt) && image.IsPrivate {
		log.Printf("[DeleteImage]image %v: wrong user %v \n", image.ID, userIdInt)
		http.Error(w, "403: The image is private", http.StatusForbidden)
		return
	} else {
		w.Header().Set("Content-Type", "image/png")
		file, err := os.Open(image.Path)

		if err != nil {
			log.Println(result.Error)
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		FileStat, _ := file.Stat()
		//Get file size as a string
		FileSize := strconv.FormatInt(FileStat.Size(), 10)
		w.Header().Set("Content-Length", FileSize)
		io.Copy(w, file)
		return
	}
}
