package adapter

import (
	"crypto/tls"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"strconv"
	"strings"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func newMqttClient(addr address, cert string, key string) mqtt.Client {
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
	opts.OnConnect = Subscribe
	return MQTT.NewClient(opts)
}

func (client *MqttClient) Sender(topic string, data interface{}) bool {
	if !client.client.IsConnected() {
		fmt.Println("INFO: Connecting to mqtt server")
		if token := client.client.Connect(); token.Wait() && token.Error() != nil {
			fmt.Println(fmt.Sprintf("ERROR: Could not connect to mqtt server, drop event: %s", token.Error().Error()))
			return false
		}
	}
	token := client.client.Publish(topic, 0, false, data)
	token.Wait()
	if token.Error() != nil {
		fmt.Println(token.Error().Error())
		return false
	} else {
		fmt.Println(fmt.Sprintf("Sent data: %s", data))
		return true
	}
}
