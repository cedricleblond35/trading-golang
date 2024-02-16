package xtb

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type XTB struct {
	URL    url.URL
	Conn   *websocket.Conn
	Status bool
}

type login struct {
	Command   string `json:"command"`
	Arguments struct {
		UserID   string `json:"userId"`
		Password string `json:"password"`
	} `json:"arguments"`
}
type argsLogin struct {
	UserID   string `json:"userId"`
	Password string `json:"password"`
}

type responseConn struct {
	Status          bool   `json:"status"`
	StreamSessionID string `json:"streamSessionId"`
}

func NewXTB() *XTB {
	return &XTB{}
}

func (xtb *XTB) Connection(host, path string) error {
	if host == "" {
		fmt.Println("error host")
	}
	if path == "" {
		fmt.Println("error path")
	}
	u := url.URL{
		Scheme: "wss",
		Host:   host,
		Path:   path,
	}
	log.Printf("connecting to %s", u.String())

	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		return err
	}
	xtb.Conn = c

	return nil
}

func (xtb *XTB) Login(user, password string) error {
	jsonData, err := json.Marshal(login{
		Command: "login",
		Arguments: argsLogin{
			UserID:   user,
			Password: password,
		},
	})
	if err != nil {
		return err
	}

	if err := xtb.Conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		return err
	}
	_, message, err := xtb.Conn.ReadMessage()
	if err != nil {
		return err
	}

	var rc responseConn
	err = json.Unmarshal(message, &rc)
	if err != nil {
		return err
	}
	xtb.Status = rc.Status
	fmt.Println("connected to xtb:", rc.Status)

	return nil
}
