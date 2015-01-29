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
	Color           int
	MapDatetimeData map[int]float64 // map[seconds]value

	InnerParent xbase.ISeries
	InnerChilds []xbase.ISeries

	// use between OnMeasure() && OnDraw()
	drawStart    time.Time
	drawEnd      time.Time
	drawStartIdx int
	drawEndIdx   int
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
	s.Init(nil, 0)
	return s
}

func (this *FloatSeries) Init(parent xbase.ISeries, color int) {
	this.DateTime = []time.Time{}
	this.Data = []float64{}
	this.MapDatetimeData = map[int]float64{}
	this.InnerParent = parent
	this.Color = color

	// [GoBug Tag00001]??
	//
	// It seems we can't do it in anonymous struct. leave it to indicator.SMA
	//
	// if parent != nil {
	// 	parent.AddChild(this)
	// }
}

func (this *FloatSeries) InitWorkaround(parent xbase.ISeries, color int) {

	this.Init(parent, color)
	if parent != nil {
		parent.AddChild(this)
	}
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

// X cordirate [datetime count], Y cordirate [min, max, num, percent]
func (this *FloatSeries) OnMeasure(start, end time.Time) (int, float64, float64, int, float64) {
	if debug {
		fmt.Println("FloatSeries ", this.Name, "symbol", this.Symbol, "OnMeasure")
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
		fmt.Println("FloatSeries OnMeasure", num, min, max, 100, 100)
	}

	return num, min, max, 100, 100
}

// indicator.IIndicator.OnDraw(ICanvas)
func (this *FloatSeries) OnDraw(canvas xbase.ICanvas) {
	if len(this.DateTime) == 0 {
		return
	}

	if debug {
		fmt.Println("FloatSeries", this.Name, "symbol", this.Symbol, "onDraw")
		fmt.Println("FloatSeries onDraw", " start:", this.drawStartIdx, " end:", this.drawEndIdx)
	}
	canvas.DrawLine(this.DateTime[this.drawStartIdx:this.drawEndIdx], this.Data[this.drawStartIdx:this.drawEndIdx], this.Color)

	// canvas.DrawBuy(table,  []time.Time, []float64,color)
	// canvas.DrawSell(table, []time.Time, []float64,color)
	// canvas.DrawSpark(table,[]time.Time, []float64,color)
	// canvas.DrawShit(table, []time.Time, []float64,color)
}

func (this *FloatSeries) OnLayout() []time.Time {
	if debug {
		fmt.Println("FloatSeries ", this.Name, "symbol", this.Symbol, "OnLayout")
	}
	/*
	 * Note:
	 *  For a = [0,1,2,3]:
	 *    a[0:1] means [0]
	 *    a[0:4] means [0,1,2,3]
	 */
	return this.DateTime[this.drawStartIdx:this.drawEndIdx]
}
