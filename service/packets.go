package service

import (
	"go.dedis.ch/onet/network"
)

var BlockProposalType network.MessageTypeID
var BootstrapType network.MessageTypeID

func init() {
	BlockProposalType = network.RegisterMessage(&BlockProposal{})
	BootstrapType = network.RegisterMessage(&Bootstrap{})
}

type BlockProposal struct {
	TrackId int
	*Block
	Count int // count of parital signatures in the array
	Signatures []*PartialSignature // Partial signature from the signer
}

type PartialSignature struct {
	Signer int
	Partial []byte
}

type Bootstrap struct {
	*Block
	Seed int
}
