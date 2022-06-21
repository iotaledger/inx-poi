package poi

import (
	iotago "github.com/iotaledger/iota.go/v3"
	"github.com/iotaledger/iota.go/v3/merklehasher"
)

type ProofRequestAndResponse struct {
	Milestone *iotago.Milestone   `json:"milestone"`
	Block     *iotago.Block       `json:"block"`
	Proof     *merklehasher.Proof `json:"proof"`
}

type ValidateProofResponse struct {
	Valid bool `json:"valid"`
}
