package xbase

const (
	CasinoDicingGame = true

	RoundOfPanelInVenetian = 16 // not sure yet

	DiceTrippleAnyOdds    = 24
	DiceTrippleNumberOdds = 150

	HotNumberThreshold  = 11
	CoolNumberThreshold = 6
)

var (
	DiceNumberOddsMap = map[int]int{
		4:  50,
		5:  18,
		6:  14,
		7:  12,
		8:  8,
		9:  6,
		10: 6,
		11: 6,
		12: 6,
		13: 8,
		14: 12,
		15: 14,
		16: 18,
		17: 50,
	}
)
