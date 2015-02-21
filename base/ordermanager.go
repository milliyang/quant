package base

import (
	"quant/base/order"
)

var (
	orderManager = map[int]*order.Order{}
)

func Send(order_ *order.Order) {
	orderManager[order_.Id] = order_
}

func Get(id int) *order.Order {
	order_, ok := orderManager[id]
	if ok {
		return order_
	}
	return nil
}
