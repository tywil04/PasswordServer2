package routes

import (
	"encoding/json"
	"net/http"
)

func GetTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"Hello": "World!"})
}
