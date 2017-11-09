package block

import (
	"bytes"
	"crypto/sha256"
	"time"
)

//BlockHash ...
type BlockHash []byte

//BlockHeader ...
type BlockHeader struct {
	Ts            time.Time
	PrevBlockHash BlockHash
	Hash          BlockHash
}

//Block ...
type Block struct {
	Data    []byte
	Headers BlockHeader
}

//HashMe hashes a block once it has data
func (b *Block) HashMe() {
	b.Headers.Hash = getHash(*b)
}

func getHash(b Block) BlockHash {
	t := []byte(b.Headers.Ts.Format(time.RFC3339Nano))
	headers := bytes.Join([][]byte{b.Headers.PrevBlockHash, b.Data, t}, []byte{})
	hash := sha256.Sum256(headers)
	h := hash[:]

	return h
}
