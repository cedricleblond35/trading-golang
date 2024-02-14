package model

type TradeTransInfo struct {
	Cmd           TradeCmd  `json:"cmd"`
	CustomComment string    `json:"customComment"`
	Expiration    int       `json:"expiration"`
	Offset        int       `json:"offset"`
	Order         int       `json:"order"`
	Price         float64   `json:"price"`
	Sl            float64   `json:"sl"`
	Symbol        string    `json:"symbol"`
	Tp            float64   `json:"tp"`
	Type          TradeType `json:"type"`
	Volume        float64   `json:"volume"`
}

type Order struct {
	Order int `json:"order"`
}

type Trade struct {
	ClosePrice       float64  `json:"close_price"`
	CloseTime        int      `json:"close_time"`
	CloseTimeString  string   `json:"close_timeString"`
	Closed           bool     `json:"closed"`
	Cmd              TradeCmd `json:"cmd"`
	Comment          string   `json:"comment"`
	Commission       float64  `json:"commission"`
	CustomComment    string   `json:"customComment"`
	Digits           int      `json:"digits"`
	Expiration       int      `json:"expiration"`
	ExpirationString string   `json:"expirationString"`
	MarginRate       float64  `json:"margin_rate"`
	Offset           int      `json:"offset"`
	OpenPrice        float64  `json:"open_price"`
	OpenTime         int      `json:"open_time"`
	OpenTimeString   string   `json:"open_timeString"`
	Order            int      `json:"order"`
	Order2           int      `json:"order2"`
	Position         int      `json:"position"`
	Profit           float64  `json:"profit"`
	Sl               float64  `json:"sl"`
	Storage          float64  `json:"storage"`
	Symbol           string   `json:"symbol"`
	TimeStamp        int      `json:"timestamp"`
	Tp               float64  `json:"tp"`
	Volume           float64  `json:"volume"`
}
