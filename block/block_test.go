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
			Headers: BlockHeader{
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
}
