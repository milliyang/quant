package provider

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

func init() {
	if debug {
		fmt.Println("quant/provider pkg init")
	}
}

var (
	ErePathNotExist = errors.New("path not exist")
)

type FileProvider struct {
	Path          string
	Symbols       []string
	cacheDatagram []*Datagram
}

func NewFileProvider(path string) *FileProvider {
	p := FileProvider{}
	p.Path = path
	return &p
}

/*
 * 1. Open Path
 */
func (this *FileProvider) Connect() error {
	if debug {
		fmt.Println("FileProvider check path", this.Path)
	}
	if pathExists(this.Path) {
		if debug {
			fmt.Println("FileProvider check path ok")
		}
		return nil
	}
	return ErePathNotExist
}

// ToDoList
// use 9 bytes Later
// stock  "STO000001"
// ETF    "ETF123456"
// option "OPT000001"
func (this *FileProvider) RequestInstrument(symbols []string) error {
	this.Symbols = symbols

	for _, oneSymbols := range this.Symbols {
		this.loadBarsFromFile(oneSymbols)
	}

	if false {
		dumpDatagram(this.cacheDatagram)
		fmt.Println("len of datagram:", len(this.cacheDatagram))
	}

	sort.Sort(SortableDatagram(this.cacheDatagram))

	// sort them
	if false {
		dumpDatagram(this.cacheDatagram)
		fmt.Println("len of datagram:", len(this.cacheDatagram))
	}
	return nil
}

// This Function Will Block, use go provider.Receive()
func (this *FileProvider) Receive(chan *Datagram) error {
	return nil
}

// Exists reports whether the named file or directory exists.
func pathExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (this *FileProvider) loadBarsFromFile(symbol string) {

	filename := ""
	if len(symbol) < 6 {
		return
	} else if len(symbol) == 6 {
		filename = symbol + ".txt"
	} else {
		// 0123456
		filename = symbol[len(symbol)-6:] + ".txt"
	}

	path := this.Path + string(os.PathSeparator) + filename

	if !pathExists(path) {
		return
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err, path)
		return
	}
	defer f.Close()

	// #date,open,high,low,close,volumn(share/10K),amount(CNY/10K)
	// 2004-01-02,1.00,1.07,1.00,1.06,2191,14946
	// 2004-01-05,1.06,1.13,1.06,1.13,4958,35839

	all, err := ioutil.ReadAll(f)
	if err == nil {
		allLine := string(all)
		lines := strings.Split(allLine, "\n")
		for _, oneLine := range lines {
			if strings.Contains(oneLine, "#") {
				continue
			}
			items := strings.Split(oneLine, ",")

			// 2007-08-08,48.80,58.95,50.80,48.80,35241892.00,1861702656.00,4.00
			if len(items) == 7 {
				datagram := &Datagram{}
				datagram.Symbol = symbol
				datagram.Time = items[0]
				datagram.Open = items[1]
				datagram.High = items[2]
				datagram.Low = items[3]
				datagram.Close = items[4]
				datagram.Volumn = items[5]
				datagram.Amount = items[6]
				this.cacheDatagram = append(this.cacheDatagram, datagram)
			}
		}
	}
}

func dumpDatagram(datagram []*Datagram) {
	for _, one := range datagram {
		fmt.Println(*one)
	}
}

// implements sort.Interface
type SortableDatagram []*Datagram

func (a SortableDatagram) Len() int      { return len(a) }
func (a SortableDatagram) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortableDatagram) Less(i, j int) bool {
	if a[i].Time == a[j].Time {
		if a[i].Symbol == a[j].Symbol {
			if debug {
				fmt.Println("duplicate record", *a[i])
				panic(a[j].Symbol)
			}
			if a[i].Open == a[j].Open {
				return a[i].Close < a[j].Close
			} else {
				return a[i].Open < a[j].Open
			}
		} else {
			return a[i].Symbol < a[j].Symbol
		}
	} else {
		return a[i].Time < a[j].Time
	}
}
