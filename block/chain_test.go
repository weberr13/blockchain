package block

import (
	"io/ioutil"
	"os"
	"testing"

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
		defer os.Remove(tmpfile.Name())
		p := NewSha256Pow(2)

		c := NewBlockChain(p, tmpfile.Name())
		So(c, ShouldNotBeNil)
		defer c.Close()
		c.AddBlock([]byte("hello world"))
		i := c.Iterator()
		blocks := []*Block{}
		for b, err := i.Next(); b != nil && err == nil; b, err = i.Next() {
			blocks = append(blocks, b)
		}
		So(len(blocks), ShouldEqual, 2)
		So(blocks[0].Data, ShouldResemble, []byte("hello world"))
	})

	Convey("Constructor", t, func() {
		tmpfile, err := ioutil.TempFile("/tmp", "testchains")
		if err != nil {
			t.FailNow()
		}
		defer os.Remove(tmpfile.Name())
		p := NewSha256Pow(2)

		c := NewBlockChain(p, tmpfile.Name())
		So(c, ShouldNotBeNil)
		defer c.Close()
		blocks := []*Block{}
		i := c.Iterator()
		for b, err := i.Next(); b != nil && err == nil; b, err = i.Next() {
			blocks = append(blocks, b)
		}
		So(blocks, ShouldNotBeEmpty)
	})

	Convey("add some stuff", t, func() {
		tmpfile, err := ioutil.TempFile("/tmp", "testchains")
		if err != nil {
			t.FailNow()
		}
		defer os.Remove(tmpfile.Name())
		p := NewSha256Pow(2)

		c := NewBlockChain(p, tmpfile.Name())
		So(c, ShouldNotBeNil)
		defer c.Close()
		c.AddBlock([]byte("hello"))
		c.AddBlock([]byte("World"))
		blocks := []*Block{}
		i := c.Iterator()
		for b, err := i.Next(); b != nil && err == nil; b, err = i.Next() {
			blocks = append(blocks, b)
		}
		So(blocks, ShouldNotBeEmpty)
		So(len(blocks), ShouldEqual, 3)
		So(blocks[1].Headers, ShouldNotBeNil)
		So(blocks[0].Headers, ShouldNotBeNil)
		So(blocks[1].Headers.Hash, ShouldNotBeNil)
		So(blocks[0].Headers.Hash, ShouldNotBeNil)
		So(blocks[0].Headers.PrevBlockHash, ShouldResemble, blocks[1].Headers.Hash)
		So(blocks[1].Headers.PrevBlockHash, ShouldResemble, blocks[2].Headers.Hash)
		So(blocks[2].Headers.PrevBlockHash, ShouldResemble, Hash{})
		for _, block := range blocks {
			proof := p.GetPOW(block)
			So(proof.Validate(), ShouldBeTrue)
		}
	})
}
