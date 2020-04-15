package config

import (
	"os"

	configutil "github.com/pangpanglabs/goutils/config"
	"github.com/pangpanglabs/goutils/echomiddleware"
	"github.com/pangpanglabs/goutils/jwtutil"
	"github.com/sirupsen/logrus"
)

var config C

func Init(appEnv string, options ...func(*C)) C {
	config.AppEnv = appEnv
	if err := configutil.Read(appEnv, &config); err != nil {
		logrus.WithError(err).Warn("load config file error")
	}

	if s := os.Getenv("JWT_SECRET"); s != "" {
		config.JwtSecret = s
		jwtutil.SetJwtSecret(s)
	}

	for _, option := range options {
		option(&config)
	}
	return config
}

func Config() C {
	return config
}

type C struct {
	Database struct {
		Driver     string
		Connection string
		Logger     struct {
			Kafka echomiddleware.KafkaConfig
		}
	}
	Queue struct {
		Broker       string
		DefaultQueue string
	}
	BehaviorLog struct {
		Kafka echomiddleware.KafkaConfig
	}
	AppEnv, JwtSecret, ServiceName string
	Debug                          bool
}
