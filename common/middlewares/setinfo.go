package middlewares

import (
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"ptc-Game/common/response"
)

// SetUserInfo get uid from jwt token,and set to request context
func SetUserInfo() iris.Handler {
	return func(c *context.Context) {
		// get jwt token from context
		parsedToken, ok := c.Values().Get("jwt").(*jwt.Token)
		if !ok {
			response.Send(c, response.ErrJwtTokenNotFound, nil)
		}

		// assert jwt claims
		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			response.Send(c, response.ErrJwtTokenNotFound, nil)
		}

		// get uid from claims
		uid, ok := claims["uid"]
		if !ok {
			response.Send(c, response.ErrJwtTokenNotFound, nil)
		}

		// convert uid to int type
		uidF, ok := uid.(float64)
		if !ok {
			response.Send(c, response.ErrJwtTokenNotFound, nil)
		}

		c.Values().Set("userinfo.uid", int(uidF))
		c.Next()
	}
}
