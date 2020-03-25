package service

import (
	"sync"
	"crypto/rand"
	"crypto/sha256"
	//"time"
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
	//rounds map[int]*RoundStorage
	// done callback
	callback func(int) // callsback number of finalized blocks
	// store signatures received for current rounds
	tmpSigs map[int]*PartialSignature
	// store block proposals received for current rounds
	tmpBlocks map[int]*BlockProposal
	// node started the Consensus
	isGenesis bool

	// for firsts tests
	receivedBlockProposals map[int]*BlockProposal

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
    	receivedBlockProposals: make(map[int]*BlockProposal),
    	broadcast: b,
	}
	return n
}

func (n *Node) AttachCallback(fn func(int)) {
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
			log.Lvl1("Received unidentified message")
	}
}

func (n *Node) NewRound(round int) {
	// generate round randomness (sha256 - 32 bytes size)
	roundRandomness := n.generateRoundRandomness(uint32(round)) // should change... seed should be based on prev block sign
	log.Lvlf2("%d - Round randomness: %s",n.c.Index, hex.EncodeToString(roundRandomness))
	// pick block proposer
	proposerPosition := n.pickBlockProposer(binary.BigEndian.Uint32(roundRandomness), n.c.N)
	log.Lvlf2("Block proposer picked - position %d of %d", proposerPosition, n.c.N)
	
	if (proposerPosition == uint32(n.c.Index)) {
		log.Lvlf1("%d - I am block proposer for round %d !", n.c.Index, round)
	} else {
		return
	}

	// generate block proposal
	
	//b.Cond.L.Lock()
	//defer b.Cond.L.Unlock()
	/*
	for b.fin.HighestRound() < p.Round-1 {
		log.Lvl1("blockmaker: waiting highest round go to ", p.Round-1)
		b.Cond.Wait()
	} */

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
	log.Lvlf2("Broadcasting block proposal for round %d", round)
	go n.broadcast(n.c.Roster.List, packet)
	n.ReceivedBlockProposal(packet)
	// start round loop which will periodically check round end conditions
	//go n.roundLoop(round)
}

func (n *Node) ReceivedBlockProposal(p *BlockProposal) {
	log.Lvl2("Processing new block proposal")

	n.chain.Append(&p.Block, false)
	if (n.callback != nil) {
		n.callback(n.chain.Length())	
	}	
	n.NewRound(p.Block.BlockHeader.Round+1)

	//if p.Block.Round < n.round {
	//	log.Lvl2("received too old block")
	//	return
	//}
	// TODO: check if its in storage

	// TODO: when implementing gossip, we have to check if already received the signature

	// n.tmpBlocks[p.Round] = append(n.tmpBlocks[p.Round], p)
	// if p.Signatures[0] = nil {
	//   // block proposal does not contain any signature
	//   log.Lvl2("received block proposal without signatures")
	//   return
	// }
	// n.tmpSigs[p.Round] = append(n.tmpSigs[p.Round], p.Signatures[0])
	//n.receivedBlockProposals[p.Round]++
}

func (n *Node) ReceivedBootstrap(b *Bootstrap) {
	log.Lvl2("Processing bootstrap message")

	log.Lvl1(n.chain.Append(&b.Block, true))

	n.NewRound(0)
}

func (n *Node) roundLoop(round int) {
	log.Lvl3("Round Loop")

	// wait block time to receive messages
	//time.Sleep(n.c.BlockTime * time.Millisecond)

	// wait on new inputs
	n.Cond.Wait()

	// check round finish conditions
		//append block to finalized blockchain

}

// generates round randomness as a byte array based on a given seed
func (n *Node) generateRoundRandomness(seed uint32) []byte {
	rHash := Suite.Hash()
	binary.Write(rHash, binary.BigEndian, uint32(seed + 1)) //TODO for testing... must change
	buff := rHash.Sum(nil)
	return buff
} 

func (n *Node) pickBlockProposer(randomness uint32, listSize int) uint32 {
	return randomness % uint32(listSize)
}

func rootHash(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
