package model

import (
	"database/sql"

	"trading/internal/config"
)

type (
	candleUS100 struct {
		ID     int           `gorm:"type:bigint(20);column:id;autoIncrement"`
		Ctm    sql.NullInt32 `gorm:"type:bigint(20);column:ctm"`
		Close  float32       `gorm:"type:decimal(13,9);column:close"`
		High   float32       `gorm:"type:decimal(13,9);column:high"`
		Low    float32       `gorm:"type:decimal(13,9);column:low"`
		Open   float32       `gorm:"type:decimal(13,9);column:open"`
		Vol    sql.NullInt32 `gorm:"type:int;column:vol"`
		Period uint16        `gorm:"type:int;column:period"`
	}
)

func NewCandle() *candleUS100 {
	return &candleUS100{}
}

func (o *candleUS100) Database() string {
	return config.Database
}

func (o *candleUS100) TableName() string {
	return "candleUS100"
}

// func (o *Candle) ListFromCharacteristics(ctx context.Context, query string, characteristics ...any) ([]*Candle, error) {
// 	var candle []Candle
// 	result :=
// }
