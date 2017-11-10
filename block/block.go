package block

import (
	"encoding/json"
	"time"

	log "github.com/cihub/seelog"
)

type POWBuilder interface {
	GetPOW(b *Block) *ProofOfWork
}

//BlockHash ...
type BlockHash []byte

//BlockHeader ...
type BlockHeader struct {
	Ts            time.Time
	PrevBlockHash BlockHash
	Hash          BlockHash
	Nonce         int64
}

//Block ...
type Block struct {
	Data    []byte
	Headers *BlockHeader
}

//NewBlock ...
func NewBlock(pow POWBuilder, d []byte, prevHash BlockHash) *Block {
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

func NewHeader(prevHash BlockHash) *BlockHeader {
	h := &BlockHeader{
		Ts:            time.Now(),
		PrevBlockHash: prevHash,
	}
	return h
}

func (b Block) GetHash() BlockHash {
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
