package shadow

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"jingxi.cn/tools/shadow/internal/pkg/controller"
)

type App struct {
	controller *controller.Controller
}

func NewApp() *App {
	return &App{
		controller: controller.NewController(),
	}
}

func (app *App) Run(httpAddr string, dir string) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := app.controller.Run(httpAddr, dir); nil != err {
			logrus.Errorf("Failed to Run http server(%s): %+v", httpAddr, err)
			stop()
		}
	}()
	<-ctx.Done()
	app.controller.Stop()
	logrus.Infof("Service Graceful Quit!")
}
