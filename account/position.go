package account

import ()

type Position struct {
	Symbol       string
	Qty          int
	CurrentPrice float64
	BuyPrice     float64
}

type Transaction struct {
	Id        int
	Symbol    string
	Qty       int
	BuyPrice  float64
	SellPrice float64
	PnL       float64
	Cost      float64
	Charge    float64 // 0.75%
}
