package service

import (
	"go.dedis.ch/onet/network"
)

var BlockProposalType network.MessageTypeID

func init() {
	BlockProposalType = network.RegisterMessage(&BlockProposal)
}

type BlockProposal struct {
	*Block
	Count                         // count of parital signatures in the array
	Signatures []PartialSignature // Partial signature from the signer
}

type PartialSignature struct {
	Partial []byte
}

type Bootstrap struct {
	*Block
	Seed int
}
