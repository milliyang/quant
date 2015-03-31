package base

import (
	"quant/base/order"
	"quant/base/xbase"
)

var (
	orderManager = map[int]*order.Order{}

	//
	dicingNewGameOrder = []*order.Order{}
)

func Send(order_ *order.Order) {
	orderManager[order_.Id] = order_

	if xbase.CasinoDicingGame {
		dicingNewGameOrder = append(dicingNewGameOrder, order_)
	}
}

func Get(id int) *order.Order {
	order_, ok := orderManager[id]
	if ok {
		return order_
	}
	return nil
}
