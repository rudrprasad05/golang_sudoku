package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Message struct {
	Data string `json:"data"`
}

func (routes *Routes) Handle404(w http.ResponseWriter, r *http.Request){
	routes.LOG.Error(fmt.Sprintf("404 Not Found: %s", r.URL.Path))
	data := Message{Data: "404 not found"}
	sendJSONResponse(w, http.StatusNotFound, data)
}


func (routes *Routes) GetHome(w http.ResponseWriter, r *http.Request){
	data := Message{Data: "hello"}
	sendJSONResponse(w, http.StatusOK, data)
}


func sendJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Encode data to JSON and send response
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}