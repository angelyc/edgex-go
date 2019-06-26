package adapter

import (
	"github.com/eclipse/paho.mqtt.golang"
)

type client interface {
	Sender(topic string, data interface{}) bool
}

type Addressable struct {
	Protocol  string `json:"protocol"`    // Protocol for the address (HTTP/TCP)
	Address   string `json:"address"`     // Address of the addressable
	Port      int    `json:"port,Number"` // Port for the address
	Publisher string `json:"publisher"`   // For message bus protocols
	User      string `json:"user"`        // User id for authentication
	Password  string `json:"password"`    // Password of the user for authentication for the addressable
	Topic     string `json:"topic"`       // Topic for message bus addressables
}

type NetInfo struct {
	Name        string //网卡名称
	DisplayName string //网卡的显示名称
	MacAddr     string //Mac 地址
	IPv4Addr    string //ipv4 地址
	IPv6Address string //ipv6 地址
	IsUp        bool   //网卡状态是否UP
}

type MqttClient struct {
	client mqtt.Client
}

type GatewayInfo struct {
	Id                    string  `json:"id"`                    //网关ID
	Name                  string  `json:"name"`                  //网关名称
	Manufacturer          string  `json:"manufacturer"`          //网关的制造商
	BaseboardSerialNumber string  `json:"baseboardSerialNumber"` //主板的序列号
	ProcessorID           string  `json:"processorID"`           //CPU ID
	RegisterTime          string  `json:"registerTime"`          //最近一次网关向云平台的注册时间
	LocationName          string  `json:"locationName"`          //地点名称
	LocationCode          string  `json:"locationCode"`          //地点编码
	NetIfs                []NetIf `json:"netIfs"`                //网卡信息
}
type NetIf struct {
	Name        string `json:"name"`        //网卡名称
	DisplayName string `json:"displayName"` //网卡的显示名称
	MacAddr     string `json:"macAddr"`     //Mac 地址
	Ipv4Addr    string `json:"ipv4Addr"`    //ipv4 地址
	Ipv6Addr    string `json:"Ipv6Addr"`    //ipv6 地址
	IsUp        bool   `json:"isUp"`        //网卡状态是否UP
}
type HeardBeatInfo struct {
	Id      string `json:"id"` // gateway id
	Timeout int    `json:"timeout"`
}

type HttpRequest struct {
	Method string `json:"method"`
	Url    string `json:"url"`
	Data   []byte `json:"data"`
}

type MqttRequest struct {
	Http      HttpRequest `json:"http"`
	AuthToken string      `json:"authToken"`
	Context   int         `json:"context"`
	Timestamp string      `json:"timestamp"`
	Timeout   int         `json:"timeout"`
}

type MqttResponse struct {
	Http    HttpResponse `json:"http"`
	Context int          `json:"context"`
}

type HttpResponse struct {
	HttpCode int    `json:"httpCode"`
	Code     int    `json:"code"`
	Result   string `json:"result"`
}
