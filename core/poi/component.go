package poi

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/dig"

	"github.com/iotaledger/hive.go/core/app"
	"github.com/iotaledger/inx-app/httpserver"
	"github.com/iotaledger/inx-app/nodebridge"
	"github.com/iotaledger/inx-poi/pkg/daemon"
	inx "github.com/iotaledger/inx/go"
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/iota.go/v3/keymanager"
)

func init() {
	CoreComponent = &app.CoreComponent{
		Component: &app.Component{
			Name:     "POI",
			Params:   params,
			DepsFunc: func(cDeps dependencies) { deps = cDeps },
			Provide:  provide,
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
			keyManager.AddKeyRange(keyRange.GetPublicKey(), iotago.MilestoneIndex(keyRange.GetStartIndex()), iotago.MilestoneIndex(keyRange.GetEndIndex()))
		}
		return outDeps{
			KeyManager:              keyManager,
			MilestonePublicKeyCount: int(deps.NodeBridge.NodeConfig.GetMilestonePublicKeyCount()),
		}
	})

}

func run() error {
	// create a background worker that handles the API
	if err := CoreComponent.Daemon().BackgroundWorker("API", func(ctx context.Context) {
		CoreComponent.LogInfo("Starting API ... done")

		e := httpserver.NewEcho(CoreComponent.Logger(), nil, ParamsPOI.DebugRequestLoggerEnabled)
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

func FetchMilestoneCone(index uint32) (iotago.BlockIDs, error) {
	CoreComponent.LogDebugf("Fetch cone of milestone %d\n", index)
	req := &inx.MilestoneRequest{
		MilestoneIndex: index,
	}
	stream, err := deps.NodeBridge.Client().ReadMilestoneConeMetadata(context.Background(), req)
	if err != nil {
		return nil, err
	}
	var blockIDs iotago.BlockIDs
	for {
		payload, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// We are done
				break
			}
			return nil, err
		}

		blockIDs = append(blockIDs, payload.UnwrapBlockID())
	}
	CoreComponent.LogDebugf("Milestone %d contained %d blocks\n", index, len(blockIDs))
	return blockIDs, nil
}
