package block

import (
	"testing"

	log "github.com/cihub/seelog"
	. "github.com/smartystreets/goconvey/convey"
)

func TestChains(t *testing.T) {
	defer log.Flush()
	Convey("Build A Chain", t, func() {
		p := NewSha256Pow(2)

		c := NewBlockChain(p)
		c.AddBlock([]byte("hello world"))
		So(len(c.blocks), ShouldEqual, 2)
		So(c.blocks[1].Data, ShouldResemble, []byte("hello world"))
	})

	Convey("Constructor", t, func() {
		p := NewSha256Pow(2)

		c := NewBlockChain(p)
		So(c.blocks, ShouldNotBeNil)
		So(c.blocks, ShouldNotBeEmpty)
	})

	Convey("add some stuff", t, func() {
		p := NewSha256Pow(2)

		c := NewBlockChain(p)
		So(c.blocks, ShouldNotBeNil)
		So(c.blocks, ShouldNotBeEmpty)
		c.AddBlock([]byte("hello"))
		c.AddBlock([]byte("World"))
		So(len(c.blocks), ShouldEqual, 3)
		So(c.blocks[1].Headers, ShouldNotBeNil)
		So(c.blocks[2].Headers, ShouldNotBeNil)
		So(c.blocks[1].Headers.Hash, ShouldNotBeNil)
		So(c.blocks[2].Headers.Hash, ShouldNotBeNil)
		So(c.blocks[2].Headers.PrevBlockHash, ShouldResemble, c.blocks[1].Headers.Hash)
		So(c.blocks[1].Headers.PrevBlockHash, ShouldResemble, c.blocks[0].Headers.Hash)
		So(c.blocks[0].Headers.PrevBlockHash, ShouldResemble, BlockHash{})

	})
}
