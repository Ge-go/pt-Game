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
	app.Router.Post("/minorCert", accountContrl.MinorCert).Name = "未成年人注册"

	//初始化docusign控制器
	//docusignRepo := repositories.NewDocusignRepository(datasource.GetDataSource())
	//docusignService := services.NewDocusignService(docusignRepo)
}
