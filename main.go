package main

import (
	"github.com/gitslagga/gitbitex-spot/conf"
	"github.com/gitslagga/gitbitex-spot/matching"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/mylog"
	"github.com/gitslagga/gitbitex-spot/pushing"
	"github.com/gitslagga/gitbitex-spot/rest"
	"github.com/gitslagga/gitbitex-spot/service"
	"github.com/gitslagga/gitbitex-spot/tasks"
	"github.com/gitslagga/gitbitex-spot/worker"
	"github.com/prometheus/common/log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	gbeConfig := conf.GetConfig()
	mylog.ConfigLoggers()

	go func() {
		log.Info(http.ListenAndServe("localhost:6060", nil))
	}()

	go models.NewBinLogStream().Start()

	matching.StartEngine()

	pushing.StartServer()

	worker.NewFillExecutor().Start()
	worker.NewBillExecutor().Start()
	products, err := service.GetProducts()
	if err != nil {
		panic(err)
	}
	for _, product := range products {
		worker.NewTickMaker(product.Id, matching.NewKafkaLogReader("tickMaker", product.Id, gbeConfig.Kafka.Brokers)).Start()
		worker.NewFillMaker(matching.NewKafkaLogReader("fillMaker", product.Id, gbeConfig.Kafka.Brokers)).Start()
		worker.NewTradeMaker(matching.NewKafkaLogReader("tradeMaker", product.Id, gbeConfig.Kafka.Brokers)).Start()
	}

	rest.StartServer()

	go tasks.StartMachineRelease()
	go tasks.StartTokenInfo()
	go tasks.StartTransactionScan()
	go tasks.StartSendToMainTask()

	select {}
}
