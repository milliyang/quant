package series

import (
	"fmt"
	"quant/base/xbase"
	"time"
)

const (
	enable_safe_mode = true
	debug            = false
)

/*
	implement interface:
	1. xbase.ISeries
	2. IIndicator


*/
type FloatSeries struct {
	Name            string
	Symbol          string
	StartTime       time.Time
	EndTime         time.Time
	DateTime        []time.Time
	Data            []float64
	MapDatetimeData map[int]float64 // map[seconds]value

	InnerParent xbase.ISeries
	InnerChilds []xbase.ISeries
}

func (this *FloatSeries) Keys() []time.Time {
	if enable_safe_mode {
		all := []time.Time{}
		for _, item := range this.DateTime {
			all = append(all, item)
		}
		return all
	}
	return this.DateTime
}

func (this *FloatSeries) Values() []float64 {
	if enable_safe_mode {
		all := []float64{}
		for _, item := range this.Data {
			all = append(all, item)
		}
		return all
	}
	return this.Data
}

func (this *FloatSeries) Count() int {
	return len(this.Data)
}

func (this *FloatSeries) Contains(datetime *time.Time) bool {
	if this.Index(datetime) == -1 {
		return false
	} else {
		return true
	}
}

func (this *FloatSeries) Index(datetime *time.Time) int {
	for idx, item := range this.DateTime {
		if item.Equal(*datetime) {
			return idx
		}
	}
	return -1
}

func (this *FloatSeries) ValueAtTime(datetime *time.Time) float64 {

	idx := this.Index(datetime)
	if idx >= 0 {
		return this.ValueAtIndex(idx)
	} else {
		fmt.Println("ValueAtTime invalid datetime: %v %v", datetime, idx)
		panic(datetime)
		return -1
	}
}

func (this *FloatSeries) ValueAtIndex(index int) float64 {
	if index >= len(this.Data) || index < 0 {
		fmt.Println("OutOfArray: %v %v", len(this.Data), index)
		panic(index)
	}
	return this.Data[index]
}

func NewFloatSeries() *FloatSeries {
	s := &FloatSeries{}
	s.Init(nil)
	return s
}

func (this *FloatSeries) Init(parent xbase.ISeries) {
	this.DateTime = []time.Time{}
	this.Data = []float64{}
	this.MapDatetimeData = map[int]float64{}
	this.InnerParent = parent

	// [GoBug Tag00001]??
	//
	// It seems we can't do it in anonymous struct. leave it to indicator.SMA
	//
	// if parent != nil {
	// 	parent.AddChild(this)
	// }
}

func (this *FloatSeries) Match(symbol string) bool {
	if this.Symbol == symbol {
		return true
	} else {
		return false
	}
}

func (this *FloatSeries) Now() time.Time {
	return this.EndTime
}

func (this *FloatSeries) AddChild(child xbase.ISeries) {
	this.InnerChilds = append(this.InnerChilds, child)
}

func (this *FloatSeries) Append(datetime *time.Time, value float64) {
	if debug {
		fmt.Println("FloatSeries.Append:", value)
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

	this.MapDatetimeData[sec] = value
	this.DateTime = append(this.DateTime, *datetime)
	this.Data = append(this.Data, value)

	// [GoBug Tag00001]??
	for _, child := range this.InnerChilds {
		child.Append(datetime, value)
	}
}

// indicator.IIndicator.OnMeasure

// num of table, X cordirate [datetime count], Y cordirate [min, max, num, percent]
func (this *FloatSeries) OnMeasure() (int, int, float64, float64, int, float64) {

	// put it here now.
	var min, max float64
	var num int

	for idx, value := range this.Data {
		num = idx
		if value < min {
			min = value
		}
		if value > max {
			value = max
		}
	}
	return 1, len(this.DateTime), min, max, num, 100
}

// indicator.IIndicator.OnDraw(ICanvas)
func (this *FloatSeries) OnDraw(canvas xbase.ICanvas) {
	fmt.Println(this.Name, "symbol", this.Symbol, "onDraw")

	canvas.DrawLine(this.DateTime, this.Data, 1)
	// canvas.DrawBar(table,  []time.Time, []bar.Bar, color)

	// canvas.DrawBuy(table,  []time.Time, []float64,color)
	// canvas.DrawSell(table, []time.Time, []float64,color)
	// canvas.DrawSpark(table,[]time.Time, []float64,color)
	// canvas.DrawShit(table, []time.Time, []float64,color)
}
