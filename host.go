package gpfx

import (
	"fmt"
	"github.com/lonverce/gpfx/service"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Host struct {
	rootScope service.LifetimeScope
	modules   []*Module
	mgr       *HostedServiceManager
}

func (app *Host) Start() {
	log.Println("[gpfx] 正在初始化应用程序")
	initCtx := &ModuleInitializer{
		ctx: app.rootScope.GetProvider(),
	}

	for _, module := range app.modules {
		if module.OnApplicationInitialize == nil {
			continue
		}
		module.OnApplicationInitialize(initCtx)
	}

	log.Println("[gpfx] 应用程序初始化完成, 正在启动服务集")
	var mgr *HostedServiceManager
	service.MustLoad(app.rootScope.GetProvider(), &mgr)
	mgr.StartAllServices()
	app.mgr = mgr
	log.Println("[gpfx] 应用程序服务集启动成功.")
}

func (app *Host) Stop() {
	log.Println("[gpfx] 即将停止应用程序服务集")
	app.mgr.StopAllServices()
	log.Println("[gpfx] 应用程序服务集已全部停止")
}

func (app *Host) Run() {
	app.Start()
	fmt.Println("[gpfx] 按下'Ctrl+C'退出")
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.Stop()
}
