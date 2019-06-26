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
