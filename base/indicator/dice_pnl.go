package indicator

import (
	"fmt"
	"quant/base/xbase"
	"time"
)

const (
	enable_safe_mode = true
)

/*
	implement interface:
	1. IIndicator

*/

// Casino Dicing Game
type DicePnl struct {
	Name      string
	Symbol    string
	StartTime time.Time
	EndTime   time.Time
	DateTime  []time.Time
	Data      []float64
	Color     int

	TotalData       []float64
	MapDatetimeData map[int]float64 // map[seconds]value
	Total           float64

	// use between OnMeasure() && OnDraw()
	drawStart    time.Time
	drawEnd      time.Time
	drawStartIdx int
	drawEndIdx   int
}

func (this *DicePnl) Keys() []time.Time {
	if enable_safe_mode {
		all := []time.Time{}
		for _, item := range this.DateTime {
			all = append(all, item)
		}
		return all
	}
	return this.DateTime
}

func (this *DicePnl) Values() []float64 {
	if enable_safe_mode {
		all := []float64{}
		for _, item := range this.Data {
			all = append(all, item)
		}
		return all
	}
	return this.Data
}

func (this *DicePnl) Count() int {
	return len(this.Data)
}

func (this *DicePnl) Contains(datetime *time.Time) bool {
	if this.Index(datetime) == -1 {
		return false
	} else {
		return true
	}
}

func (this *DicePnl) Index(datetime *time.Time) int {
	for idx, item := range this.DateTime {
		if item.Equal(*datetime) {
			return idx
		}
	}
	return -1
}

func (this *DicePnl) ValueAtTime(datetime *time.Time) float64 {

	idx := this.Index(datetime)
	if idx >= 0 {
		return this.ValueAtIndex(idx)
	} else {
		fmt.Println("ValueAtTime invalid datetime: %v %v", datetime, idx)
		panic(datetime)
		return -1
	}
}

func (this *DicePnl) ValueAtIndex(index int) float64 {
	if index >= len(this.Data) || index < 0 {
		fmt.Println("OutOfArray: %v %v", len(this.Data), index)
		panic(index)
	}
	return this.Data[index]
}

func NewDicePnl(amt float64) *DicePnl {
	s := &DicePnl{}
	s.Init(0)
	s.Total = amt
	return s
}

func (this *DicePnl) Init(color int) {
	this.DateTime = []time.Time{}
	this.Data = []float64{}
	this.MapDatetimeData = map[int]float64{}
	this.Color = color
}

func (this *DicePnl) Match(symbol string) bool {
	if this.Symbol == symbol {
		return true
	} else {
		return false
	}
}

func (this *DicePnl) Now() time.Time {
	return this.EndTime
}

// pnl may be negative

func (this *DicePnl) UpdatePnl(datetime *time.Time, pnl float64) {
	if debug {
		fmt.Println("DicePnl.Append:", pnl)
	}

	if len(this.Data) == 0 {
		this.StartTime = *datetime
		this.EndTime = *datetime
	} else {
		this.EndTime = *datetime
	}

	sec := int(datetime.Unix())
	oldvalue, ok := this.MapDatetimeData[sec]
	if ok {
		// can not append duplicate record
		// second is min unit "now".
		fmt.Println("can not append duplicate record: %v %v", oldvalue, ok)
		panic(datetime)
	}

	this.Total += pnl
	this.MapDatetimeData[sec] = pnl

	this.DateTime = append(this.DateTime, *datetime)
	this.Data = append(this.Data, pnl)
	this.TotalData = append(this.TotalData, this.Total)
}

// indicator.IIndicator.OnMeasure

// X cordirate [datetime count], Y cordirate [min, max, num, percent]
func (this *DicePnl) OnMeasure(start, end time.Time) (int, float64, float64, int, float64) {
	if debug {
		fmt.Println("DicePnl ", this.Name, "symbol", this.Symbol, "OnMeasure")
	}
	this.drawStart = start
	this.drawEnd = end

	this.drawStartIdx = 0
	this.drawEndIdx = 0

	// put it here now.
	var min, max float64
	var num int

	min = 10000
	max = -10000

	for idx, oneTime := range this.DateTime {
		if oneTime.Before(this.drawStart) {
			continue
		} else if oneTime.After(this.drawEnd) {
			continue
		}

		if num == 0 {
			this.drawStartIdx = idx
		}
		this.drawEndIdx = idx + 1
		num++

		value := this.Data[idx]
		if value < min {
			min = value
		}
		if value > max {
			value = max
		}
	}

	if debug {
		fmt.Println("DicePnl OnMeasure", num, min, max, 100, 100)
	}

	return num, min, max, 100, 100
}

// indicator.IIndicator.OnDraw(ICanvas)
func (this *DicePnl) OnDraw(canvas xbase.ICanvas) {
	if len(this.DateTime) == 0 {
		return
	}

	if debug {
		fmt.Println("DicePnl", this.Name, "symbol", this.Symbol, "onDraw")
		fmt.Println("DicePnl onDraw", " start:", this.drawStartIdx, " end:", this.drawEndIdx)
	}

	canvas.DrawTextAtPrice(this.DateTime[this.drawStartIdx:this.drawEndIdx], []string{}, 10.5, this.Color)
}

func (this *DicePnl) OnLayout() []time.Time {
	if debug {
		fmt.Println("DicePnl ", this.Name, "symbol", this.Symbol, "OnLayout")
	}
	/*
	 * Note:
	 *  For a = [0,1,2,3]:
	 *    a[0:1] means [0]
	 *    a[0:4] means [0,1,2,3]
	 */
	return this.DateTime[this.drawStartIdx:this.drawEndIdx]
}
