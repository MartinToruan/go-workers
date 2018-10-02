package usecase

import (
	"math/rand"
	"workers/examples/workers_tdk/pkg/domain"
	"workers/examples/workers_tdk/pkg/lib/config"
	w "workers/pkg/worker"
)

type Order struct {
	UserID    int `json:"user_id"`
	Quantity  int `json:"quantity"`
	ProductID int `json:"product_id"`
}

type OrderUsecase struct {
	cfg    config.Config
	worker w.Worker
	od     domain.OrderDomain
}

func InitOrderUsecase(cfg config.Config, work w.Worker, order domain.OrderDomain) *OrderUsecase {
	return &OrderUsecase{
		cfg:    cfg,
		worker: work,
		od:     order,
	}
}

func (o *OrderUsecase) PutNewOrder(order Order) (string, error) {

	// create new Order entity of order domain
	newOrder := domain.Order{
		OrderID:   rand.Intn(100),
		Quantity:  order.Quantity,
		ProductID: order.ProductID,
	}

	job := func() error {
		if err := o.od.CreateOrder(&newOrder); err != nil {
			return err
		}
		return nil
	}

	go o.worker.PushJob(uint(newOrder.OrderID), uint8(1), job)

	return newOrder.Invoice, nil
}
