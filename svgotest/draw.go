package main

import (
	"fmt"
	"quant/base/indicator"
	"quant/base/series"
	"quant/svgo"
	"time"
)

func DrawData() {
	fSeries := series.NewFloatSeries()

	sma := indicator.NewSMA(fSeries, 16)
	ema := indicator.NewEMA(fSeries, 16)
	smaOfSma := indicator.NewSMA(sma, 16)
	macd := indicator.NewMACD(fSeries, 16, 26, 9)

	now := time.Now()
	var i64 float64
	for i := 0; i < 10; i++ {
		i64 = float64(i)
		now = now.Add(time.Second)
		fSeries.Append(&now, i64)
	}

	fmt.Println("?")
	printslice("orignial:", fSeries.Values())
	printslice("sma6:", sma.Values())
	printslice("sma6 of sma6", smaOfSma.Values())
	printslice("ema", ema.Values())

	fmt.Println("macd:")
	printslice("short ema:", macd.ShortEmaValues())
	printslice("long ema:", macd.LongEmaValues())
	printslice("diff:", macd.DiffValues())
	printslice("dem:", macd.DemValues())
	printslice("osc:", macd.OscValues())

	//
	fmt.Println("draw fSeries")
	svgo.Draw(fSeries)
	svgo.Draw(sma)
	svgo.Draw(smaOfSma)
	svgo.Draw(ema)
	svgo.Draw(macd)

	fmt.Println("dodrawing")
	svgo.TestDoDrawing()
}

func printslice(tag string, values []float64) {
	if len(tag) < 7 {
		fmt.Print(tag, "\t\t[")
	} else {
		fmt.Print(tag, "\t[")
	}
	for _, one := range values {
		fmt.Printf("%0.3f ", one)
	}
	fmt.Print("]\n")
}
