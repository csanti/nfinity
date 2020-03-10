package service

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"sync"
)

type BlockChain struct {
	sync.Mutex
	all    []*Block
	last   *Block
	length int
}

// BlockHeader represents all the information regarding a block
type BlockHeader struct {
	Round      int    // round of the block
	Owner      int    // index of the owner of the block
	Root       string // hash of the data
	Randomness int64  // randomness of the round
	PrvHash    string // hash of the previous block
	PrvSig     []byte // signature of the previous block (i.e. notarization)
}

// Block represents how a block is stored locally
// Block is first sent from a block maker
type Block struct {
	BlockHeader
	Blob []byte // the actual content
}

// Hash returns the hash in hexadecimal of the header
func (h *BlockHeader) Hash() string {
	hash := Suite.Hash()
	binary.Write(hash, binary.BigEndian, h.Owner)
	binary.Write(hash, binary.BigEndian, h.Round)
	hash.Write([]byte(h.PrvHash))
	hash.Write([]byte(h.Root))
	hash.Write(h.PrvSig)
	buff := hash.Sum(nil)
	return hex.EncodeToString(buff)
}

func rootHash(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// GenesisBlock is the first block of the chain
var GenesisBlock = &Block{
	BlockHeader: BlockHeader{
		Round: 0,
		Owner: -1,
		Root:  "6afbc27f4ae8951a541be53038ca20d3a9f18f60a38b1dc2cd48a46ff5d26ace",
		// sha256("hello world")
		PrvHash: "a948904f2f0f479b8f8197694b30184b0d2ed1c1cd2a1ec0fb85d299a192a447",
		// echo "hello world" | sha256sum | sha256sum
		PrvSig: []byte("3605ff73b6faec27aa78e311603e9fe2ef35bad82ccf46fc707814bfbdcc6f9e"),
	},
	Blob: []byte("Hello Genesis"),
}
