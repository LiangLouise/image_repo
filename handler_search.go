package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type SearchReplyImage struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type SearchReply struct {
	Images  []SearchReplyImage `json:"images"`
	HasNext bool               `json:"hasNext"`
}

func (s *Server) SearchImages(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("userid")
	text := r.FormValue("text")
	page := r.FormValue("page")

	userIdInt, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		http.Error(w, "400 bad request: require the userId", http.StatusBadRequest)
		return
	}

	pageInt, err := strconv.ParseInt(page, 10, 32)
	if err != nil {
		http.Error(w, "400 bad request: require the userId", http.StatusBadRequest)
		return
	}

	var images []Image
	// Select image from images where title like %text% and (Isprivate = false or OwnerId = userId)
	s.db.Where("Title LIKE ? AND (IsPrivate = ? OR OwnerId = ?)", "%"+text+"%", false, userIdInt).
		Offset(int(pageInt-1) * 10).
		Limit(int(pageInt) * 10).
		Find(&images)

	rep := SearchReply{
		Images:  make([]SearchReplyImage, len(images)),
		HasNext: len(images) == 10,
	}

	for i := 0; i < len(images); i++ {
		rep.Images[i].Name = images[i].Name
		rep.Images[i].Path = fmt.Sprintf("/image/%v?userid=%v", images[i].ID, userId)
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rep)
	return
}
