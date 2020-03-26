package service

import (
	"go.dedis.ch/onet/log"
)

type RoundStorage struct {
	Round int
	Randomness uint32
	ProposerIndex int
	ReceivedValidBlock bool // true if we received a block from valid proposer
	Block     *Block
	BlockHash string
	FinalSig  []byte         // when notarization happenned
	Sigs      map[int]*PartialSignature // all signatures for the blob received so far
	SigCount int
	Finalized bool           // true if already notarized
	ReceivedBlockProposals int
	SentBlockProposals int
	StoredBlockProposals int
	TmpBlockProposals map[int]*BlockProposal
	//pub       *share.PubPoly
}


func NewRoundStorage(round int, randomness uint32) *RoundStorage {
	return &RoundStorage {
		Round: round,
		Randomness: randomness,
		Sigs: make(map[int]*PartialSignature),
		TmpBlockProposals: make(map[int]*BlockProposal),
	}
}

func (rs *RoundStorage) StoreValidBlock(b *Block) {
	rs.Block = b
	rs.BlockHash = b.BlockHeader.Hash()
	rs.ReceivedValidBlock = true
}

// stores block proposal after checking that it is the latest iteration of a particular sender
// also deletetes previous block proposal iterations from the sender
func (rs *RoundStorage) StoreBlockProposal(p *BlockProposal) {
	rs.ReceivedBlockProposals++
	source := int(p.TrackId / 10)

	// check if a newer iteration from the same sender was previously received
	for i := p.TrackId; i < source*10 + 10; i++ {
		_, exists := rs.TmpBlockProposals[p.TrackId]
		if exists {
			// later iteration exists
			return
		}
	}
	// delete all previsou iterations, this will make round loop faster
	for i := source*10 ; i < p.TrackId ; i++ {
		_, exists := rs.TmpBlockProposals[i]
		if exists {
			delete(rs.TmpBlockProposals, i)
			rs.StoredBlockProposals--
		}
	}
	rs.StoredBlockProposals++
	rs.TmpBlockProposals[p.TrackId] = p
}

func (rs *RoundStorage) ProcessBlockProposals() ([]*PartialSignature, bool) {
	initialSigCount := rs.SigCount
	for _ , bp := range rs.TmpBlockProposals {
		for _, ps := range bp.Signatures {
			rs.Sigs[ps.Signer] = ps
		}
	}
	rs.SigCount = len(rs.Sigs)
	// TODO i could save one map conversion if i save the array in memory and use it again when there is no info change
	sigsArray := rs.mapToArray(rs.Sigs)
	log.Lvlf2("Finished processing block proposals - sign count = %d (%d new)",rs.SigCount, rs.SigCount - initialSigCount)
	if (rs.SigCount - initialSigCount) > 0 {
		return sigsArray, true
	} else {
		return sigsArray, false
	}
}

func (rs *RoundStorage) mapToArray(m map[int]*PartialSignature) []*PartialSignature {
	array := make([]*PartialSignature, len(m))
	for _, p := range m {
		array = append(array,p)
	}
	return array
}

