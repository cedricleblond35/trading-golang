package candlecollector

import (
	"context"
)

type EMA struct {
	Arrondi uint
	Symbol  string
}

func NewEMA(ctx context.Context, arrondi uint, symbol string) *EMA {
	return &EMA{
		Arrondi: arrondi,
		Symbol:  symbol,
	}
}

// EMA = (CLOSE (i) * P) + (EMA (i - 1) * (1 - P))
// (prix de clôture – EMA du jour précédent) × constante pondérant la moyenne mobile exponentielle en décimale + EMA du jour précédent
// CLOSE (i) – prix de clôture de la période actuelle ;
// EMA (i - 1) – valeur de la Moyenne Mobile de la période précédente ;
// P – pourcentage d'utilisation de la valeur du prix.
// α = 2 / (n + 1)
func (ema *EMA) EMA(duration uint) {
	// name := "EMA" + strconv.FormatInt(int64(duration), 10)
	// nameSMA = "SMA" + strconv.FormatInt(int64(duration), 10)
}
