package adapter

import (
	"encoding/json"
	"net/http"
)

func restGetStatics(w http.ResponseWriter, _ *http.Request) {
	res := "success"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&res)
}
func restClearStatics(w http.ResponseWriter, _ *http.Request) {
	res := "success"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&res)
}