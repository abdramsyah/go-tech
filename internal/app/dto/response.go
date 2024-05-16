package dto

import "go-tech/internal/app/constant"

type (
	HttpResponse struct {
		ProcessTime string      `json:"process_time"`
		Status      string      `json:"status"`
		Code        string      `json:"code"`
		Message     string      `json:"message"`
		Data        interface{} `json:"data,omitempty"`
	}
)

func SuccessHttpResponse(code, message string, data interface{}) HttpResponse {
	return HttpResponse{
		Status:  constant.RespSuccessStatus,
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func FailedHttpResponse(code, message string, data interface{}) HttpResponse {
	return HttpResponse{
		Status:  constant.RespFailedStatus,
		Code:    code,
		Message: message,
		Data:    data,
	}
}
