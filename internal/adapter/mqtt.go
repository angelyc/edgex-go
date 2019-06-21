package adapter

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"strconv"
	"strings"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type mqttClient struct {
	client mqtt.Client
	topic  string
}

func newMqttClient(addr address, cert string, key string) client {
	protocol := strings.ToLower(addr.Protocol)
	opts := MQTT.NewClientOptions()
	broker := protocol + "://" + addr.Address + ":" + strconv.Itoa(addr.Port) + addr.Path
	opts.AddBroker(broker)
	opts.SetClientID(addr.Publisher)
	opts.SetUsername(addr.User)
	opts.SetPassword(addr.Password)
	opts.SetAutoReconnect(false)

	if protocol == "tcps" || protocol == "ssl" || protocol == "tsl" {
		cert, err := tls.LoadX509KeyPair(cert, key)
		if err != nil {
			fmt.Println("ERROR: failed loading x509 data")
			return nil
		}

		tlsconfig := &tls.Config{
			ClientCAs:          nil,
			InsecureSkipVerify: true,
			Certificates:       []tls.Certificate{cert},
		}
		opts.SetTLSConfig(tlsconfig)
	}

	sender := &mqttClient{
		client: MQTT.NewClient(opts),
		topic:  addr.Topic,
	}
	return sender
}

func (client *mqttClient) Sender(data []byte, ctx context.Context) bool {
	if !client.client.IsConnected() {
		fmt.Println("INFO: Connecting to mqtt server")
		if token := client.client.Connect(); token.Wait() && token.Error() != nil {
			fmt.Println(fmt.Sprintf("ERROR: Could not connect to mqtt server, drop event: %s", token.Error().Error()))
			return false
		}
	}
	token := client.client.Publish(client.topic, 0, false, data)
	token.Wait()
	if token.Error() != nil {
		fmt.Println(token.Error().Error())
		return false
	} else {
		fmt.Println(fmt.Sprintf("Sent data: %X", data))
		return true
	}
}
