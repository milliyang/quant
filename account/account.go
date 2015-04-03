package account

import (
	"quant/base/indicator"
	"quant/base/order"
)

type Account struct {
	Name     string
	Position map[string]*Position
	trans    map[int]*Transaction

	InitialWealth   float64
	availableAmount float64
	FinalWealth     float64
	PnL             float64

	PanddingOrders map[int]*order.Order // key:id

	// CasinoDicingGame
	IndicatorPNL         *indicator.DicePNL
	IndicatorPerformance *indicator.DicePerformance
}

func NewAccount(name string, amount float64) *Account {
	acc := &Account{}
	acc.Name = name
	acc.InitialWealth = amount
	acc.availableAmount = amount
	acc.IndicatorPNL = indicator.NewDicePNL(amount)
	acc.IndicatorPerformance = indicator.NewDicePerformance(amount)
	return acc
}

func (this *Account) CheckOrder(order_ *order.Order) bool {
	if order_.Status != order.StatusPendingNew {
		panic(order_)
	}

	if order_.Side == order.SideBuy {
		// BUY
		order_.Status = order.StatusNew
		order_.Price = 0 // this price is unknown until tomorrow
		return true
	} else {
		// SELL

		pos, ok := this.Position[order_.Symbol]
		if ok {
			if pos.Qty >= order_.Qty {
				order_.Status = order.StatusNew
				return true
			} else {
				order_.Status = order.StatusRejected
				return false
			}
		}
		return false
	}
}

func (this *Account) HandleOrder(order_ *order.Order) bool {
	switch order_.Status {
	case order.StatusNew:
		return this.handleNewOrder(order_)
	case order.StatusRejected:
		return this.handleRejectedOrder(order_)
	case order.StatusFilled:
		return this.handleFilledOrder(order_)
	default:
		panic(order_)
	}
}

func (this *Account) handleNewOrder(order_ *order.Order) bool {
	if order_.Side == order.SideBuy {
		this.PanddingOrders[order_.Id] = order_

		amt := float64(order_.Qty) * order_.Price

		this.availableAmount -= amt

		return true
	} else if order_.Side == order.SideSell {
		pos, ok := this.Position[order_.Symbol]
		if !ok {
			panic(order_)
		}

		pos.Qty -= order_.Qty

		// Send Order
		return true
	} else {
		panic(order_)
	}
}

func (this *Account) handleRejectedOrder(order_ *order.Order) bool {
	o, ok := this.PanddingOrders[order_.Id]
	if !ok {
		panic(order_)
	}
	o.Status = order_.Status

	if order_.Side == order.SideBuy {
		// BUY failed.
		return true
	} else if order_.Side == order.SideSell {
		pos, ok := this.Position[order_.Symbol]
		if !ok {
			panic(order_)
		}
		pos.Qty += o.Qty

		// Generate on Transaction Record

		return true
	} else {
		panic(order_)
	}
}

func (this *Account) handleFilledOrder(order_ *order.Order) bool {
	o, ok := this.PanddingOrders[order_.Id]
	if !ok {
		panic(order_)
	}
	o.Status = order_.Status

	if order_.Side == order.SideBuy {
		// BUY
		pos, ok := this.Position[order_.Symbol]
		if !ok {
			pos := &Position{}
			pos.Symbol = order_.Symbol
			pos.BuyPrice = order_.Price
			pos.Qty = order_.Qty
		} else {
			pos.BuyPrice = (float64(pos.Qty)*pos.BuyPrice + float64(order_.Qty)*order_.Price) / float64(pos.Qty+order_.Qty)
			pos.Qty += order_.Qty
		}

		//trans := Transaction{}
		return true
	} else if order_.Side == order.SideSell {

		// Generate on Transaction Record
		return true
	} else {
		panic(order_)
	}
}
