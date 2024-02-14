package model

type Pivot struct {
	R4 float64
	R3 float64
	R2 float64
	R1 float64
	P  float64
	S1 float64
	S2 float64
	S3 float64
	S4 float64
}

func NewPivot() *Pivot {
	return &Pivot{}
}

func (pivot *Pivot) PivotClassic(high, low, close float64) *Pivot {
	pivot.P = (high + low + close) / 3
	pivot.S1 = (pivot.P * 2) - high
	pivot.S2 = pivot.P - (high - low)
	pivot.S3 = pivot.P - (high-low)*2
	pivot.R1 = (pivot.P * 2) - low
	pivot.R2 = pivot.P + (high - low)
	pivot.R3 = pivot.P + (high-low)*2

	return pivot
}

func (pivot *Pivot) PivotFibonacci(high, low, close float64) *Pivot {
	pivot.P = (high + low + close) / 3
	pivot.S1 = pivot.P - 0.382*(high-low)
	pivot.S2 = pivot.P - 0.618*(high-low)
	pivot.S3 = pivot.P - 1.000*(high-low)
	pivot.R1 = pivot.P + 0.382*(high-low)
	pivot.R2 = pivot.P + 0.618*(high-low)
	pivot.R3 = pivot.P + 1.000*(high-low)

	return pivot
}

func (pivot *Pivot) PivotWoodie(high, low, close float64) *Pivot {
	pivot.P = (high + low + 2*close) / 4
	pivot.S1 = (2 * pivot.P) - high
	pivot.S2 = pivot.P - (high - low)
	pivot.R1 = (2 * pivot.P) - low
	pivot.R2 = pivot.P + (high - low)

	return pivot
}

func (pivot *Pivot) PivotCamarilla(high, low, close float64) *Pivot {
	pivot.P = (high + low + close) / 3
	pivot.S1 = close - (high-low)*1.1/12
	pivot.S2 = close - (high-low)*1.1/6
	pivot.S3 = close - (high-low)*1.1/4
	pivot.S4 = close - (high-low)*1.1/2
	pivot.R4 = (high-low)*1.1/2 + close
	pivot.R3 = (high-low)*1.1/4 + close
	pivot.R2 = (high-low)*1.1/6 + close
	pivot.R1 = (high-low)*1.1/12 + close

	return pivot
}
