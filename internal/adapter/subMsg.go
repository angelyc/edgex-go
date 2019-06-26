package adapter

import (
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"net/http"
	"strings"
)

var subTopic string
var pubTopic string
var qos byte

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	//fmt.Printf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())
	// {"http":{"method":"GET","url":"/gateway-adapter/api/v1/gateway-info"},"context":1561000330711,"timestamp":"2019-06-25 17:24:10","timeout":120}
	var mreq MqttRequest
	var mres MqttResponse
	json.Unmarshal(message.Payload(), &mreq)
	mres.Http.Code = 1001
	var resp *http.Response
	var err error
	/*until := time.Now().Truncate(time.Millisecond * time.Duration(mr.Timeout)).String()
	a := strings.Compare(mr.Timestamp, until)
	fmt.Printf("Received message on topic: %d\n", a)*/
	//apiIndex := strings.Index(mr.Http.Url, mdclient.ApiBase)
	a := strings.SplitN(mreq.Http.Url, "/", 3)
	if len(a) == 3 {
		resp, err = doGet("http://localhost:48011/" + a[2])
		if err == nil {
			mres.Http.HttpCode = resp.StatusCode
			mres.Http.Code = 0
			var result interface{}
			json.NewDecoder(resp.Body).Decode(&result)
			b, _ :=json.Marshal(result)
			mres.Http.Result = string(b)
			/*b, _:= ioutil.ReadAll(resp.Body)
			json.Unmarshal(b, &)
			*/
			mres.Context = mreq.Context
			fmt.Printf("http response: %s\n", mres.Http.Result)
			/*mres = MqttResponse{Http: HttpResponse{HttpCode: resp.StatusCode,
				Code:        0,
				Description: "",
				Result:      data,},
				Context: mreq.Context,
			}*/

		}
	}
	msg, _ := json.Marshal(mres)
	publisher.Sender(pubTopic, msg)
}

func Subscribe(c MQTT.Client) {
	if token := c.Subscribe(subTopic, qos, onMessageReceived); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func SubInit() {
	subTopic = fmt.Sprintf("/gateway/%s/send", strings.ToUpper(GatewayId))
	pubTopic = fmt.Sprintf("/gateway/%s/recv", strings.ToUpper(GatewayId))
	qos = 1
}
