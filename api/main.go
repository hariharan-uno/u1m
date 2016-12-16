package main

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
)

// this service's configuration
type specification struct {
	Bind string `envconfig:"bind" default:":8080"`
	DB   string `envconfig:"db"`
}

func main() {
	var err error

	// Set up our logging options
	logger := logrus.New()
	logger.Formatter = new(logrus.TextFormatter)
	logger.Level = logrus.DebugLevel

	var spec specification
	err = envconfig.Process("APP", &spec)
	if err != nil {
		logger.Fatalln(err)
	}
	logger.Info(spec)

	s, err := NewServer()
	if err != nil {
		logger.Fatalln(err)
	}
	s.log = logger
	for {
		if err = s.ConnectDB(spec.DB); err != nil {
			logger.Warnln("Problem connecting to DB", err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	defer s.Close()

	logger.Info("Starting API server")
	err = http.ListenAndServe(spec.Bind, s.Router())
	if err != nil {
		logger.Errorln(err)
	}
}
