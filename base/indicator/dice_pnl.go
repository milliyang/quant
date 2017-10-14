package indicator

import (
	"fmt"
	"quant/base/xbase"
	"time"
)

// Casino Dicing Game
type DicePNL struct {
	Total     float64
	TotalData []float64

	BaseIndicator
}

func NewDicePNL(amt float64) *DicePNL {
	s := &DicePNL{}
	s.Init(0)
	s.Total = amt
	return s
}

// pnl may be negative
func (this *DicePNL) UpdateData(datetime *time.Time, pnl float64) {
	if debug {
		fmt.Println("DicePNL.Append:", pnl)
	}
	this.BaseIndicator.UpdateData(datetime, pnl)

	this.Total += pnl
	this.TotalData = append(this.TotalData, this.Total)
}

// indicator.IIndicator.OnDraw(ICanvas)
func (this *DicePNL) OnDraw(canvas xbase.ICanvas) {
	if len(this.DateTime) == 0 {
		return
	}

	if debug {
		fmt.Println("DicePNL", this.Name, "symbol", this.Symbol, "onDraw")
		fmt.Println("DicePNL onDraw", " start:", this.drawStartIdx, " end:", this.drawEndIdx)
	}

	pnlResult := []string{}
	for i := 0; i < len(this.Data); i++ {
		result := fmt.Sprintf("%d(%+d)", int(this.TotalData[i]), int(this.Data[i]))
		pnlResult = append(pnlResult, result)
	}
	canvas.DrawTextAtPrice(this.DateTime[this.drawStartIdx:this.drawEndIdx], pnlResult, 10.5, this.Color)
}
