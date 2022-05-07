package account

import (
	"github.com/kataras/iris/v12/mvc"
	"ptc-Game/account/repositories"
	"ptc-Game/account/services"
	"ptc-Game/account/web/controllers"
	"ptc-Game/common/datasource"
)

func Account(app *mvc.Application) {
	// 初始化
	//conf := config.GetConfig()
	repo := repositories.NewAccountRepository(datasource.GetDataSource())
	accountService := services.NewAccountService(repo)
	youtubeService := services.NewYoutubeService(repo)
	app.Register(accountService)
	accountContrl := &controllers.AccountController{accountService, youtubeService}

	// 账号管理
	app.Router.Post("/register", accountContrl.Register).Name = "注册"
	app.Router.Post("/verifyEmail", accountContrl.VerifyEmail).Name = "邮箱校验"
	app.Router.Post("/login", accountContrl.Login).Name = "账号登录"
	app.Router.Post("/resetPassword", accountContrl.ResetPassword).Name = "密码找回"
	app.Router.Get("/captcha", accountContrl.GetCaptcha).Name = "获取图片验证码"
	app.Router.Get("/tags", accountContrl.GetUserTag).Name = "获取用户标签"
	app.Router.Get("/youtubeLogin", accountContrl.YoutubeLogin).Name = "youtube账号登录"
	app.Router.Post("/minorCert", accountContrl.MinorCert).Name = "未成年人注册"
	app.Router.Post("/checkEmailAndPassword", accountContrl.CheckEmailAndPassword).Name = "验证邮箱和密码"

	//初始化docusign控制器
	//docusignRepo := repositories.NewDocusignRepository(datasource.GetDataSource())
	//docusignService := services.NewDocusignService(docusignRepo)
}
