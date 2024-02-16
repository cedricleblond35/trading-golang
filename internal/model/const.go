package model

type Impact string

const (
	Low    Impact = "1"
	Medium Impact = "2"
	High   Impact = "3"
)

type MarginMode int

const (
	MarginModeForex        MarginMode = 101
	MarginModeCFDLeveraged MarginMode = 102
	MarginModeCFD          MarginMode = 103
)

type Period int

const (
	PERIOD_M1  Period = 1
	PERIOD_M5  Period = 5
	PERIOD_M15 Period = 15
	PERIOD_M30 Period = 30
	PERIOD_H1  Period = 60
	PERIOD_H4  Period = 240
	PERIOD_D1  int    = 1440
	PERIOD_W1  Period = 10080
	PERIOD_MN1 Period = 43200
)

type ProfitMode int

const (
	ProfitModeForex ProfitMode = 5
	ProfitModeCFD   ProfitMode = 6
)

func (p ProfitMode) String() string {
	switch p {
	case ProfitModeForex:
		return "Forex"
	case ProfitModeCFD:
		return "CFD"
	}
	return ""
}

type QuoteID int

const (
	Fixed QuoteID = iota
	Float
	Depth
	Cross
)

func (q QuoteID) String() string {
	switch q {
	case Fixed:
		return "Fixed"
	case Float:
		return "Float"
	case Depth:
		return "Depth"
	case Cross:
		return "Cross"
	}
	return ""
}

type TradeCmd int

const (
	BuyCmd TradeCmd = iota
	SellCmd
	BuyLimitCmd
	SellLimitCmd
	BuyStopCmd
	SellStopCmd
	BalanceCmd
	CreditCmd
)

func (t TradeCmd) String() string {
	switch t {
	case BuyCmd:
		return "buy"
	case SellCmd:
		return "Sell"
	case BuyLimitCmd:
		return "BuyLimit"
	case SellLimitCmd:
		return "SellLimit"
	case BuyStopCmd:
		return "BuyStop"
	case SellStopCmd:
		return "SellStop"
	case BalanceCmd:
		return "Balance"
	case CreditCmd:
		return "Credit"
	}
	return ""
}

type TradeType int

const (
	Open TradeType = iota
	Pending
	Close
	Modify
	Delete
)

func (t TradeType) String() string {
	switch t {
	case Open:
		return "Open"
	case Pending:
		return "Pending"
	case Close:
		return "Close"
	case Modify:
		return "Modify"
	case Delete:
		return "Delete"
	}
	return ""
}

type WeekDay int

const (
	Monday WeekDay = iota + 1
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func (w WeekDay) String() string {
	switch w {
	case Monday:
		return "Monday"
	case Tuesday:
		return "Tuesday"
	case Wednesday:
		return "Wednesday"
	case Friday:
		return "Friday"
	case Saturday:
		return "Saturday"
	case Sunday:
		return "Sunday"
	}
	return ""
}
