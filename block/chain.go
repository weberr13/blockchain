package block

import (
	"fmt"

	"github.com/boltdb/bolt"
	log "github.com/cihub/seelog"
)

const blocksBucket = "myblocks"

//Chain of blocks (simple)
type Chain struct {
	Tip []byte
	pow POWBuilder
	db  *bolt.DB
}

type ChainIterator struct {
	current Hash
	db      *bolt.DB
}

func (c *Chain) Iterator() *ChainIterator {
	return &ChainIterator{c.Tip, c.db}
}

func (i *ChainIterator) Next() (block *Block, err error) {
	err = i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		blockBytes := b.Get(i.current)
		block, err = DeserializeBlock(blockBytes)

		return err
	})

	if err != nil {
		return nil, err
	}

	i.current = block.Headers.PrevBlockHash

	return block, nil
}

//Close will shut down the db connection/etc
func (c *Chain) Close() error {
	return c.db.Close()
}

//AddBlock to chain
func (c *Chain) AddBlock(data []byte) (err error) {
	var lastHash []byte

	err = c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	newBlock := NewBlock(c.pow, data, lastHash)

	err = c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		j, err := newBlock.Serialize()
		if err != nil {
			return err
		}
		err = b.Put(newBlock.Headers.Hash, j)
		if err != nil {
			return err
		}
		err = b.Put([]byte("l"), newBlock.Headers.Hash)
		if err != nil {
			return err
		}
		c.Tip = newBlock.Headers.Hash
		return nil
	})
	return err
}

//NewBlockChain with valid genesis
func NewBlockChain(pow POWBuilder, dbfile string) *Chain {
	var tip []byte
	db, err := bolt.Open(dbfile, 0600, nil)
	if err != nil {
		log.Error("cannot create chain: ", err)
		return nil
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			gen := NewGenesisBlock(pow)
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				return err
			}
			j, err := gen.Serialize()
			if err != nil {
				return err
			}
			err = b.Put(gen.Headers.Hash, j)
			if err != nil {
				return err
			}
			err = b.Put([]byte("l"), gen.Headers.Hash)
			tip = gen.Headers.Hash
			return nil
		}

		tip = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		db.Close()
		log.Error("cannot update chain: ", err)
		return nil
	}

	return &Chain{Tip: tip, db: db, pow: pow}
}

//Walk the chain performing action and fininshing when done
func (c Chain) Walk(action func(b *Block) error, done func(b *Block) bool) error {
	i := c.Iterator()

	for b, err := i.Next(); b != nil && err == nil && !done(b); b, err = i.Next() {
		proof := c.pow.GetPOW(b)
		if !proof.Validate() {
			return fmt.Errorf("chain corrupt")
		}
		if b.IsGenesis() {
			return nil
		}
		err = action(b)
		if err != nil {
			return err
		}
	}
	return nil
}
