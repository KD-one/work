package serializer

// Response 通用Response结构体
type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Msg   string      `json:"message"`
	Error string      `json:"error,omitempty"`
}
