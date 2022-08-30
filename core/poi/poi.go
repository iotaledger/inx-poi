package poi

import (
	"bytes"
	"context"
	"crypto"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/iotaledger/inx-app/httpserver"
	"github.com/iotaledger/iota.go/v3/merklehasher"

	// import implementation.
	_ "golang.org/x/crypto/blake2b"
)

func createProof(c echo.Context) (*ProofRequestAndResponse, error) {

	blockID, err := httpserver.ParseBlockIDParam(c, ParameterBlockID)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(CoreComponent.Daemon().ContextStopped(), 5*time.Second)
	defer cancel()

	metadata, err := deps.NodeBridge.BlockMetadata(ctx, blockID)
	if err != nil {
		return nil, err
	}

	msIndex := metadata.GetReferencedByMilestoneIndex()
	if msIndex == 0 {
		return nil, errors.WithMessagef(httpserver.ErrInvalidParameter, "block %s is not referenced by a milestone", blockID.ToHex())
	}

	ms, err := deps.NodeBridge.Milestone(ctx, msIndex)
	if err != nil {
		return nil, err
	}

	block, err := deps.NodeBridge.Block(ctx, blockID)
	if err != nil {
		return nil, err
	}

	blockIDs, err := FetchMilestoneCone(msIndex)
	if err != nil {
		return nil, err
	}

	//nolint:nosnakecase // crypto package uses underscores
	hasher := merklehasher.NewHasher(crypto.BLAKE2b_256)

	proof, err := hasher.ComputeProof(blockIDs, blockID)
	if err != nil {
		return nil, err
	}

	hash := proof.Hash(hasher)

	if !bytes.Equal(hash, ms.Milestone.InclusionMerkleRoot[:]) {
		return nil, errors.WithMessage(echo.ErrInternalServerError, "valid proof cannot be created")
	}

	return &ProofRequestAndResponse{
		Milestone: ms.Milestone,
		Block:     block,
		Proof:     proof,
	}, nil
}

func validateProof(c echo.Context) (*ValidateProofResponse, error) {

	req := &ProofRequestAndResponse{}
	if err := c.Bind(req); err != nil {
		return nil, errors.WithMessagef(httpserver.ErrInvalidParameter, "invalid request, error: %s", err)
	}

	if req.Proof == nil || req.Milestone == nil || req.Block == nil {
		return nil, errors.WithMessage(httpserver.ErrInvalidParameter, "invalid request")
	}

	// Hash the contained block to get the ID
	blockID, err := req.Block.ID()
	if err != nil {
		return nil, err
	}

	// Check if the contained proof contains the blockID
	containsValue, err := req.Proof.ContainsValue(blockID)
	if err != nil {
		return nil, err
	}
	if !containsValue {
		return &ValidateProofResponse{Valid: false}, nil
	}

	// Verify the contained Milestone signatures
	keySet := deps.KeyManager.PublicKeysSetForMilestoneIndex(req.Milestone.Index)
	if err := req.Milestone.VerifySignatures(deps.MilestonePublicKeyCount, keySet); err != nil {
		//nolint:nilerr // false positive
		return &ValidateProofResponse{Valid: false}, nil
	}

	//nolint:nosnakecase // crypto package uses underscores
	hash := req.Proof.Hash(merklehasher.NewHasher(crypto.BLAKE2b_256))

	return &ValidateProofResponse{Valid: bytes.Equal(hash, req.Milestone.InclusionMerkleRoot[:])}, nil
}
