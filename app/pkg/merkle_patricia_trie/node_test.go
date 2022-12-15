package merkle_patricia_trie

import (
// 	// "bytes"
// 	"fmt"
// 	"sync"
// 	// "io"
// 	// "os"
// 	// "strings"
// 	"testing"

// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/stretchr/testify/assert"
// 	// pb "gitlab.com/meta-node/meta-node/pkg/proto"

// 	pb "gitlab.com/meta-node/meta-node/pkg/proto"
// 	"gitlab.com/meta-node/meta-node/pkg/storage"
// 	"google.golang.org/protobuf/proto"
// 	// "google.golang.org/protobuf/proto"
)

// func TestValueNodeString(t *testing.T) {
// 	valueNode := &ValueNode{
// 		Value: []byte("123"),
// 		flag : Flag{
// 			cacheHash :[]byte("456"),
// 			dirty     :false,
// 		},
// 	}
// 	fmt.Println(valueNode.StringWithPadding(""))
// }
// func TestShortNodeString(t *testing.T) {
// 	valueNode := &ValueNode{
// 		Value: []byte("123"),
// 		flag : Flag{
// 			cacheHash :[]byte("456"),
// 			dirty     :false,
// 		},
// 	}
// 	shortNode := ShortNode{
// 		Key:   []byte("789"),
// 		Value: valueNode,
// 		flag : Flag{
// 			cacheHash :[]byte("456"),
// 			dirty     :true,
// 		},
// 	}
// 	fmt.Println(shortNode.StringWithPadding(""))
// }
// func TestFullNodeString(t *testing.T) {
// 	valueNode := &ValueNode{
// 		Value: []byte("123"),
// 		flag : Flag{
// 			cacheHash :[]byte("456"),
// 			dirty     :false,
// 		},
// 	}
// 	shortNode := ShortNode{
// 		Key:   []byte("789"),
// 		Value: valueNode,
// 		flag : Flag{
// 			cacheHash :[]byte("456"),
// 			dirty     :true,
// 		},
// 	}
// 	fullNode := FullNode{}
// 	fullNode.Children[0] = &shortNode
// 	fullNode.Value = []byte("123")
// 	fullNode.flag.dirty = true
// 	fmt.Println(fullNode.StringWithPadding(""))
// }
// func TestNodeSerialize(t *testing.T) {
// 	// store := storage.NewMemoryDb()

// 	valueNode := &ValueNode{
// 		Value: []byte("123"),
// 		flag : NewFlag(),
// 	}
// 	//serialize valueNode
// 	kq := valueNode.Serialize()
// 	protoNode := &pb.MPTNode{
// 		Type: pb.MPTNODE_TYPE_VALUE,
// 		Data: []byte("123"),
// 	}
// 	vN, _ := proto.Marshal(protoNode)

// 	assert.Equal(t, vN, kq)

// 	//serialize shortNode
// 	shortNode := ShortNode{
// 		Key:   []byte("789"),
// 		Value: valueNode,
// 		flag : NewFlag(),
// 	}
// 	kq1 := shortNode.Serialize()

// 	protoSN := &pb.MPTShortNode{
// 		Key:   []byte("789"),
// 		Value: valueNode.Hash(),
// 	}
// 	bSN, _ := proto.Marshal(protoSN)
// 	protoNode1 := &pb.MPTNode{
// 		Type: pb.MPTNODE_TYPE_SHORT,
// 		Data: bSN,
// 	}
// 	bN, _ := proto.Marshal(protoNode1)
// 	assert.Equal(t, kq1, bN)

// 	//serialize FullNode
// 	fullNode := FullNode{}
// 	fullNode.Children[0] = &shortNode
// 	fullNode.Value = []byte("123")
// 	fullNode.flag.dirty = true

// 	kq2 := fullNode.Serialize()
// 	bNodeHashes := [16][]byte{}
// 	wg := sync.WaitGroup{}
// 	for i := 0; i < TOTAL_CHILD_NODE; i++ {
// 		if fullNode.Children[i] == nil {
// 			bNodeHashes[i] = common.Hash{}.Bytes()
// 		} else {
// 			wg.Add(1)
// 			go func(idx int) {
// 				bNodeHashes[idx] = common.Hash{}.Bytes()
// 				bNodeHashes[idx] = fullNode.Children[idx].Hash()
// 				wg.Done()
// 			}(i)
// 		}
// 	}
// 	wg.Wait()
// 	protoFN := &pb.MPTFullNode{
// 		Nodes: bNodeHashes[:],
// 		Value: fullNode.Value,
// 	}
// 	bFN, _ := proto.Marshal(protoFN)
// 	protoNode2 := &pb.MPTNode{
// 		Type: pb.MPTNODE_TYPE_FULL,
// 		Data: bFN,
// 	}
// 	bN3, _ := proto.Marshal(protoNode2)
// 	assert.Equal(t, kq2, bN3)

// }
// func TestNodeDeserialize(t *testing.T) {
// store := storage.NewMemoryDb()

// valueNode := &ValueNode{
// 	Value: []byte("123"),
// 	flag : NewFlag(),
// }
// valueNode.Save(store)

// //Deserialize valueNode
// vN:=valueNode.Serialize()
// protoNode := &pb.MPTNode{}
// err := proto.Unmarshal(vN, protoNode)
// assert.Nil(t,err)

// valueNode1 :=&ValueNode{}
// err=valueNode1.Deserialize(protoNode.Data,store)
// assert.Equal(t, valueNode1.Value, []byte("123"))

// //Deserialize shortNode
// shortNode := ShortNode{
// 	Key:   []byte("789"),
// 	Value: valueNode,
// 	flag : NewFlag(),
// }
// shortNode.Save(store)
// bN:=shortNode.Serialize()
// protoNode1 := &pb.MPTNode{}
// err = proto.Unmarshal(bN, protoNode1)
// shortNode1 :=&ShortNode{}
// err=shortNode1.Deserialize(protoNode1.Data,store)
// valueNode2 := shortNode1.Value
// assert.Equal(t, valueNode2.(*ValueNode).Value, []byte("123"))

// //Deserialize FullNode
// fullNode := FullNode{}
// fullNode.Children[0] = &shortNode
// fullNode.Value = []byte("")
// fullNode.flag.dirty = true

// fullNode.Save(store)
// bN2:=fullNode.Serialize()
// protoNode2 := &pb.MPTNode{}
// err = proto.Unmarshal(bN2, protoNode2)
// fullNode2 :=&FullNode{}
// err=fullNode2.Deserialize(protoNode2.Data,store)
// sN := fullNode2.Children[0]
// vN3:=sN.(*ShortNode).Value
// assert.Equal(t, vN3.(*ValueNode).Value, []byte("123"))


// }


