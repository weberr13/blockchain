package block

//Chain of blocks (simple)
type Chain struct {
	blocks []*Block
	pow    POWBuilder
}

//AddBlock to chain
func (c *Chain) AddBlock(data []byte) {
	if c.blocks == nil {
		c.blocks = []*Block{}
		c.blocks = append(c.blocks, NewGenesisBlock(c.pow))
	}
	prev := c.blocks[len(c.blocks)-1]
	new := NewBlock(c.pow, data, prev.GetHash())
	c.blocks = append(c.blocks, new)
}

//NewBlockChain with valid genesis
func NewBlockChain(pow POWBuilder) *Chain {
	return &Chain{blocks: []*Block{NewGenesisBlock(pow)}, pow: pow}
}
