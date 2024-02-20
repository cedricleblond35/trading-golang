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

// df is the dataframe, n is the period, f is the factor; f=3, n=7 are commonly used.
func (s *ST) Calcul(candles []CandleUS100) *ST {
	f := 10
	n := 3

	for i, c := range candles {
		fmt.Println(c)
		fmt.Println(i)
		c := candle{
			High: c.High.Float64,
			Low:  c.Low.Float64,
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
		tr := math.Max(math.Max(hl, hpc), lpc)

		s.candles[i].TR = tr
		// df[i].TR = tr
		if i == n-1 {
			sumTR := 0.0
			for j := 0; j < n-1; j++ {
				sumTR += s.candles[j].TR
			}
			s.candles[i].ATR = sumTR / float64(n-1)
		} else if i > n-1 {
			s.candles[i].ATR = (s.candles[i-1].ATR*(float64(n)-1) + tr) / float64(n)
		}

		s.candles[i].UpperBasic = (s.candles[i].High+s.candles[i].Low)/2 + (float64(f) * s.candles[i].ATR)
		s.candles[i].LowerBasic = (s.candles[i].High+s.candles[i].Low)/2 - (float64(f) * s.candles[i].ATR)
		s.candles[i].UpperBand = s.candles[i].UpperBasic
		s.candles[i].LowerBand = s.candles[i].LowerBasic

		if i >= n {
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
	}

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
