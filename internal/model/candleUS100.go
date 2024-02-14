package model

import (
	"database/sql"

	"trading/internal/config"
)

type (
	CandleUS100 struct {
		ID     int           `gorm:"type:bigint(20);column:id;autoIncrement"`
		Ctm    sql.NullInt64 `gorm:"type:bigint(20);column:ctm;index,unique,composite:candleuniq_id"`
		Close  sql.NullInt32 `gorm:"type:int(8);column:close"`
		High   sql.NullInt32 `gorm:"type:int(8);column:high"`
		Low    sql.NullInt32 `gorm:"type:int(8);column:low"`
		Open   sql.NullInt32 `gorm:"type:int(8);column:open"`
		Vol    sql.NullInt32 `gorm:"type:int;column:vol"`
		Period sql.NullInt16 `gorm:"type:int;column:period;index,unique,composite:candleuniq_id"`
	}
)

func NewCandle() *CandleUS100 {
	return &CandleUS100{}
}

func (o *CandleUS100) Database() string {
	return config.Database
}

func (o *CandleUS100) TableName() string {
	return "candleUS100"
}