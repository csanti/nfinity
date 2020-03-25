package service

import (
	"sync"
	"crypto/rand"
	"crypto/sha256"
	"time"
	"encoding/binary"
	"encoding/hex"

	"go.dedis.ch/onet"
	"go.dedis.ch/onet/log"
	"go.dedis.ch/onet/network"
)

type Node struct {
	*sync.Cond
	*onet.ServiceProcessor

	// config
	c *Config
	// current round number
	round int
	// finalized Chain
	chain *BlockChain
	// information of previous rounds
	rounds map[int]*RoundStorage
	// done callback
	callback func(int) // callsback number of finalized blocks
	// node started the Consensus
	isGenesis bool

	broadcast BroadcastFn
}

func NewNodeProcess(c *onet.Context, conf *Config, b BroadcastFn) *Node {
	// need to create chain first
	chain := new(BlockChain)
	n := &Node {
		ServiceProcessor: onet.NewServiceProcessor(c),
		Cond: sync.NewCond(new(sync.Mutex)),
		chain: chain,
		c:     conf,
    	broadcast: b,
    	rounds: make(map[int]*RoundStorage),
	}
	return n
}

func (n *Node) AttachCallback(fn func(int)) {
	// usually only attached to one of the nodes to notify a higher layer of the progress
	n.callback = fn
}

func (n *Node) StartConsensus() {
	log.Lvl1("Staring consensus")
	n.isGenesis = true
	packet := &Bootstrap{
		Block: GenesisBlock,
		Seed:  1234,
	}
	log.Lvl2("Starting consensus, sending bootstrap..")
	// send bootstrap message to all nodes
	go n.broadcast(n.c.Roster.List, packet)
	n.ReceivedBootstrap(packet)
}

func (n *Node) Process(e *network.Envelope) {
	n.Cond.L.Lock()
	defer n.Cond.L.Unlock()
	defer n.Cond.Broadcast()
	switch inner := e.Msg.(type) {
		case *BlockProposal:
			n.ReceivedBlockProposal(inner)
		case *Bootstrap:			
			n.ReceivedBootstrap(inner)			
		default:
			log.Lvl2("Received unidentified message")
	}
}

func (n *Node) NewRound(round int) {
	// new round can only be called after previous round is finished, so this is safe
	n.round = round
	// generate round randomness (sha256 - 32 bytes size)
	roundRandomness := n.generateRoundRandomness(round) // should change... seed should be based on prev block sign
	log.Lvlf2("%d - Round randomness: %s",n.c.Index, hex.EncodeToString(roundRandomness))
	// create the round storage
	n.rounds[round] = NewRoundStorage(round, binary.BigEndian.Uint32(roundRandomness))
	// pick block proposer
	proposerPosition := n.pickBlockProposer(binary.BigEndian.Uint32(roundRandomness), n.c.N)
	log.Lvlf2("Block proposer picked - position %d of %d", proposerPosition, n.c.N)
	n.rounds[round].ProposerIndex = proposerPosition
	
	// start round loop which will periodically check round end conditions
	go n.roundLoop(round)

	// check if node is proposer, if not: returns
	if (proposerPosition == n.c.Index) {
		log.Lvlf2("%d - I am block proposer for round %d !", n.c.Index, round)
	} else {
		return
	}

	// generate block proposal
	oldBlock := n.chain.Head()
	blob := make([]byte, n.c.BlockSize)
	rand.Read(blob)
	hash := rootHash(blob)
	header := BlockHeader {
		Round:      round,
		Owner:      n.c.Index,
		Root:       hash,
		Randomness: binary.BigEndian.Uint32(roundRandomness),
		PrvHash:    oldBlock.BlockHeader.Hash(),
		PrvSig:     oldBlock.BlockHeader.Signature,
	}
	blockProposal := Block {
		BlockHeader: header,
		Blob:        blob,
	}
	packet := &BlockProposal {
		Block: blockProposal,
		// need to add the proposer signature to the array
	}
	log.Lvlf3("Broadcasting block proposal for round %d", round)
	go n.broadcast(n.c.Roster.List, packet)
	n.ReceivedBlockProposal(packet)
}

func (n *Node) ReceivedBlockProposal(p *BlockProposal) {
	if p.Block.Round < n.round {
		log.Lvl2("received too old block")
		return
	}
	blockRound := p.Block.BlockHeader.Round
	if (p.Block.BlockHeader.Owner != n.rounds[blockRound].ProposerIndex) {
		log.Lvl2("received block with invalid proposer")
		return
	}
	// TODO check owner signature
	// is from valid proposer, we have to check if already received a block from him.
	if (!n.rounds[blockRound].ReceivedValidBlock) {
		n.rounds[blockRound].StoreValidBlock(&p.Block)		
	} else if (p.Block.BlockHeader.Hash() != n.rounds[blockRound].BlockHash) {
		log.Lvl1("Received two different blocks from proposer!")
		// TODO handle malicious case
		return
	}
	n.rounds[blockRound].StoreBlockProposal(p)
}

func (n *Node) ReceivedBootstrap(b *Bootstrap) {
	log.Lvl3("Processing bootstrap message... starting consensus")
	// add genesis and start new round
	if(n.chain.Append(&b.Block, true) != 1) {
		panic("this should never happen")
	}
	n.NewRound(0)
}

func (n *Node) roundLoop(round int) {
	log.Lvlf3("Starting round %d loop",round)
	// 
	defer func() {
		log.Lvlf3("Exiting round %d loop",round)
		n.NewRound(round+1)
		delete(n.rounds, round)
	}()
	var times int = 0
	for {
		// wait block time to receive messages
		time.Sleep(time.Duration(n.c.GossipTime) * time.Millisecond)

		// wait on new inputs
		//n.Cond.Wait()

		if times == 3 {
			return
		}
		n.rounds[round].ProcessBlockProposals()

		times++
		// iterate received proposals
			// use trackid to only check last iteration from each source
		//



	}
	// check round finish conditions
		//append block to finalized blockchain

	return
	/*
	n.chain.Append(&p.Block, false)
	if (n.callback != nil) {
		n.callback(n.chain.Length())	
	}	
	n.NewRound(p.Block.BlockHeader.Round+1)
	*/

}

// generates round randomness as a byte array based on a given seed
func (n *Node) generateRoundRandomness(seed int) []byte {
	rHash := Suite.Hash()
	binary.Write(rHash, binary.BigEndian, seed + 1) //TODO for testing... must change
	buff := rHash.Sum(nil)
	return buff
} 

func (n *Node) pickBlockProposer(randomness uint32, listSize int) int {
	return int(randomness) % listSize
}

func rootHash(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
