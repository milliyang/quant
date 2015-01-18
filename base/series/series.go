package series

import (
	"fmt"
	"time"
)

const (
	safe_guard = true
	debug      = false
)

type Series interface {
	Keys() []time.Time
	Values() []float64
	Count() float64
	//
	Index(datetime *time.Time) int
	ValueAtIndex(index int) float64
	ValueAtTime(datetime *time.Time) float64

	// event driven interface
	Append(datetime *time.Time, value float64)
}

type FloatSeries struct {
	StartTime       time.Time
	EndTime         time.Time
	DateTime        []time.Time
	Data            []float64
	MapDatetimeData map[int]float64 // map[seconds]value
}

func (this *FloatSeries) Keys() []time.Time {
	if safe_guard {
		all := []time.Time{}
		for _, item := range this.DateTime {
			all = append(all, item)
		}
		return all
	}
	return this.DateTime
}

func (this *FloatSeries) Values() []float64 {
	if safe_guard {
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

func (this *FloatSeries) Index(datetime *time.Time) int {
	for idx, item := range this.DateTime {
		if item.Equal(*datetime) {
			return idx
		}
	}
	return -1
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
	s.Init()
	return s
}

func (this *FloatSeries) Init() {
	this.DateTime = []time.Time{}
	this.Data = []float64{}
	this.MapDatetimeData = map[int]float64{}
}

func (this *FloatSeries) Append(datetime *time.Time, value float64) {
	if debug {
		fmt.Println("FloatSeries.Append")
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
		fmt.Println("can not append duplicate record: %v %v", oldvalue, ok)
		panic(datetime)
	}

	this.MapDatetimeData[sec] = value
	this.DateTime = append(this.DateTime, *datetime)
	this.Data = append(this.Data, value)
}
