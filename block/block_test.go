package block

import (
	"testing"
	"time"

	log "github.com/cihub/seelog"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBlocks(t *testing.T) {
	defer log.Flush()
	Convey("Hash something", t, func() {
		b := Block{
			Headers: &BlockHeader{
				Ts:            time.Now(),
				PrevBlockHash: []byte{0x00},
			},
			Data: []byte("hello world"),
		}
		ahash := getHash(b)
		b.HashMe()
		So(b.Headers.Hash, ShouldNotBeNil)
		So(b.Headers.Hash, ShouldResemble, ahash)
		b.HashMe()
		So(b.Headers.Hash, ShouldResemble, ahash)
	})
	Convey("Construction", t, func() {
		d := []byte("hello world")
		prev := BlockHash([]byte{})
		b := NewBlock(d, prev)
		So(b.Data, ShouldResemble, d)
		So(b.Headers, ShouldNotBeNil)
		So(b.Headers.Ts, ShouldNotEqual, time.Time{})
		So(b.Headers.PrevBlockHash, ShouldResemble, prev)
		So(b.Headers.Hash, ShouldNotBeNil)
		So(b.Headers.Hash, ShouldNotBeEmpty)
	})
}
