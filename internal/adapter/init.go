package adapter

import (
	"fmt"
	"github.com/edgexfoundry/edgex-go/internal/adapter/interfaces"
	"github.com/edgexfoundry/edgex-go/internal/pkg/config"
	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
	"github.com/edgexfoundry/edgex-go/internal/pkg/db/mongo"
	"github.com/edgexfoundry/edgex-go/internal/pkg/db/redis"
	"github.com/edgexfoundry/edgex-go/internal/pkg/telemetry"
	mdclients "github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"sync"
	"time"
)

var Configuration *ConfigurationStruct
//var registryClient registry.Client
var dbClient interfaces.DBClient
var LoggingClient logger.LoggingClient
//var nc notifications.NotificationsClient
var publisher client
var GatewayId string
func Retry(useRegistry bool, useProfile string, timeout int, wait *sync.WaitGroup, ch chan error) {
	until := time.Now().Add(time.Millisecond * time.Duration(timeout))
	for time.Now().Before(until) {
		var err error
		//When looping, only handle configuration if it hasn't already been set.
		if Configuration == nil {
			Configuration, err = initializeConfiguration(useProfile)
			if err != nil {
				ch <- err
			} else {
				// Initialize notificationsClient based on configuration
				initializeClients()
				// Setup Logging
				logTarget := setLoggingTarget()
				LoggingClient = logger.NewClient(mdclients.CoreMetaDataServiceKey,
					Configuration.Logging.EnableRemote, logTarget, Configuration.Writable.LogLevel)
			}
		}

		//Only attempt to connect to database if configuration has been populated
		if Configuration != nil {
			err := connectToDatabase()
			if err != nil {
				ch <- err
			} else {
				break
			}
		}
		time.Sleep(time.Second * time.Duration(1))
	}
	close(ch)
	wait.Done()

	return
}

func Init() bool {
	if Configuration == nil {
		return false
	}

	go telemetry.StartCpuUsageAverage()

	return true
}

func connectToDatabase() error {
	var err error

	dbClient, err = newDBClient()
	if err != nil {
		dbClient = nil
		return fmt.Errorf("couldn't create database client: %v", err.Error())
	}

	return nil
}

// Return the dbClient interface
func newDBClient() (interfaces.DBClient, error) {
	switch Configuration.Databases["Primary"].Type {
	case db.MongoDB:
		dbConfig := db.Configuration{
			Host:         Configuration.Databases["Primary"].Host,
			Port:         Configuration.Databases["Primary"].Port,
			Timeout:      Configuration.Databases["Primary"].Timeout,
			DatabaseName: Configuration.Databases["Primary"].Name,
			Username:     Configuration.Databases["Primary"].Username,
			Password:     Configuration.Databases["Primary"].Password,
		}
		return mongo.NewClient(dbConfig)
	case db.RedisDB:
		dbConfig := db.Configuration{
			Host: Configuration.Databases["Primary"].Host,
			Port: Configuration.Databases["Primary"].Port,
		}
		return redis.NewClient(dbConfig) //TODO: Verify this also connects to Redis
	default:
		return nil, db.ErrUnsupportedDatabase
	}
}
func initializeConfiguration(useProfile string) (*ConfigurationStruct, error) {
	//We currently have to load configuration from filesystem first in order to obtain RegistryHost/Port
	configuration := &ConfigurationStruct{}
	err := config.LoadFromFile(useProfile, configuration)
	if err != nil {
		return nil, err
	}

	return configuration, nil
}

func Destruct() {
	if dbClient != nil {
		dbClient.CloseSession()
		dbClient = nil
	}
}
func initializeClients() {
	// Create notification client
	cert := ""
	key := ""
	InitGatewayID()
	SubInit()
	mqttClient := newMqttClient(Configuration.MQTT, cert, key)
	publisher = &MqttClient{
		client: mqttClient,
	}

	Register()
}
func setLoggingTarget() string {
	if Configuration.Logging.EnableRemote {
		return Configuration.Clients["Logging"].Url() + mdclients.ApiLoggingRoute
	}
	return Configuration.Logging.File
}
