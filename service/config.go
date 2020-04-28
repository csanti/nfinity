package service

import (	
	"go.dedis.ch/kyber"
	"go.dedis.ch/kyber/share"
	"github.com/csanti/onet"
	"github.com/csanti/onet/network"
)

var ConfigType network.MessageTypeID

func init() {
	ConfigType = network.RegisterMessage(&Config{})
}

// Config holds all the parameters for the consensus protocol
type Config struct {
	Roster *onet.Roster // participants
	Index  int          // index of the node receiving this config
	N      int          // length of participants
	Public    []kyber.Point   // to reconstruct public polynomial
	Share     *share.PriShare // private share
	Threshold int             // threshold of the threshold sharing scheme
	BlockSize int             // the size of the block in bytes
	BlockTime int             // blocktime in seconds
	GossipTime int
	GossipPeers int 		// number of neighbours that each node will gossip messages to
	CommunicationMode int  	// 0 for broadcast, 1 for gossip
	MaxRoundLoops int // maximum times a node can loop on a round before alerting
}
