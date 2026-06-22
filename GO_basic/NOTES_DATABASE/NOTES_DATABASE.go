package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type note struct {
	Name      string    `json:"name"`
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

var NOTES []note
var nextID = 1

func createNote(w http.ResponseWriter, r *http.Request) {
	var newnote note
	err := json.NewDecoder(r.Body).Decode(&newnote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newnote.ID = nextID
	nextID++
	newnote.CreatedAt = time.Now()
	NOTES = append(NOTES, newnote)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newnote)
}
func deleteNote(w http.ResponseWriter, r *http.Request) {
	var tempid int
	err := json.NewDecoder(r.Body).Decode(&tempid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for i := range NOTES {
		if NOTES[i].ID == tempid {
			NOTES = append(NOTES[:i], NOTES[i+1:]...)
			break
		}
	}
	w.Write([]byte("Note deleted successfully"))
}
func editNote(w http.ResponseWriter, r *http.Request) {
	type editRequest struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	}
	var request editRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for i := range NOTES {
		if NOTES[i].Name == request.Name {
			NOTES[i].Content = request.Content
		}
	}
	w.Write([]byte("Note Sucessfully Edited"))
	json.NewEncoder(w).Encode(request)
}
func showallNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(NOTES)
}
func main() {
	http.HandleFunc("/createNote", createNote)
	http.HandleFunc("/deleteNote", deleteNote)
	http.HandleFunc("/editNote", editNote)
	http.HandleFunc("/getNotes", showallNote)

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
