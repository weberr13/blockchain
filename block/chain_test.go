package block

import (
	"testing"
	"os"
	"io/ioutil"

	log "github.com/cihub/seelog"
	. "github.com/smartystreets/goconvey/convey"
)

func TestChains(t *testing.T) {

	defer log.Flush()
	Convey("Build A Chain", t, func() {
		tmpfile, err := ioutil.TempFile("/tmp", "testchains")
		if err != nil {
			t.FailNow()
		}
		defer os.Remove(tmpfile)
		p := NewSha256Pow(2)

		c := NewBlockChain(p, tmpfile)
		So(c, ShouldNotBeNil)
		defer c.Close()
		c.AddBlock([]byte("hello world"))
		So(len(c.blocks), ShouldEqual, 2)
		So(c.blocks[1].Data, ShouldResemble, []byte("hello world"))
	})

	Convey("Constructor", t, func() {
		tmpfile, err := ioutil.TempFile("/tmp", "testchains")
		if err != nil {
			t.FailNow()
		}
		defer os.Remove(tmpfile)
		p := NewSha256Pow(2)

		c := NewBlockChain(p, tmpfile)
		So(c, ShouldNotBeNil)
		defer c.Close()
		So(c.blocks, ShouldNotBeNil)
		So(c.blocks, ShouldNotBeEmpty)
	})

	Convey("add some stuff", t, func() {
		tmpfile, err := ioutil.TempFile("/tmp", "testchains")
		if err != nil {
			t.FailNow()
		}
		defer os.Remove(tmpfile)
		p := NewSha256Pow(2)

		c := NewBlockChain(p)
		So(c, ShouldNotBeNil)
		defer c.Close()
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
		for _, block := range c.blocks {
			proof := p.GetPOW(block)
			So(proof.Validate(), ShouldBeTrue)
		}
	})
}
