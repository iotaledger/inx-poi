package nodebridge

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gohornet/hornet/pkg/keymanager"
	"github.com/gohornet/hornet/pkg/model/milestone"
	"github.com/iotaledger/hive.go/logger"
	"github.com/iotaledger/hive.go/serializer/v2"
	inx "github.com/iotaledger/inx/go"
	iotago "github.com/iotaledger/iota.go/v3"
)

type NodeBridge struct {
	logger     *logger.Logger
	conn       *grpc.ClientConn
	client     inx.INXClient
	NodeConfig *inx.NodeConfiguration
}

func NewNodeBridge(ctx context.Context, address string, logger *logger.Logger) (*NodeBridge, error) {

	conn, err := grpc.Dial(address,
		grpc.WithChainUnaryInterceptor(grpc_retry.UnaryClientInterceptor(), grpc_prometheus.UnaryClientInterceptor),
		grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	client := inx.NewINXClient(conn)
	retryBackoff := func(_ uint) time.Duration {
		return 1 * time.Second
	}

	logger.Info("Connecting to node and reading protocol parameters...")
	nodeConfig, err := client.ReadNodeConfiguration(ctx, &inx.NoParams{}, grpc_retry.WithMax(5), grpc_retry.WithBackoff(retryBackoff))
	if err != nil {
		return nil, err
	}

	return &NodeBridge{
		logger:     logger,
		conn:       conn,
		client:     client,
		NodeConfig: nodeConfig,
	}, nil
}

func (n *NodeBridge) Run(ctx context.Context) {
	c, cancel := context.WithCancel(ctx)
	defer cancel()
	<-c.Done()
	n.conn.Close()
}

func (n *NodeBridge) MilestonePublicKeyCount() int {
	return int(n.NodeConfig.GetMilestonePublicKeyCount())
}

func (n *NodeBridge) KeyManager() *keymanager.KeyManager {
	keyManager := keymanager.New()
	for _, keyRange := range n.NodeConfig.GetMilestoneKeyRanges() {
		keyManager.AddKeyRange(keyRange.GetPublicKey(), milestone.Index(keyRange.GetStartIndex()), milestone.Index(keyRange.GetEndIndex()))
	}
	return keyManager
}

func (n *NodeBridge) NodeStatus() (confirmedIndex milestone.Index, pruningIndex milestone.Index) {
	status, err := n.client.ReadNodeStatus(context.Background(), &inx.NoParams{})
	if err != nil {
		return 0, 0
	}
	return milestone.Index(status.GetConfirmedMilestone().GetMilestoneIndex()), milestone.Index(status.GetPruningIndex())
}

func (n *NodeBridge) BlockMetadataForBlockID(blockID iotago.BlockID) (*inx.BlockMetadata, error) {
	return n.client.ReadBlockMetadata(context.Background(), inx.NewBlockId(blockID))
}

func (n *NodeBridge) BlockForBlockID(blockID iotago.BlockID) (*iotago.Block, error) {
	inxMsg, err := n.client.ReadBlock(context.Background(), inx.NewBlockId(blockID))
	if err != nil {
		return nil, err
	}
	return inxMsg.UnwrapBlock(serializer.DeSeriModeNoValidation, nil)
}

func (n *NodeBridge) Milestone(index uint32) (*iotago.Milestone, error) {
	req := &inx.MilestoneRequest{
		MilestoneIndex: index,
	}
	m, err := n.client.ReadMilestone(context.Background(), req)
	if err != nil {
		return nil, err
	}
	milestone := &iotago.Milestone{}
	if _, err := milestone.Deserialize(m.GetMilestone().GetData(), serializer.DeSeriModeNoValidation, nil); err != nil {
		return nil, err
	}
	return milestone, nil
}

func (n *NodeBridge) FetchMilestoneCone(index uint32) (iotago.BlockIDs, error) {
	fmt.Printf("Fetch cone of milestone %d\n", index)
	req := &inx.MilestoneRequest{
		MilestoneIndex: index,
	}
	stream, err := n.client.ReadMilestoneConeMetadata(context.Background(), req)
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
	fmt.Printf("Milestone %d contained %d blocks\n", index, len(blockIDs))
	return blockIDs, nil
}

func (n *NodeBridge) RegisterAPIRoute(route string, bindAddress string) error {
	bindAddressParts := strings.Split(bindAddress, ":")
	if len(bindAddressParts) != 2 {
		return fmt.Errorf("Invalid address %s", bindAddress)
	}
	port, err := strconv.ParseInt(bindAddressParts[1], 10, 32)
	if err != nil {
		return err
	}

	apiReq := &inx.APIRouteRequest{
		Route: route,
		Host:  bindAddressParts[0],
		Port:  uint32(port),
	}

	if err != nil {
		return err
	}
	_, err = n.client.RegisterAPIRoute(context.Background(), apiReq)
	return err
}

func (n *NodeBridge) UnregisterAPIRoute(route string) error {
	apiReq := &inx.APIRouteRequest{
		Route: route,
	}
	_, err := n.client.UnregisterAPIRoute(context.Background(), apiReq)
	return err
}
