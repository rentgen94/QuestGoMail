package management

type Response struct {
	msg    string      `json:"msg"`
	errMsg string      `json:"errMsg"`
	data   interface{} `json:"data"`
}
