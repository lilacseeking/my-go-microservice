// api/response/response.go
package response

import "net/http"

// Response 定义标准响应结构
type Response struct {
	Code    int         `json:"code"`           // 业务状态码
	Message string      `json:"message"`        // 提示信息
	Data    interface{} `json:"data,omitempty"` // 返回数据，omitempty 表示空值时省略
}

// Success 返回成功响应
func Success(data interface{}) *Response {
	return &Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	}
}

// Error 返回错误响应
func Error(code int, message string) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}
