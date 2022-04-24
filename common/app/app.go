package app

import (
	"context"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/host"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/middleware/requestid"
	"ptc-Game/common/pkg/accesslog"
	"ptc-Game/common/pkg/config"
	"strings"
	"time"
)

// App 封装了server初始化时的一些操作
type App struct {
	Config *config.Config
	Iris   *iris.Application
}

type AppConfig struct {
	Config *config.Config
}

// New 根据配置返回 初始化后的App
func New(config *config.Config) (*App, error) {
	app := &App{
		Config: config,
	}
	if config.App.RunMode == "dev" {
		app.Iris = iris.Default()
	} else {
		app.Iris = app.newAPiServer(config)
	}

	if err := app.init(); err != nil {
		return nil, err
	}
	return app, nil
}

func (this *App) init() error {
	govalidator.SetFieldsRequiredByDefault(true)
	govalidator.SetNilPtrAllowedByRequired(true)

	return nil
}

//Run 启动server,默认支持graceful shutdown
func (this *App) Run() error {
	// 实现graceful shutdown
	idleConnsClosed := make(chan struct{})
	iris.RegisterOnInterrupt(func() {
		timeout := 10 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		// close all hosts
		this.Iris.Shutdown(ctx)
		close(idleConnsClosed)
	})

	// 启动 API server
	this.Iris.Listen(
		this.Config.App.Addr,
		iris.WithoutInterruptHandler,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
	<-idleConnsClosed
	return nil
}

// newAPiServer 注册一些常用中间件 (recover,requestId,accessLog)
func (this *App) newAPiServer(config *config.Config) *iris.Application {
	app := iris.New()

	logLevel := strings.ToLower(config.Log.LoggerLevel)
	app.Logger().Info(fmt.Sprintf(`Log level set to "%s"`, logLevel))

	ac := accesslog.GetAccessLog(config.App.Name, config.App.RunMode, config.App.InstantKey)
	app.ConfigureHost(func(su *host.Supervisor) {
		su.RegisterOnShutdown(
			func() {
				ac.Close()
			})
	})
	ac.AddOutput(app.Logger().Printer)
	app.UseRouter(ac.Handler)

	// request id 中间件
	app.UseRouter(requestid.New())

	// recover 中间件
	app.UseRouter(recover.New())

	return app
}

// UseSwagger 启动swagger API 文档
func (this *App) UseSwagger() {
	swaggerUI := swagger.WrapHandler(swaggerFiles.Handler, func(c *swagger.Config) {
		c = &swagger.Config{
			URL:         fmt.Sprintf("http://%s/swagger/doc.json", this.Config.App.Addr),
			DeepLinking: true,
		}
	})
	this.Iris.Get("/swagger", swaggerUI)
	this.Iris.Get("/swagger/{any:path}", swaggerUI)
}
