package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
)

func main() {
	http.HandleFunc("/create", createScript)
	http.ListenAndServe(":8080", nil)
}

type Script struct {
	UUID     uuid.UUID
	Server   string
	Queries  []string
	Headers  map[string]string
	CalcFunc string
}

func createScript(w http.ResponseWriter, r *http.Request) {
	// Get the server and headers from the request headers
	server := r.Header.Get("X-Server")
	if server == "" {
		http.Error(w, "Missing 'X-Server' in request headers", http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get("X-Header-Content-Type")
	if contentType == "" {
		http.Error(w, "Missing 'X-Header-Content-Type' in request headers", http.StatusBadRequest)
		return
	}

	headers := map[string]string{
		"Content-Type": contentType,
	}

	// Parse the request body to get the queries and calcFunc
	var requestBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	calcFunc, ok := requestBody["calcFunc"].(string)
	if !ok {
		http.Error(w, "Missing 'calcFunc' in request body", http.StatusBadRequest)
		return
	}

	queriesInterface, ok := requestBody["queries"].([]interface{})
	if !ok {
		http.Error(w, "Missing 'queries' in request body", http.StatusBadRequest)
		return
	}

	// Convert the queries to a slice of strings
	queries := make([]string, len(queriesInterface))
	for i, queryInterface := range queriesInterface {
		query, ok := queryInterface.(string)
		if !ok {
			http.Error(w, "Invalid 'queries' in request body", http.StatusBadRequest)
			return
		}
		queries[i] = query
	}

	id := uuid.New()
	script := Script{
		UUID:     id,
		Server:   server,
		Queries:  queries,
		Headers:  headers,
		CalcFunc: calcFunc,
	}

	// Convert the script to JSON
	jsonData, err := json.Marshal(script)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the JSON data to a file
	err = ioutil.WriteFile(fmt.Sprintf("%s.json", id), jsonData, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Script created with UUID: %s\n", id)
}
