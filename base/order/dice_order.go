package order

const (
	DiceBetTypeBig int = 0 + iota
	DiceBetTypeSmall
	DiceBetTypeEven
	DiceBetTypeSingle
	DiceBetTypeNumber // [4,5,...16,17]
	DiceBetTypeTriple
	DiceBetTypeTripleNumber
	DiceBetTypeFaceNumber // [1,2,3,4,5,6]
)

func NewDiceOrder(betType, money int, text string) *Order {
	o := &Order{}
	o.Id = newId()
	o.Symbol = "DICE"
	o.DiceBetType = betType
	o.DiceBetAmt = money
	o.Text = text
	return o
}

/*
	For: DiceBetTypeNumber, DiceBetTypeTripleNumber, DiceBetTypeFaceNumber
*/
func NewDiceOrderWithNumber(betType, money, number int, text string) *Order {
	o := &Order{}
	o.Id = newId()
	o.Symbol = "DICE"
	o.DiceBetNumber = number
	o.DiceBetType = betType
	o.DiceBetAmt = money
	o.Text = text
	return o
}
