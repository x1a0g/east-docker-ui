package API

type ApiResponseObject struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (a *ApiResponseObject) Fail(code string, Msg string) {
	a.Code = code
	a.Msg = Msg
}
func (a *ApiResponseObject) Success(code string, Msg string) {
	a.Code = code
	a.Msg = Msg
}

func (a *ApiResponseObject) Success4data(data interface{}) {
	a.Code = SUCCESS.GetCode()
	a.Msg = SUCCESS.GetName()
	a.Data = data
}
