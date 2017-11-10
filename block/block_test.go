package block

import (
	"testing"
	"time"

	log "github.com/cihub/seelog"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBlocks(t *testing.T) {
	defer log.Flush()
	Convey("Construction", t, func() {
		d := []byte("hello world")
		prev := BlockHash([]byte{})
		p := NewSha256Pow(2)
		b := NewBlock(p, d, prev)
		So(b.Data, ShouldResemble, d)
		So(b.Headers, ShouldNotBeNil)
		So(b.Headers.Ts, ShouldNotEqual, time.Time{})
		So(b.Headers.PrevBlockHash, ShouldResemble, prev)
		So(b.Headers.Hash, ShouldNotBeNil)
		So(b.Headers.Hash, ShouldNotBeEmpty)
	})
}
