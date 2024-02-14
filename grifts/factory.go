package grifts

import (
	"strconv"
	"time"

	user "trading/internal/model"
)

func CreateFakeUser() *user.User{
	t := time.Now()
	user := user.NewUser()
	user.Name = "name" + strconv.FormatInt(time.Now().UnixNano(), 36)
	user.Email = user.Name+"@truc.com"
	user.UserBroker = user.Name+"@userbroker.com"
	user.PassBrocker = "pass"+strconv.FormatInt(time.Now().UnixNano(), 36)
	user.LastConnection = t
	user.CreatedAt = t

	return user
}
