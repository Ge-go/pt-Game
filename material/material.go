package material

import (
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/mvc"
	"ptc-Game/common/datasource"
	"ptc-Game/common/middlewares"
	"ptc-Game/common/pkg/config"
	"ptc-Game/material/repositories"
	"ptc-Game/material/services"
	"ptc-Game/material/web/controllers"
)

// Material 注册素材模块MVC
func Material(app *mvc.Application) {
	// 认证中间件: jwt 认证
	conf := config.GetConfig()

	//controller 实例
	repo := repositories.NewMaterialPermissionRepository(datasource.GetMaterialDataSource())
	materialRepo := repositories.NewMaterialRepository(datasource.GetMaterialDataSource(), datasource.GetDataSource())
	permissionService := services.NewMaterialPermissionService(repo)
	materialService := services.NewMaterialService(materialRepo)
	app.Register(permissionService)
	app.Register(materialService)

	materialCtr := &controllers.MaterialController{
		Service:    materialService,
		PerService: permissionService,
	}
	app.Router.PartyFunc("/", func(material router.Party) {
		material.Use(middlewares.JwtHandler(conf.App.JwtSecret), middlewares.SetUserInfo(), middlewares.VerifyPermissions(datasource.GetDataSource()))

		material.Get("/mateHome", materialCtr.MateHome).Name = "素材库首页"
	})
}
