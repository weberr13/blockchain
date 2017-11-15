package merkle

import (
	"crypto/sha256"
	"testing"

	log "github.com/cihub/seelog"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNodes(t *testing.T) {
	defer log.Flush()
	z := sha256.Sum256([]byte{0x00})
	o := sha256.Sum256([]byte{0x01})
	var ZeroHash, OneHash, ZeroOne []byte

	ZeroHash = z[:]
	OneHash = o[:]
	zo := sha256.Sum256(append(ZeroHash, OneHash...))
	ZeroOne = zo[:]
	Convey("Construction", t, func() {
		left := NewNode(nil, nil, []byte{0x00})
		So(left, ShouldNotBeNil)
		So(left.Left, ShouldBeNil)
		So(left.Right, ShouldBeNil)
		So(left.Data, ShouldResemble, ZeroHash)

		right := NewNode(nil, nil, []byte{0x01})
		So(right, ShouldNotBeNil)
		So(right.Left, ShouldBeNil)
		So(right.Right, ShouldBeNil)
		So(right.Data, ShouldResemble, OneHash)

		root := NewNode(left, right, nil)
		So(root, ShouldNotBeNil)
		So(root.Left, ShouldEqual, left)
		So(root.Right, ShouldEqual, right)
		So(root.Data, ShouldResemble, ZeroOne)
	})
	Convey("Leaf detection", t, func() {
		So(IsLeaf(nil, nil), ShouldBeTrue)
		So(IsLeaf(nil, &Node{}), ShouldBeFalse)
		So(IsLeaf(&Node{}, &Node{}), ShouldBeFalse)
		So(IsLeaf(&Node{}, nil), ShouldBeFalse)
	})
}

func TestTrees(t *testing.T) {
	defer log.Flush()
	data := [][]byte{[]byte{0x00}, []byte{0x01}}
	z := sha256.Sum256(data[0])
	o := sha256.Sum256(data[1])
	var ZeroHash, OneHash, ZeroOne, ZeroZero []byte

	ZeroHash = z[:]
	OneHash = o[:]
	zo := sha256.Sum256(append(ZeroHash, OneHash...))
	ZeroOne = zo[:]
	zz := sha256.Sum256(append(ZeroHash, ZeroHash...))
	ZeroZero = zz[:]
	Convey("Construction", t, func() {
		tree := NewTree(data)
		So(tree.Root.Left, ShouldNotBeNil)
		So(tree.Root.Right, ShouldNotBeNil)
		So(tree.Root.Data, ShouldResemble, ZeroOne)
		So(tree.Root.Left.Left, ShouldBeNil)
		So(tree.Root.Left.Right, ShouldBeNil)
		So(tree.Root.Left.Data, ShouldResemble, ZeroHash)
		So(tree.Root.Right.Left, ShouldBeNil)
		So(tree.Root.Right.Right, ShouldBeNil)
		So(tree.Root.Right.Data, ShouldResemble, OneHash)
	})
	Convey("Asymmetric Tree", t, func() {
		tree := NewTree([][]byte{data[0]})
		So(tree.Root.Left, ShouldNotBeNil)
		So(tree.Root.Right, ShouldNotBeNil)
		So(tree.Root.Data, ShouldResemble, ZeroZero)
		So(tree.Root.Left.Left, ShouldBeNil)
		So(tree.Root.Left.Right, ShouldBeNil)
		So(tree.Root.Left.Data, ShouldResemble, ZeroHash)
		So(tree.Root.Right.Left, ShouldBeNil)
		So(tree.Root.Right.Right, ShouldBeNil)
		So(tree.Root.Right.Data, ShouldResemble, ZeroHash)
	})
	Convey("Biger tree", t, func() {
		tree := NewTree(append(data, data...))
		So(tree.Root.Left, ShouldNotBeNil)
		So(tree.Root.Right, ShouldNotBeNil)
		So(tree.Root.Left.Left, ShouldNotBeNil)
		So(tree.Root.Left.Right, ShouldNotBeNil)
		So(tree.Root.Right.Left, ShouldNotBeNil)
		So(tree.Root.Right.Right, ShouldNotBeNil)
		So(tree.Root.Left.Left.Left, ShouldBeNil)
		So(tree.Root.Left.Left.Right, ShouldBeNil)
		So(tree.Root.Left.Right.Left, ShouldBeNil)
		So(tree.Root.Left.Right.Right, ShouldBeNil)
		So(tree.Root.Right.Left.Left, ShouldBeNil)
		So(tree.Root.Right.Left.Right, ShouldBeNil)
		So(tree.Root.Right.Right.Left, ShouldBeNil)
		So(tree.Root.Right.Right.Right, ShouldBeNil)
	})

}
