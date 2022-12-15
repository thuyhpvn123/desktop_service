package merkle_patricia_trie

import (
	"fmt"
	"testing"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"gitlab.com/meta-node/meta-node/pkg/storage"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
		"google.golang.org/protobuf/proto"


)

func initEmptyTrie() *Trie {
	return New(&FullNode{
		Value: nil,
		flag:  NewFlag(),
	}, nil)
}
func TestHashRoot(t *testing.T) {
	trie := initEmptyTrie()
	root, hash, err := trie.HashRoot()
	fmt.Printf("root %v\n hash %v\n err %v\n", root, hash, err)
	assert.Nil(t, err)
}

func TestTrieSetGet(t *testing.T) {
	trie := initEmptyTrie()
	// key := common.FromHex("f1f1f1")
	// value := common.FromHex("f2f2f2f2")
	key := common.FromHex("0")
	value := common.FromHex("f")

	err := trie.Set(key, value)
	assert.Nil(t, err)
	getValue, err := trie.Get(key)
	assert.Nil(t, err)
	assert.Equal(t, value, getValue)

	root, hash, err := trie.HashRoot()
	fmt.Printf("root %v\n hash %v\n err %v\n test %v\n", root, hash, err, crypto.Keccak256Hash(nil))
}
// func TestTrieDelete(t *testing.T) {
// 	trie := initEmptyTrie()
// 	befDel := fmt.Sprintf("%v",trie)
// 	fmt.Println(befDel)

// 	key := common.FromHex("f1f1f1")
// 	value := common.FromHex("f2f2f2f2")

// 	err := trie.Set(key, value)
// 	assert.Nil(t, err)

// 	fmt.Println(fmt.Sprintf("%v",trie))
// 	err2 := trie.Delete(key)
// 	assert.Nil(t,err2)
// 	aftDel :=fmt.Sprintf("%v",trie)
// 	fmt.Println(aftDel)
// 	assert.Equal(t, befDel, aftDel)

// }
func TestSetCommitGet(t *testing.T) {
	store := storage.NewMemoryDb()
	trie := New(&FullNode{
		Value: nil,
		flag:  NewFlag(),
	}, store)

	key := common.FromHex("123456")
	value := common.FromHex("A")
	trie.Set(key, value)
	trie.Commit(store)
	hashedKey := crypto.Keccak256(value)
	v, ok :=store.Get(hashedKey)
	protoNode := &pb.MPTNode{}
	err := proto.Unmarshal(v, protoNode)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, common.Bytes2Hex(protoNode.Data), "0a")
	assert.Nil(t,ok)

	err1 := trie.Set([]byte("134567"), []byte("B"))
	if err1 != nil {
		t.Error(err1.Error())
	}
	trie.Commit(store)
	hashedKey = crypto.Keccak256([]byte("B"))
	v1, ok1 :=store.Get(hashedKey)
	protoNode1 := &pb.MPTNode{}
	err = proto.Unmarshal(v1, protoNode1)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, string(protoNode1.Data), "B")
	assert.Nil(t,ok1)

}
// func TestFetchNodeFromStorage(t *testing.T) {
// 	store := storage.NewMemoryDb()
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

// 	trie := New(&fullNode, store)

// 	trie.Commit(store)
// 	_, err :=FetchNodeFromStorage(fullNode.Hash(),store)
// 	assert.Nil(t,err)

// }


