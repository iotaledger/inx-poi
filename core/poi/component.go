package poi

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"go.uber.org/dig"

	"github.com/gohornet/inx-poi/pkg/daemon"
	"github.com/gohornet/inx-poi/pkg/nodebridge"
	"github.com/iotaledger/hive.go/app"
)

func init() {
	CoreComponent = &app.CoreComponent{
		Component: &app.Component{
			Name:     "POI",
			Params:   params,
			DepsFunc: func(cDeps dependencies) { deps = cDeps },
			Run:      run,
		},
	}
}

var (
	CoreComponent *app.CoreComponent
	deps          dependencies
)

type dependencies struct {
	dig.In
	NodeBridge *nodebridge.NodeBridge
}

func newEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	return e
}

func run() error {
	// create a background worker that handles the API
	if err := CoreComponent.Daemon().BackgroundWorker("API", func(ctx context.Context) {
		CoreComponent.LogInfo("Starting API ... done")

		e := newEcho()
		setupRoutes(e)
		go func() {
			CoreComponent.LogInfof("You can now access the API using: http://%s", ParamsPOI.BindAddress)
			if err := e.Start(ParamsPOI.BindAddress); err != nil && !errors.Is(err, http.ErrServerClosed) {
				CoreComponent.LogWarnf("Stopped REST-API server due to an error (%s)", err)
			}
		}()

		if err := deps.NodeBridge.RegisterAPIRoute(APIRoute, ParamsPOI.BindAddress); err != nil {
			CoreComponent.LogWarnf("Error registering INX api route (%s)", err)
		}

		<-ctx.Done()
		CoreComponent.LogInfo("Stopping API ...")

		if err := deps.NodeBridge.UnregisterAPIRoute(APIRoute); err != nil {
			CoreComponent.LogWarnf("Error unregistering INX api route (%s)", err)
		}

		shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := e.Shutdown(shutdownCtx); err != nil {
			CoreComponent.LogWarn(err)
		}
		shutdownCtxCancel()
		CoreComponent.LogInfo("Stopping API ... done")
	}, daemon.PriorityStopRestAPI); err != nil {
		CoreComponent.LogPanicf("failed to start worker: %s", err)
	}

	return nil
}
