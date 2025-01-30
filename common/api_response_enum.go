package API

type ApiResponseEnum string

const (
	AIREADY_EXISTS   ApiResponseEnum = "-4:密码重复"
	SUCCESS          ApiResponseEnum = "0:成功"
	ERROR_PARAM      ApiResponseEnum = "-2:参数错误"
	ERROR_DATA_EMPTY ApiResponseEnum = "-3:无数据"
	FAIL             ApiResponseEnum = "-1:失败"
	LOGIN_FAIL       ApiResponseEnum = "1001:用户名密码错误"
)

func (a ApiResponseEnum) GetCode() string {
	codeMap := map[ApiResponseEnum]string{
		AIREADY_EXISTS:   "-4",
		SUCCESS:          "0",
		ERROR_PARAM:      "-2",
		ERROR_DATA_EMPTY: "-3",
		FAIL:             "-1",
		LOGIN_FAIL:       "1001",
	}
	return codeMap[a]
}
func (a ApiResponseEnum) GetName() string {
	codeMsgMap := map[ApiResponseEnum]string{
		AIREADY_EXISTS:   "密码重复",
		SUCCESS:          "成功",
		ERROR_PARAM:      "参数错误",
		ERROR_DATA_EMPTY: "无数据",
		FAIL:             "失败",
		LOGIN_FAIL:       "用户名密码错误",
	}
	return codeMsgMap[a]
}
