package provider

import (
	"fmt"
	"github.com/milliyang/dice"
	"github.com/milliyang/dice/utils"
	"sort"
	"time"
)

const (
	yearLayout     = "yyyy-mm-dd"
	yyyymmddLayout = "%4d-%02d-%02d"
	floatLayout    = "%d.00"

	ROUND = 2 * 60 * 8 // 2 * 60 minute * 8hour

	FACES       = 6
	NUM_OF_DICE = 3
)

var (
	// [Important!!!]
	// TIME_ORIGIN is related with  quant/canvas/TIME_ORIGIN
	// and may affect drawing.
	// it's better not change
	TIME_ORIGIN = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.Local)

	rollDate = TIME_ORIGIN //time.Now()
)

type DiceProvider struct {
	Rounds        int
	Symbols       string
	cacheDatagram []*Datagram `json:"-" `
}

func NewDiceProvider() *DiceProvider {
	p := DiceProvider{}
	p.Rounds = ROUND
	p.Symbols = "DICE"
	return &p
}

/*
 * 1. Open Path
 */
func (this *DiceProvider) Connect() error {
	if debug {
		fmt.Println("DiceProvider rounds:", this.Rounds)
	}
	return nil
}

func (this *DiceProvider) RequestInstrument(symbols []string) error {
	return nil
}

// This Function Will Block, use go provider.Receive()
func (this *DiceProvider) Receive(outChan chan *Datagram) error {

	var allDiceRolls []*dice.DiceRoll

	for i := 0; i < this.Rounds; i++ {

		data := &Datagram{}

		data.Time = fmt.Sprintf(yyyymmddLayout, rollDate.Year(), rollDate.Month(), rollDate.Day())
		data.Symbol = this.Symbols
		data.Event = EventBar
		data.Amount = "1000.00" // in case crash
		data.Volumn = "1000.00"

		diceRoll := utils.CasinoRoll()

		// sort
		var points sort.IntSlice
		points = diceRoll.Rolls
		sort.Sort(points)

		//total := data.DiceA + data.DiceB + data.DiceC
		// Hacking
		// fill it
		data.DiceA = points[0]
		data.DiceB = points[1]
		data.DiceC = points[2]

		data.Open = fmt.Sprintf(floatLayout, points[1])
		data.Close = fmt.Sprintf(floatLayout, points[0]+points[1]+points[2])

		data.High = fmt.Sprintf(floatLayout, points[0]+points[1]+points[2])
		data.Low = fmt.Sprintf(floatLayout, 1)

		outChan <- data

		// Prepare next round
		rollDate = rollDate.AddDate(0, 0, 1)

		//fmt.Println(data)
		//utils.JsonPrint(data)

		allDiceRolls = append(allDiceRolls, diceRoll)
	}

	// close the channel at last, so that range Chan can finish!!
	close(outChan)
	return nil
}
