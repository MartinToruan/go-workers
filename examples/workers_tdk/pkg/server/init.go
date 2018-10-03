package server

import (
	"fmt"
	"workers/examples/workers_tdk/pkg/domain"
	"workers/examples/workers_tdk/pkg/lib/config"
	"workers/examples/workers_tdk/pkg/usecase"
	w "workers/pkg/worker"

	"github.com/tokopedia/tdk/go/app"
)

var workers w.Worker
var cfg config.Config

var orderDomain domain.OrderDomain
var orderUsecase *usecase.OrderUsecase

// We do all the wire up in this Init() function
// please return any error if you fail to initialize something
func Init(app *app.App) error {
	cfg = config.GetConfig()

	//set Workers
	workers = w.NewWorkers(100, 5052)
	workers.Run()

	//set Worker Listener
	go func() {
		for err := range workers.PollJob() {
			fmt.Println(err)
		}
	}()

	orderDomain = domain.InitOrderDomain(
		domain.OrderResource{},
	)

	orderUsecase = usecase.InitOrderUsecase(
		cfg,
		workers,
		orderDomain,
	)
	return nil
}
