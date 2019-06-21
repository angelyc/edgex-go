package main

import (
	"flag"
	"fmt"
	"github.com/edgexfoundry/edgex-go"
	"github.com/edgexfoundry/edgex-go/internal"
	"github.com/edgexfoundry/edgex-go/internal/adapter"
	"github.com/edgexfoundry/edgex-go/internal/pkg/correlation"
	"github.com/edgexfoundry/edgex-go/internal/pkg/startup"
	"github.com/edgexfoundry/edgex-go/internal/pkg/usage"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/gorilla/context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)


func main() {
	start := time.Now()
	var useRegistry bool
	var useProfile string

	flag.BoolVar(&useRegistry, "registry", false, "Indicates the service should use registry service.")
	flag.BoolVar(&useRegistry, "r", false, "Indicates the service should use registry service.")
	flag.StringVar(&useProfile, "profile", "", "Specify a profile other than default.")
	flag.StringVar(&useProfile, "p", "", "Specify a profile other than default.")
	flag.Usage = usage.HelpCallback
	flag.Parse()

	params := startup.BootParams{UseRegistry: useRegistry, UseProfile: useProfile, BootTimeout: internal.BootTimeoutDefault}
	startup.Bootstrap(params, adapter.Retry, logBeforeInit)

	ok := adapter.Init()
	if !ok {
		logBeforeInit(fmt.Errorf("%s: Service bootstrap failed!", adapter.GatewayAdapterKey))
		os.Exit(1)
	}

	adapter.LoggingClient.Info("Service dependencies resolved...")
	adapter.LoggingClient.Info(fmt.Sprintf("Starting %s %s ", adapter.GatewayAdapterKey, edgex.Version))

	http.TimeoutHandler(nil, time.Millisecond*time.Duration(adapter.Configuration.Service.Timeout),
		               "Gateway adapter request timed out")
	adapter.LoggingClient.Info(adapter.Configuration.Service.StartupMsg)

	errs := make(chan error, 2)
	listenForInterrupt(errs)
	startHttpServer(errs, adapter.Configuration.Service.Port)

	// Time it took to start service
	adapter.LoggingClient.Info("Service started in: " + time.Since(start).String())
	adapter.LoggingClient.Info("Listening on port: " + strconv.Itoa(adapter.Configuration.Service.Port))
	c := <-errs
	adapter.Destruct()
	adapter.LoggingClient.Warn(fmt.Sprintf("terminating: %v", c))

	os.Exit(0)
}

func logBeforeInit(err error) {
	l := logger.NewClient(adapter.GatewayAdapterKey, false, "", models.InfoLog)
	l.Error(err.Error())
}

func listenForInterrupt(errChan chan error) {
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt)
		errChan <- fmt.Errorf("%s", <-c)
	}()
}

func startHttpServer(errChan chan error, port int) {
	go func() {
		correlation.LoggingClient = adapter.LoggingClient //Not thrilled about this, can't think of anything better ATM
		r := adapter.LoadRestRoutes()
		errChan <- http.ListenAndServe(":"+strconv.Itoa(port), context.ClearHandler(r))
	}()
}
