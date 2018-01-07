package bean

/**
 des    统一返回数据格式
 author liupengfei
**/
type ResInfo struct {
	Code string      `json:"code"`
	Data interface{} `json:"data"`
	Desc string      `json:"desc"`
}
