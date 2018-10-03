package domain

import (
	"errors"
	"fmt"
)

type Order struct {
	OrderID   int
	ProductID int
	Quantity  int
	Invoice   string
}

type OrderDomain struct {
	resource OrderResourceItf
}

func InitOrderDomain(rsc OrderResourceItf) OrderDomain {
	return OrderDomain{
		resource: rsc,
	}
}

func (d OrderDomain) CreateOrder(order *Order) error {
	// lets generate invoice before we apply it to DB
	idinvoice := order.OrderID
	order.Invoice = fmt.Sprintf("INV/%d", idinvoice)

	// apply it to database
	if err := d.resource.InsertOrder(order); err != nil {
		return err
	}

	return nil
}

type OrderResourceItf interface {
	InsertOrder(*Order) error
}

type OrderResource struct {
}

func (rsc OrderResource) InsertOrder(order *Order) error {

	fmt.Printf(`try to save order to DB :
				orderID %d
				productID %d
				Quantity %d
				Invoice %s
		`,
		order.OrderID,
		order.ProductID,
		order.Quantity,
		order.Invoice)

	if order.OrderID > 75 {
		return errors.New("error : simulate error when orderID > 75 the job should be retried")
	}
	fmt.Println("success : order has been stored to DB!")
	return nil
}
