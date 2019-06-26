package adapter

import (
	"encoding/json"
	"net/http"
)

func restGetGatewayInfo(w http.ResponseWriter, _ *http.Request) {
	res := ReadFile()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&res)
}

func restSetGatewayInfoName(w http.ResponseWriter, r *http.Request) {
	var name string

	err := json.NewDecoder(r.Body).Decode(&name)
	if err != nil {
		LoggingClient.Error(err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
}

func restSetGatewayInfoLocation(w http.ResponseWriter, _ *http.Request) {
	res := ReadFile()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&res)
}
