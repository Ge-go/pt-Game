package main

import (
	"flag"
	"github.com/kataras/iris/v12/mvc"
	"ptc-Game/account"
	"ptc-Game/common/app"
	"ptc-Game/common/datasource"
	"ptc-Game/common/pkg/config"
	"ptc-Game/common/pkg/email"
	"ptc-Game/common/pkg/logiclog"
	_ "ptc-Game/docs" // 初始化swagger
)

var confFile = flag.String("c", "", "配置文件名称")

// @title Swagger API文档
// @version 1.0
// @description 海外内容创作者平台-用户端 API文档

// @contact.name chenguangWang
// @contact.url http://www.swagger.io/support
// @contact.email chenguangWang

// @host localhost:8000

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /
func main() {
	flag.Parse()
	//初始化配置文件
	conf, err := config.Init(*confFile)
	if err != nil {
		panic(err)
	}

	//初始化日志(logic log)
	logiclog.InitConfig(conf.App.Name, conf.App.RunMode, conf.App.InstantKey, conf.Log.LoggerLevel)

	// 初始化 Datasource
	_, err = datasource.NewDataSources(conf)
	if err != nil {
		panic(err)
	}

	// 初始化 MaterialSource
	_, err = datasource.NewMaterialDataSources(conf)
	if err != nil {
		panic(err)
	}

	// 初始化邮箱
	_, err = email.NewEmailSender(conf.EmailAPI.Domain, conf.EmailAPI.Uri, conf.EmailAPI.GameId, conf.EmailAPI.Sigkey, conf.EmailAPI.ChannelId, conf.EmailAPI.Source)
	if err != nil {
		panic(err)
	}

	// 初始化app(封装access log,swagger,graceful shutdown等功能)
	app, err := app.New(conf)
	if err != nil {
		panic(err)
	}

	if conf.App.UseSwagger {
		app.UseSwagger()
	}

	// 注册路由
	// 注册路由
	apiv1 := app.Iris.Party("/api/v1")
	//apiv2 := app.Iris.Party("/api/v2")
	mvc.Configure(apiv1.Party("/user"), account.Account) // 账号模块
	//mvc.Configure(apiv1.Party("/index"), index.Index)          // 首页模块
	//mvc.Configure(apiv1.Party("/userInfo"), userinfo.UserInfo) // 个人信息模块
	//mvc.Configure(apiv1.Party("/task"), task.Task)             // 任务模块
	//mvc.Configure(apiv1.Party("/message"), message.Message)    //消息模块
	//mvc.Configure(apiv2.Party("/material"), material.Material) //素材模块
	//mvc.Configure(apiv2.Party("/user"), user_v2.User)          //素材模块
	//mvc.Configure(apiv2.Party("/task"), task_v2.Task)          //素材模块s

	app.Run()
}
