package order

import (
	"errors"
	"time"
)

const (
	SideBuy int = 0 + iota
	SideSell
)

const (
	TypeMarket int = 0 + iota
	TypeLimit
	TypeStop
	TypeStopLimit
	TypeTrail
	TypeTrailLimit
	TypeMarketOnClose
)

const (
	StatusPendingNew int = 0 + iota
	StatusNew
	StatusPartiallyFilled
	StatusFilled
	StatusPendingCancel
	StatusCancelled
	StatusExpired
	StatusPendingReplace
	StatusReplaced
	StatusRejected
)

var (
	ErrFieldNoReady = errors.New("bar field not ready")

	orderUniqueSeq = 0
)

func newId() int {
	/*
	 * in the same thread, no need to lock now.
	 */
	orderUniqueSeq++
	return orderUniqueSeq
}

type Order struct {
	Id          int       // required
	Datetime    time.Time //
	Symbol      string    // required
	Side        int       // required
	Type        int       // required
	Status      int       // required
	Amt         float64   // required
	Qty         int
	Price       float64
	StopPrice   float64
	TrailingAmt float64
	AvgPrice    float64
	Account     string
	Text        string

	// Casino Dicing Game
	DiceBetType   int
	DiceBetNumber int
	DiceBetAmt    int
}

func NewMarketOrder(side, qty int, text string) *Order {
	o := &Order{}
	o.Id = newId()
	o.Side = side
	o.Qty = qty
	o.Text = text
	o.Type = TypeMarket
	return o
}

func NewLimitOrder(side, qty int, text string) *Order {
	o := &Order{}
	o.Id = newId()
	o.Side = side
	o.Qty = qty
	o.Text = text
	o.Type = TypeLimit
	return o
}

func NewStopOrder(side, qty int, text string) *Order {
	o := &Order{}
	o.Id = newId()
	o.Side = side
	o.Qty = qty
	o.Text = text
	o.Type = TypeStop
	return o
}

func (this *Order) ToString() string {
	return ""
}
