package main

import (
	"app/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", HelloServer)
	http.HandleFunc("/bot", runBot)
	http.ListenAndServe(":8080", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	fmt.Printf("Hello, %s!", r.URL.Path[1:])
}

func runBot(w http.ResponseWriter, r *http.Request) {
	var request models.Request

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Received message: %+v", request)
}
