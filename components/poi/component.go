package poi

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/dig"

	"github.com/iotaledger/hive.go/app"
	"github.com/iotaledger/inx-app/pkg/httpserver"
	"github.com/iotaledger/inx-app/pkg/nodebridge"
	"github.com/iotaledger/inx-poi/pkg/daemon"
	inx "github.com/iotaledger/inx/go"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/iota.go/v3/keymanager"
)

func init() {
	Component = &app.Component{
		Name:     "POI",
		Params:   params,
		DepsFunc: func(cDeps dependencies) { deps = cDeps },
		Provide:  provide,
		Run:      run,
	}
}

var (
	Component *app.Component
	deps      dependencies
)

type dependencies struct {
	dig.In
	NodeBridge              *nodebridge.NodeBridge
	KeyManager              *keymanager.KeyManager
	MilestonePublicKeyCount int `name:"milestonePublicKeyCount"`
}

func provide(c *dig.Container) error {

	type inDeps struct {
		dig.In
		NodeBridge *nodebridge.NodeBridge
	}

	type outDeps struct {
		dig.Out
		KeyManager              *keymanager.KeyManager
		MilestonePublicKeyCount int `name:"milestonePublicKeyCount"`
	}

	return c.Provide(func(deps inDeps) outDeps {
		keyManager := keymanager.New()
		for _, keyRange := range deps.NodeBridge.NodeConfig.GetMilestoneKeyRanges() {
			keyManager.AddKeyRange(keyRange.GetPublicKey(), keyRange.GetStartIndex(), keyRange.GetEndIndex())
		}

		return outDeps{
			KeyManager:              keyManager,
			MilestonePublicKeyCount: int(deps.NodeBridge.NodeConfig.GetMilestonePublicKeyCount()),
		}
	})

}

func run() error {
	// create a background worker that handles the API
	if err := Component.Daemon().BackgroundWorker("API", func(ctx context.Context) {
		Component.LogInfo("Starting API ... done")

		e := httpserver.NewEcho(Component.Logger(), nil, ParamsRestAPI.DebugRequestLoggerEnabled)

		Component.LogInfo("Starting API server ...")

		setupRoutes(e)
		go func() {
			Component.LogInfof("You can now access the API using: http://%s", ParamsRestAPI.BindAddress)
			if err := e.Start(ParamsRestAPI.BindAddress); err != nil && !errors.Is(err, http.ErrServerClosed) {
				Component.LogErrorfAndExit("Stopped REST-API server due to an error (%s)", err)
			}
		}()

		ctxRegister, cancelRegister := context.WithTimeout(ctx, 5*time.Second)

		advertisedAddress := ParamsRestAPI.BindAddress
		if ParamsRestAPI.AdvertiseAddress != "" {
			advertisedAddress = ParamsRestAPI.AdvertiseAddress
		}

		routeName := strings.Replace(APIRoute, "/api/", "", 1)

		if err := deps.NodeBridge.RegisterAPIRoute(ctxRegister, routeName, advertisedAddress, APIRoute); err != nil {
			Component.LogErrorfAndExit("Registering INX api route failed: %s", err)
		}
		cancelRegister()

		Component.LogInfo("Starting API server ... done")
		<-ctx.Done()
		Component.LogInfo("Stopping API ...")

		ctxUnregister, cancelUnregister := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelUnregister()

		//nolint:contextcheck // false positive
		if err := deps.NodeBridge.UnregisterAPIRoute(ctxUnregister, routeName); err != nil {
			Component.LogWarnf("Unregistering INX api route failed: %s", err)
		}

		shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCtxCancel()

		//nolint:contextcheck // false positive
		if err := e.Shutdown(shutdownCtx); err != nil {
			Component.LogWarn(err)
		}

		Component.LogInfo("Stopping API ... done")
	}, daemon.PriorityStopRestAPI); err != nil {
		Component.LogPanicf("failed to start worker: %s", err)
	}

	return nil
}

func FetchMilestoneCone(ctx context.Context, index uint32) (iotago.BlockIDs, error) {
	Component.LogDebugf("Fetch cone of milestone %d\n", index)

	fetchContext, cancel := context.WithCancel(ctx)
	defer cancel()

	var blockIDs iotago.BlockIDs
	if err := deps.NodeBridge.MilestoneConeMetadata(fetchContext, cancel, index, func(metadata *inx.BlockMetadata) {
		blockIDs = append(blockIDs, metadata.UnwrapBlockID())
	}); err != nil {
		return nil, err
	}

	Component.LogDebugf("Milestone %d contained %d blocks\n", index, len(blockIDs))

	return blockIDs, nil
}
