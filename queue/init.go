package queue

import (
	"hublabs/coupon-api/config"
	"os"
	"sync"

	"github.com/RichardKnop/machinery/v1"
	mcnf "github.com/RichardKnop/machinery/v1/config"
	"github.com/go-xorm/xorm"
	"github.com/labstack/gommon/log"
)

var (
	db     *xorm.Engine
	once   sync.Once
	Server *machinery.Server
)

const (
	serviceName = "coupon-api"
)

func Init(e *xorm.Engine) {
	var err error

	once.Do(func() {
		db = e
	})

	Server, err = startServer()
	if err != nil {
		panic(err)
	}

	// The second argument is a consumer tag
	// Ideally, each worker should have a unique tag (worker1, worker2 etc)
	worker := Server.NewWorker(serviceName, 0)

	err = worker.Launch()
	log.Info(err)
	os.Exit(0)
}

func startServer() (*machinery.Server, error) {
	c := mcnf.Config{
		Broker:        config.Config().Queue.Broker,
		ResultBackend: config.Config().Queue.Broker,
		DefaultQueue:  config.Config().Queue.DefaultQueue,
	}

	// Create server instance
	server, err := machinery.NewServer(&c)
	if err != nil {
		return nil, err
	}

	// Register tasks
	tasks := map[string]interface{}{
		"send_coupon": SendCoupon,
	}

	return server, server.RegisterTasks(tasks)
}
