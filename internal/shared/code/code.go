package code

import "fmt"

type AppCode struct {
	Code     string
	HttpCode int
}

func New(status int) AppCode {
	return AppCode{
		Code:     fmt.Sprintf("HTTP_ERROR_%d", status),
		HttpCode: status,
	}
}

func Parse(status int) AppCode {
	if code, ok := HttpCodeAppCode[status]; ok {
		return code
	}

	return New(status)
}
