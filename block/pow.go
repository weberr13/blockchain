package block

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"
)

const maxNonce = math.MaxInt64

type Sha256POW struct {
	targetBits int64
}

type ProofOfWork struct {
	block      *Block
	target     *big.Int
	targetBits int64
}

func NewSha256Pow(target int64) *Sha256POW {
	return &Sha256POW{targetBits: target}
}

func (p Sha256POW) GetPOW(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-p.targetBits))

	return &ProofOfWork{b, target, p.targetBits}
}

func IntToHex(i int64) []byte {
	return []byte(strconv.FormatInt(i, 16))
}

func (pow ProofOfWork) prepareData(nonce int64) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.Headers.PrevBlockHash,
			pow.block.Data,
			[]byte(pow.block.Headers.Ts.Format(time.RFC3339Nano)),
			IntToHex(pow.targetBits),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

func (pow ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Headers.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(pow.target) == -1
}

func (pow ProofOfWork) Run() (int64, []byte, error) {
	var hashInt big.Int
	var hash [32]byte
	nonce := int64(0)

findNonce:
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			break findNonce
		}
		nonce++
	}
	if nonce >= maxNonce || nonce < 0 {
		return 0, nil, fmt.Errorf("Max nonce exceeded")
	}
	return nonce, hash[:], nil
}
