package block

import (
	"encoding/json"
	"time"

	log "github.com/cihub/seelog"
)

type POWBuilder interface {
	GetPOW(b *Block) *ProofOfWork
}

//Hash ...
type Hash []byte

//BlockHeader ...
type BlockHeader struct {
	Ts            time.Time
	PrevBlockHash Hash
	Hash          Hash
	Nonce         int64
}

//Block ...
type Block struct {
	Data    []byte
	Headers *BlockHeader
}

//IsGenesis block
func (b Block) IsGenesis() bool {
	return b.Headers.PrevBlockHash == nil 
}
//NewBlock ...
func NewBlock(pow POWBuilder, d []byte, prevHash Hash) *Block {
	b := &Block{
		Headers: NewHeader(prevHash),
		Data:    d,
	}
	p := pow.GetPOW(b)
	nonce, hash, err := p.Run()
	if err != nil {
		log.Error("failed to build block: ", err)
	}
	b.Headers.Hash = hash[:]
	b.Headers.Nonce = nonce
	return b
}

//NewGenisisBlock to start a chain
func NewGenesisBlock(pow POWBuilder) *Block {
	return NewBlock(pow, []byte("In the beginning..."), []byte{})
}

func NewHeader(prevHash Hash) *BlockHeader {
	h := &BlockHeader{
		Ts:            time.Now(),
		PrevBlockHash: prevHash,
	}
	return h
}

func (b Block) GetHash() Hash {
	return b.Headers.Hash
}

func (b Block) Serialize() ([]byte, error) {
	return json.Marshal(b)
}

func DeserializeBlock(d []byte) (b *Block, err error) {
	b = &Block{}
	err = json.Unmarshal(d, b)
	return b, err
}
