package response

import (
	"fmt"
	"net/http"
)

type ErrorNo struct {
	HTTPStatusCode int
	ServiceCode    int
	Message        string
}

func (this *ErrorNo) Error() string {
	return this.Message
}

func (this *ErrorNo) Status() int {
	return this.HTTPStatusCode
}

func NewBadRequest(reason string) *ErrorNo {
	return &ErrorNo{
		HTTPStatusCode: http.StatusBadRequest,
		ServiceCode:    400000,
		Message:        fmt.Sprintf("Bad Request:%v", reason),
	}
}

func NewUnauthorized(reason string) *ErrorNo {
	return &ErrorNo{
		HTTPStatusCode: http.StatusUnauthorized,
		ServiceCode:    400001,
		Message:        fmt.Sprintf("Unauthorized:%v", reason),
	}
}

func NewForbidden(reason string) *ErrorNo {
	return &ErrorNo{
		HTTPStatusCode: http.StatusForbidden,
		ServiceCode:    400003,
		Message:        fmt.Sprintf("Forbidden：%v", reason),
	}
}

func NewInternal(reason string) *ErrorNo {
	return &ErrorNo{
		HTTPStatusCode: http.StatusInternalServerError,
		ServiceCode:    500000,
		Message:        fmt.Sprintf("Internal Server Error:%v", reason),
	}
}

func NewSuccess(reason string) *ErrorNo {
	return &ErrorNo{
		HTTPStatusCode: http.StatusOK,
		ServiceCode:    0,
		Message:        reason,
	}
}
