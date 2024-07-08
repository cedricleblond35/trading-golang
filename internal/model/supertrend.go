package model

import (
	"fmt"
	"math"
)

type ST struct {
	multiplicateur float64
	periode        int
	duration       int
	candles        []candle
}

type candle struct {
	Date       string
	Close      float64
	High       float64
	Low        float64
	Open       float64
	TR         float64
	ATR        float64
	UpperBand  float64
	LowerBand  float64
	UpperBasic float64
	LowerBasic float64
	SuperTrend float64
}

func NewSupertrend() *ST {
	return &ST{}
}

func absDiff(a, b float64) float64 {
	return math.Abs(a - b)
}

func (s *ST) Get() []candle {
	return s.candles
}

func max(values ...float64) float64 {
	maxVal := values[0]
	for _, v := range values {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}

// func calculateTR(df *DataFrame) {
// 	length := len(df.high)
// 	df.TR = make([]float64, length)

// 	for i := 0; i < length; i++ {
// 		high := df.high[i]
// 		low := df.low[i]
// 		var prevClose float64
// 		if i > 0 {
// 			prevClose = df.close[i-1]
// 		}

// 		hl := absDiff(high, low)
// 		hpc := absDiff(high, prevClose)
// 		lpc := absDiff(low, prevClose)

// 		df.TR[i] = max(hl, hpc, lpc)
// 	}
// }

// df is the dataframe, n is the period, f is the factor; f=3, n=7 are commonly used.
func (s *ST) Calcul(candles []CandleUS100) *ST {
	f := 5  // coefficient permettant de pondérer la volatilité mesurée par l'ATR. La formule standard utilise un coefficient de 3
	n := 10 // période de calcul de l'ATR, typiquement 10

	for i, c := range candles {
		fmt.Println("----------------------------------------")
		c := candle{
			High:  c.High.Float64,
			Low:   c.Low.Float64,
			Date:  c.Date,
			Close: c.Close.Float64,
		}
		s.candles = append(s.candles, c)

		// s.candles[i].High = c.High.Float64
		// s.candles[i].Low = c.Low.Float64
		var prevClose float64
		if i == 0 {
			continue
		}
		prevClose = s.candles[i-1].Close

		hl := math.Abs(s.candles[i].High - s.candles[i].Low)
		hpc := math.Abs(s.candles[i].High - prevClose)
		lpc := math.Abs(s.candles[i].Low - prevClose)

		// arrondi
		tr := math.Round(math.Max(math.Max(hl, hpc), lpc)*100) / 100
		s.candles[i].TR = tr

		if i == n {
			sumTR := 0.0
			for j := 1; j < n+1; j++ {
				fmt.Printf("-------------->Premier J:%+v, TR:%v\n", j, s.candles[j].TR)
				sumTR += s.candles[j].TR
				fmt.Println("sumTR", sumTR)
			}
			s.candles[i].ATR = sumTR / float64(n)
			fmt.Println("i:", i, "date:", s.candles[i].Date, "    prevClose", prevClose, "HL", hl, "  hpc:", hpc, " lpc:", lpc, " tr", tr, " ATR:", s.candles[i].ATR)
			fmt.Println("******** Premier ATR calculer")
		} else if i > n {
			s.candles[i].ATR = (s.candles[i-1].ATR*(float64(n)-1) + tr) / float64(n)
			fmt.Println("******** ATR calculé suivant")
			fmt.Println("i:", i, "date:", s.candles[i].Date, "    prevClose", prevClose, "HL", hl, "  hpc:", hpc, " lpc:", lpc, " tr", tr, " ATR:", s.candles[i].ATR)

		}

		// Calculation of SuperTrend
		s.candles[i].UpperBasic = (s.candles[i].High+s.candles[i].Low)/2 + (float64(f) * s.candles[i].ATR)
		s.candles[i].LowerBasic = (s.candles[i].High+s.candles[i].Low)/2 - (float64(f) * s.candles[i].ATR)
		s.candles[i].UpperBand = s.candles[i].UpperBasic
		s.candles[i].LowerBand = s.candles[i].LowerBasic

	}

	for i := n; i < len(s.candles); i++ {
		if s.candles[i-1].Close <= s.candles[i-1].UpperBand {
			s.candles[i].UpperBand = math.Min(s.candles[i].UpperBasic, s.candles[i-1].UpperBand)
		} else {
			s.candles[i].UpperBand = s.candles[i].UpperBasic
		}

		if s.candles[i-1].Close >= s.candles[i-1].LowerBand {
			s.candles[i].LowerBand = math.Max(s.candles[i].LowerBasic, s.candles[i-1].LowerBand)
		} else {
			s.candles[i].LowerBand = s.candles[i].LowerBasic
		}
	}

	fmt.Println("UpperBasic:", s.candles[len(candles)-10].UpperBasic)
	fmt.Println("LowerBasic:", s.candles[len(candles)-10].LowerBasic)

	for i := range s.candles {
		if i >= n {
			if s.candles[i-1].SuperTrend == s.candles[i-1].UpperBand && s.candles[i].Close <= s.candles[i].UpperBand {
				s.candles[i].SuperTrend = s.candles[i].UpperBand
			} else if s.candles[i-1].SuperTrend == s.candles[i-1].UpperBand && s.candles[i].Close > s.candles[i].UpperBand {
				s.candles[i].SuperTrend = s.candles[i].LowerBand
			} else if s.candles[i-1].SuperTrend == s.candles[i-1].LowerBand && s.candles[i].Close >= s.candles[i].LowerBand {
				s.candles[i].SuperTrend = s.candles[i].LowerBand
			} else if s.candles[i-1].SuperTrend == s.candles[i-1].LowerBand && s.candles[i].Close < s.candles[i].LowerBand {
				s.candles[i].SuperTrend = s.candles[i].UpperBand
			}
		}
	}

	return s
}
