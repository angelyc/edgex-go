package adapter

import (
	"encoding/json"
	"github.com/robfig/cron"
	"strconv"
)

var Register2CloudTime string
var heartBeatChannel chan bool
var interval int
func Register() {
	interval = 30
	heartCron(interval)
}

func heardBeat() {
	gi := HeardBeatInfo{Id: GatewayId, Timeout: interval * 4}
	msg, _ := json.Marshal(gi)
	publisher.Sender(HeardBeatTopic, msg)
}

func heartCron(interval int) {
	cron := cron.New()
	cron.Start()
	cron.AddFunc("@every " +strconv.Itoa(interval) + "s", heardBeat)
	heardBeat()
}


