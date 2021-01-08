package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func saveImage() {

}

func downloadImageByUrl(url string) {

}

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

	isPrivate := r.FormValue("isPrivate")
	isPrivateB, err := strconv.ParseBool(isPrivate)
	if err != nil {
		fmt.Fprintf(w, "401 bad request: isPrivate needs to be a bool value")
		w.WriteHeader(401)
		return
	}

	log.Printf("[AddOneImage] Uploaded File: %+v isPrivate: %+v\n", header.Filename, isPrivateB)

	tempFile, err := ioutil.TempFile("./images", "upload-*.png")
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
	log.Printf("[AddOneImage] File stored: %+v size %+v\n", tempFile.Name(), n)

	w.WriteHeader(200)
	fmt.Fprintf(w, "Successfully Uploaded File\n")
	return
}

func AddImages(w http.ResponseWriter, r *http.Request) {

}
