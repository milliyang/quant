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

type Datagram struct {
	Symbol string
	Time   string
	Open   string
	High   string
	Low    string
	Close  string
	Volumn string
	Amount string
}

type IProvider interface {
	Connect() error
	RequestInstrument([]string) error
	Receive(chan *Datagram) error
}
