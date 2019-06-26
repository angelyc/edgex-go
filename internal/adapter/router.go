package adapter

import (
	"encoding/json"
	"github.com/edgexfoundry/edgex-go/internal/pkg/telemetry"
	mdclient "github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/gorilla/mux"
	"net/http"
)

func LoadRestRoutes() *mux.Router {
	r := mux.NewRouter()

	// Ping Resource
	r.HandleFunc(mdclient.ApiPingRoute, pingHandler).Methods(http.MethodGet)

	// Configuration
	r.HandleFunc(mdclient.ApiConfigRoute, configHandler).Methods(http.MethodGet)

	// Metrics
	r.HandleFunc(mdclient.ApiMetricsRoute, metricsHandler).Methods(http.MethodGet)

	b := r.PathPrefix(mdclient.ApiBase).Subrouter()

	b.HandleFunc("/"+Gateway_Info, restGetGatewayInfo).Methods(http.MethodGet)
	loadStaticsRoutes(b)
	return r
}

func loadStaticsRoutes(b *mux.Router) {
	// /api/v1/" + STATICS
	b.HandleFunc("/"+STATICS, restClearStatics).Methods(http.MethodPut)
	b.HandleFunc("/"+STATICS, restGetStatics).Methods(http.MethodGet)
}

func pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	msg := "I am GOGOGO!!!"
	data2 := []byte(msg)
	publisher.Sender("gateway-adapter", data2)
	w.Write([]byte("pong"))
}

func configHandler(w http.ResponseWriter, _ *http.Request) {
	encode(Configuration, w)
}

func metricsHandler(w http.ResponseWriter, _ *http.Request) {
	s := telemetry.NewSystemUsage()

	encode(s, w)

	return
}

// Helper function for encoding things for returning from REST calls
func encode(i interface{}, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	err := enc.Encode(i)
	// Problems encoding
	if err != nil {
		LoggingClient.Error("Error encoding the data: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
