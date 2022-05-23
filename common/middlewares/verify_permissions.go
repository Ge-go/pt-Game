package middlewares

import (
	ctx "context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"ptc-Game/account/datamodels"
	"ptc-Game/common/conf"
	"ptc-Game/common/datasource"
	"ptc-Game/common/pkg/logiclog"
	"ptc-Game/common/response"
	"time"
)

// VerifyPermissions 用户被冻结,需要登录的操作不让用户看:
func VerifyPermissions(ds *datasource.DataSources) iris.Handler {
	return func(c *context.Context) {

		// get uid from context
		uid, err := c.Values().GetInt("userinfo.uid")
		if err != nil {
			logiclog.CtxLogger(c).Warn("VerifyPermissions err: can not found uid from request context")
			response.Send(c, response.ErrForbidden, nil)
			return
		}

		// verify whether user has access permission
		var isLocked int8
		timeoutCtx, cancel := ctx.WithTimeout(ctx.Background(), 30*time.Second)
		defer cancel()

		err = ds.DB.WithContext(timeoutCtx).Model(datamodels.User{}).
			Select("is_locked").
			Where("id = ?", uid).
			Find(&isLocked).Error

		if err != nil {
			logiclog.CtxLogger(c).Errorf("VerifyPermissions Get is_locked err: %+v", err)
			response.Send(c, response.ErrInternalServerError, nil)
			return
		}

		// 用户被冻结
		if isLocked == 1 {
			response.Send(c, response.GetMessage(response.AccountHasBeenFrozen, conf.EN), nil)
			return
		}

		c.Next()
	}
}
