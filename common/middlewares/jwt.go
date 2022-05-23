package middlewares

import (
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/pkg/errors"
	"ptc-Game/common/response"
)

// JwtHandler jwt认证中间件
func JwtHandler(secret string) iris.Handler {
	mySecret := []byte(secret)

	return jwt.New(jwt.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return mySecret, nil
		},
		Expiration:    true,
		SigningMethod: jwt.SigningMethodHS256,
		ErrorHandler: func(c iris.Context, err error) {

			// response when jwt token expired
			var jwtErr *jwtgo.ValidationError
			if errors.As(err, &jwtErr) && jwtErr.Errors == jwtgo.ValidationErrorExpired {
				response.Send(c, response.ErrJwtTokenExpired, nil)
				return
			}
			response.Send(c, response.ErrJwtTokenNotFound, nil)
		},
	}).Serve
}
