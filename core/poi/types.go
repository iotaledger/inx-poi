package poi

import (
	"github.com/gohornet/hornet/pkg/whiteflag"
	iotago "github.com/iotaledger/iota.go/v3"
)

type ProofRequestAndResponse struct {
	Milestone *iotago.Milestone         `json:"milestone"`
	Block     *iotago.Block             `json:"block"`
	Proof     *whiteflag.InclusionProof `json:"proof"`
}

type ValidateProofResponse struct {
	Valid bool `json:"valid"`
}
