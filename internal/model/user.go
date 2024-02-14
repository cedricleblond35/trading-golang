package model

import (
	"time"

	"trading/internal/config"
)

type (
	User struct {
		ID             int       `gorm:"type:bigint(20);column:id;autoIncrement"`
		Name           string    `gorm:"type:varchar(20);column:name"`
		Email          string    `gorm:"type:varchar(50);column:email"`
		UserBroker     string    `gorm:"type:varchar(50);column:userbrocker"`
		PassBrocker    string    `gorm:"type:varchar(50);column:passBrocker"`
		LastConnection time.Time `gorm:"type:datetime;column:lastConnection"`
		CreatedAt      time.Time `gorm:"type:datetime;column:created"`
	}
)

func NewUser() *User {
	return &User{}
}

func (o *User) Database() string {
	return config.Database
}

func (o *User) TableName() string {
	return "user"
}
