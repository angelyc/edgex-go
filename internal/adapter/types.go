package adapter

import "context"

type client interface {
	Sender(data []byte, ctx context.Context) bool
}

type Addressable struct {
	Protocol    string `json:"protocol"`    // Protocol for the address (HTTP/TCP)
	Address     string `json:"address"`     // Address of the addressable
	Port        int    `json:"port,Number"` // Port for the address
	Publisher   string `json:"publisher"`   // For message bus protocols
	User        string `json:"user"`        // User id for authentication
	Password    string `json:"password"`    // Password of the user for authentication for the addressable
	Topic       string `json:"topic"`       // Topic for message bus addressables
}