package provider

import (
	"encoding/json"
	"fmt"
	"github.com/milliyang/dice"
	"os"
	"sort"
	"time"
)

const (
	yearLayout     = "yyyy-mm-dd"
	yyyymmddLayout = "%4d-%02d-%02d"
	floatLayout    = "%d.00"

	ROUND = 600
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

	for i := 0; i < this.Rounds; i++ {

		data := &Datagram{}

		data.Time = fmt.Sprintf(yyyymmddLayout, rollDate.Year(), rollDate.Month(), rollDate.Day())
		data.Symbol = this.Symbols
		data.Event = EventBar
		data.Amount = "1000.00" // in case crash
		data.Volumn = "1000.00"

		diceRoll := casinoRoll()
		// fill it
		data.DiceA = diceRoll.Rolls[0]
		data.DiceB = diceRoll.Rolls[1]
		data.DiceC = diceRoll.Rolls[2]

		//total := data.DiceA + data.DiceB + data.DiceC
		// Hacking
		//data.High = fmt.Sprintf(floatLayout, total)

		var points sort.IntSlice
		points = diceRoll.Rolls
		sort.Sort(points)

		data.Low = fmt.Sprintf(floatLayout, points[0])
		data.Open = fmt.Sprintf(floatLayout, points[1])
		data.Close = fmt.Sprintf(floatLayout, points[2])
		data.High = fmt.Sprintf(floatLayout, points[2])

		outChan <- data

		// Prepare next round
		rollDate = rollDate.AddDate(0, 0, 1)

		//fmt.Println(data)
		//JsonPrint(data)

	}

	// close the channel at last, so that range Chan can finish!!
	close(outChan)
	return nil
}

/*
RollP() generates a new DiceRoll based on the specified parameters.
*/
func casinoRoll() *dice.DiceRoll {
	return dice.RollP(3, 6, 0, false)
}

func JsonPrint(obj interface{}) {
	b, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	os.Stdout.Write(b)
	fmt.Println("")
}
