package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	BAR_OPEN   = 0
	BAR_HIGH   = 1
	BAR_LOW    = 2
	BAR_CLOSE  = 3
	BAR_VOLUMN = 4
	BAR_AMOUNT = 5
	BAR_FACTOR = 6
)

var (
	FMT_FLOAT2 = "%s,%.2f,%.2f,%.2f,%.2f,%.0f,%.0f\n"
	FMT_FLOAT3 = "%s,%.3f,%.3f,%.3f,%.3f,%.0f,%.0f\n"

	NUM_10K = 10000.0
)

// 2014-01-28,264.190,266.849,265.963,262.860,9731806.000,58202184.000,44.327
type Bar struct {
	Date  string   // Trade Date
	Items []string // {OPEN,HIGH,LOW,CLOSE, Volumn,Amount,Factor}
}

// date,open,high,low,close,volumn(share/10K),amount(CNY/10K)
// 2014-01-28,264.190,266.849,265.963,262.860,9731806.000,58202184.00
type NumBar struct {
	Date  string    // Trade Date
	Items []float64 // {OPEN,HIGH,LOW,CLOSE, Volumn,Amount}
}

func (this *Bar) toString() string {
	line := this.Date + "," + strings.Join(this.Items, ",") + "\n"
	return line
}

func (this *Bar) hasFactor() bool {
	if len(this.Items) == 7 {
		return true
	} else {
		return false
	}
}

func (this *Bar) getFactor() float64 {
	newestFactor := this.Items[len(this.Items)-1]
	factor, _ := strconv.ParseFloat(newestFactor, 32)
	return factor
}

func (this *Bar) selfCheck() bool {
	if this.Date == "" {
		return false
	} else if len(this.Date) == 10 {
		//2014-01-28
		return true
	} else {
		return false
	}
}

// Fxxking SINA
// open, high, close, low, volumn, amouunt, factor ->
// open, high, low, close, volumn, amount, factor
func (this *Bar) fxxkingSina() {
	if len(this.Items) == 7 || len(this.Items) == 6 {
		swap := this.Items[3]
		this.Items[3] = this.Items[2]
		this.Items[2] = swap
	}
}

// date,open,high,low,close,volumn(share/10K),amount(CNY/10K)\n
func (this *Bar) toNumBar(factor float64) NumBar {
	numBar := NumBar{}
	numBar.Date = this.Date
	open, _ := strconv.ParseFloat(this.Items[BAR_OPEN], 32)
	high, _ := strconv.ParseFloat(this.Items[BAR_HIGH], 32)
	low, _ := strconv.ParseFloat(this.Items[BAR_LOW], 32)
	close_, _ := strconv.ParseFloat(this.Items[BAR_CLOSE], 32)
	// 2007-08-08,48.80,58.95,50.80,48.80,35241892.00,1861702656.00
	volumn, _ := strconv.ParseFloat(this.Items[BAR_VOLUMN], 32)
	amount, _ := strconv.ParseFloat(this.Items[BAR_AMOUNT], 32)

	numBar.Items = append(numBar.Items, open/factor)
	numBar.Items = append(numBar.Items, high/factor)
	numBar.Items = append(numBar.Items, low/factor)
	numBar.Items = append(numBar.Items, close_/factor)
	numBar.Items = append(numBar.Items, volumn/NUM_10K)
	numBar.Items = append(numBar.Items, amount/NUM_10K)
	return numBar
}

func (this *NumBar) toString() string {
	line := fmt.Sprintf(FMT_FLOAT2, this.Date,
		this.Items[BAR_OPEN],
		this.Items[BAR_HIGH],
		this.Items[BAR_LOW],
		this.Items[BAR_CLOSE],
		this.Items[BAR_VOLUMN],
		this.Items[BAR_AMOUNT])
	return line
}

// 创业板的代码是300打头的股票代码
// 沪市A股
// 沪市A股的代码是以600、601或603打头
// 沪市B股
// 沪市B股的代码是以900打头
// 深市A股
// 深市A股的代码是以000打头
// 中小板
// 中小板的代码是002打头
// 深圳B股
// 深圳B股的代码是以200打头
