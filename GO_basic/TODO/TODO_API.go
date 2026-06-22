package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Todo struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
	Done bool   `json:"done"`
}

var Todos []Todo
var nextID int = 1

func addTodo(w http.ResponseWriter, r *http.Request) {
	var newTodo Todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newTodo.ID = nextID
	nextID++
	Todos = append(Todos, newTodo)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTodo)
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Todos)
}

func main() {
	http.HandleFunc("/addTodo", addTodo)
	http.HandleFunc("/getTodos", getTodos)

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
