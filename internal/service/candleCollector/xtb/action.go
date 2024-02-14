package xtb

type Login struct {
	Command   string `json:"command"`
	Arguments struct {
		UserID   string `json:"userId"`
		Password string `json:"password"`
	} `json:"arguments"`
}
type ArgsLogin struct {
	UserID   string `json:"userId"`
	Password string `json:"password"`
}


type ResponseConn struct {
	Status          bool   `json:"status"`
	StreamSessionID string `json:"streamSessionId"`
}

///////////////////////////////////////////////////////
type GetChartRangeRequest struct {
	Command   string `json:"command"`
	Arguments struct {
		Info struct {
			End    int64  `json:"end"`
			Period int    `json:"period"`
			Start  int64  `json:"start"`
			Symbol string `json:"symbol"`
			Ticks  int    `json:"ticks"`
		} `json:"info"`
	} `json:"arguments"`
}
type ArgsGetChartRangeRequest struct {
	Info struct {
		End    int64  `json:"end"`
		Period int    `json:"period"`
		Start  int64  `json:"start"`
		Symbol string `json:"symbol"`
		Ticks  int    `json:"ticks"`
	} `json:"info"`
}
type ArgsInfoGetChartRangeRequest struct {
	End    int64  `json:"end"`
	Period int    `json:"period"`
	Start  int64  `json:"start"`
	Symbol string `json:"symbol"`
	Ticks  int    `json:"ticks"`
}


// /////////////////////////////////////////////////////////////////////
type GetChartLastRequest struct {
	Command   string `json:"command"`
	Arguments struct {
		Info struct {
			Period int    `json:"period"`
			Start  int64  `json:"start"`
			Symbol string `json:"symbol"`
		} `json:"info"`
	} `json:"arguments"`
}
type ArgsGetChartLastRequest struct {
	Info struct {
		Period int    `json:"period"`
		Start  int64  `json:"start"`
		Symbol string `json:"symbol"`
	} `json:"info"`
}
type ArgsInfogetChartLastRequest struct {
	Period int    `json:"period"`
	Start  int64  `json:"start"`
	Symbol string `json:"symbol"`
}
///////////////////////////////////////////////////////////////////////
type ResponseChartLastRequest struct {
	Status     bool `json:"status"`
	ReturnData struct {
		Digits    int `json:"digits"`
		RateInfos []struct {
			Close     float64 `json:"close"`
			Ctm       int32   `json:"ctm"`
			CtmString string  `json:"ctmString"`
			High      float64 `json:"high"`
			Low       float64 `json:"low"`
			Open      float64 `json:"open"`
			Vol       float64 `json:"vol"`
		} `json:"rateInfos"`
	} `json:"returnData"`
}
