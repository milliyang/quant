package provider

import (
	"errors"
	"fmt"
)

var (
	debug = true

	ErrConnectFail      = errors.New("connect fail")
	ErrPackageInvalid   = errors.New("package invalid")
	ErrConnectionBroken = errors.New("connection broken")
)

func init() {
	if debug {
		fmt.Println("quant/provider pkg init")
	}
}

const (
	EventBarOpen int = 1 + iota
	EventBar
	EventBarSlice
)

// Note:
//   1. Datagram is Sorted in Time < Event < Symbol < Open < Close
//   2. so that everyday, we always get Datagrams in order:
//      >> EventBarOpen, EventBar, EventBar, EventBar ... EventBarSlice
//
type Datagram struct {
	Time   string
	Event  int
	Symbol string
	Open   string
	High   string
	Low    string
	Close  string
	Volumn string
	Amount string

	// Casino Dicing Game
	DiceA int
	DiceB int
	DiceC int
}

type IProvider interface {
	Connect() error
	RequestInstrument([]string) error
	Receive(chan *Datagram) error
}
