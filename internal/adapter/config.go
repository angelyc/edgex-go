package adapter

import (
	"github.com/edgexfoundry/edgex-go/internal/pkg/config"
)

type ConfigurationStruct struct {
	Service   config.ServiceInfo
	Clients   map[string]config.ClientInfo
	Logging   config.LoggingInfo
	Writable  WritableInfo
	Databases map[string]config.DatabaseInfo
	MQTT      address
}


type WritableInfo struct {
	LogLevel string
}

type address struct {
	Publisher string
	User      string
	Password  string
	Address   string
	Port      int
	Protocol  string
	Path      string
}

/*
type address struct {
	Publisher string `json:"publisher"`
	User      string `json:"user"`
	Password  string `json:"password"`
	Address   string `json:"address"`
	Port      int    `json:"port,Number"`
	Protocol  string `json:"protocol"`
	Path      string `json:"path"`
	Topic     string `json:"topic"`
}*/
