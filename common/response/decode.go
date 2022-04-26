package response

import (
	"github.com/kataras/iris/v12"
	"github.com/pkg/errors"
	"net/http"
)

var (
	// common
	OK                     = &ErrorNo{HTTPStatusCode: http.StatusOK, ServiceCode: 0, Message: "OK"}
	ErrInternalServerError = &ErrorNo{HTTPStatusCode: http.StatusInternalServerError, ServiceCode: 500000, Message: "Internal Server Error"}
	ErrBadRequest          = &ErrorNo{HTTPStatusCode: http.StatusBadRequest, ServiceCode: 400000, Message: "Bad Request"}
	ErrUnauthorized        = &ErrorNo{HTTPStatusCode: http.StatusUnauthorized, ServiceCode: 400001, Message: "Unauthorized"}
	ErrForbidden           = &ErrorNo{HTTPStatusCode: http.StatusForbidden, ServiceCode: 400003, Message: "Forbidden"}
	//ErrInvalidParam        = &ErrorNo{HTTPStatusCode: http.StatusBadRequest, ServiceCode: 400003, Message: "Invalid Param"}

	//StatusBadRequest
	ErrInvalidCaptcha = &ErrorNo{HTTPStatusCode: http.StatusBadRequest, ServiceCode: 400002, Message: "Invalid Captcha"}
	//DifferentPassword = &ErrorNo{HTTPStatusCode: http.StatusBadRequest, ServiceCode: 400004, Message: "passwords are inconsistent"}

	//StatusUnauthorized
	ErrInvalidPassword    = &ErrorNo{HTTPStatusCode: http.StatusUnauthorized, ServiceCode: 400003, Message: "Invalid Account or Password"}
	VerificationCodeError = &ErrorNo{HTTPStatusCode: http.StatusUnauthorized, ServiceCode: 400005, Message: "Verification code error"}
	IrregularPassword     = &ErrorNo{HTTPStatusCode: http.StatusUnauthorized, ServiceCode: 400006, Message: "The password contains at least numbers and letters, and the length is greater than 8"}
	//ErrEmailHasExisted     = &ErrorNo{HTTPStatusCode: http.StatusInternalServerError, ServiceCode: 400010, Message: "email has existed"}
	//ErrEmailCode           = &ErrorNo{HTTPStatusCode: http.StatusInternalServerError, ServiceCode: 400011, Message: "invalid email code"}
	ErrAccountWasLocked    = &ErrorNo{HTTPStatusCode: http.StatusInternalServerError, ServiceCode: 400012, Message: "account was locked"}
	ErrEmailIsVerifying    = &ErrorNo{HTTPStatusCode: http.StatusInternalServerError, ServiceCode: 400013, Message: "email is verifying"}
	ErrEmailNotRegister    = &ErrorNo{HTTPStatusCode: http.StatusInternalServerError, ServiceCode: 400014, Message: "email has not register"}
	ErrGoogleAcountNotBind = &ErrorNo{HTTPStatusCode: http.StatusInternalServerError, ServiceCode: 400015, Message: "google account has not bind"}
	ErrJwtTokenExpired     = &ErrorNo{HTTPStatusCode: http.StatusUnauthorized, ServiceCode: 410000, Message: "jwt token expired"}

	//StatusForbidden
	ErrJwtTokenNotFound = &ErrorNo{HTTPStatusCode: http.StatusForbidden, ServiceCode: 400020, Message: "Jwt Token Not Found"}

	//StatusInternalServerError
	ErrRoleExist = &ErrorNo{HTTPStatusCode: http.StatusInternalServerError, ServiceCode: 500001, Message: "Role Has Existed"}
)

func DecodeErr(c iris.Context, err error) (int, int, string) {
	if err == nil {
		return OK.HTTPStatusCode, OK.ServiceCode, OK.Message
	}

	var e *ErrorNo
	if errors.As(err, &e) {
		//switch e.HTTPStatusCode {
		//case
		//	http.StatusNotAcceptable:
		return e.HTTPStatusCode, e.ServiceCode, e.Message
		//}
	}

	return ErrInternalServerError.HTTPStatusCode, ErrInternalServerError.ServiceCode, err.Error()
}
