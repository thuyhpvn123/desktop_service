package merkle_patricia_trie

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"gitlab.com/meta-node/meta-node/pkg/storage"
	"google.golang.org/protobuf/proto"
)

const TOTAL_CHILD_NODE = 16

type (
	Node interface {
		Hash() []byte
		CachedHash() []byte
		Serialize() []byte
		Save(storage.IStorage)
		StringWithPadding(padding string) string
		Flag() (bool)
		Deserialize( []byte,storage.IStorage) (err error)
 	}

	Flag struct {
		cacheHash []byte
		dirty     bool
	}

	FullNode struct {
		Children [TOTAL_CHILD_NODE]Node
		Value    []byte
		flag     Flag
	}

	ShortNode struct {
		Key   []byte
		Value Node
		flag  Flag
	}

	ValueNode struct {
		Value []byte
		flag  Flag
	}
)


func NewFlag() Flag {
	return Flag{
		dirty: true,
	}
}
	
func FetchNodeFromStorage(nodeHash []byte, storage storage.IStorage ) (Node,error){ //nodeHash là hash của node
	bN, err := storage.Get(nodeHash)
	if err != nil {
		return nil , err
	}
	protoNode:= &pb.MPTNode{}

	err = proto.Unmarshal(bN , protoNode)
	if err != nil {
		return nil , err
	}

	switch  protoNode.Type{
		case pb.MPTNODE_TYPE_FULL:
			fullNode := &FullNode{}
			err = fullNode.Deserialize(protoNode.Data,storage)
			fullNode.flag.dirty= false
			fullNode.flag.cacheHash = nodeHash
			fmt.Println(fullNode.StringWithPadding(""))	
			return fullNode, nil
		
		case pb.MPTNODE_TYPE_SHORT:
			shortNode := &ShortNode{}
			err = shortNode.Deserialize(protoNode.Data,storage)
			shortNode.flag.dirty= false
			shortNode.flag.cacheHash = nodeHash
			fmt.Println(shortNode.StringWithPadding(""))
			return shortNode, nil


		case pb.MPTNODE_TYPE_VALUE:
			valueNode := &ValueNode{}
			err = valueNode.Deserialize(protoNode.Data,storage)
			valueNode.flag.dirty = false
			valueNode.flag.cacheHash = nodeHash
			fmt.Println(valueNode.StringWithPadding(""))	

			return valueNode, nil
		default:
			return nil, nil

	}
}

// Full Node
func (n *FullNode) Flag() ( bool)  { return  n.flag.dirty }

func (n *FullNode) Hash() []byte {
	if n.flag.dirty {
		hash := crypto.Keccak256(n.Serialize())
		n.flag.cacheHash = hash[:]
		n.flag.dirty = false
		return hash
	}
	return n.flag.cacheHash
}

func (n *FullNode) CachedHash() []byte {
	return n.flag.cacheHash
}

func (n *FullNode) Serialize() []byte {
	bNodeHashes := [TOTAL_CHILD_NODE][]byte{}
	wg := sync.WaitGroup{}
	for i := 0; i < TOTAL_CHILD_NODE; i++ {
		if n.Children[i] == nil {
			bNodeHashes[i] = common.Hash{}.Bytes()
		} else {
			wg.Add(1)
			go func(idx int) {
				bNodeHashes[idx] = n.Children[idx].Hash()
				wg.Done()
			}(i)
		}
	}
	wg.Wait()
	protoFN := &pb.MPTFullNode{
		Nodes: bNodeHashes[:],
		Value: n.Value,
	}
	bFN, _ := proto.Marshal(protoFN)
	protoNode := &pb.MPTNode{
		Type: pb.MPTNODE_TYPE_FULL,
		Data: bFN,
	}
	bN, _ := proto.Marshal(protoNode)

	return bN
}
func (n*FullNode) Deserialize( bFN []byte,storage storage.IStorage) (err error) {
	protoFN := &pb.MPTFullNode{}
	err = proto.Unmarshal(bFN , protoFN)
	if err != nil {
		return err
	}
	for i, v := range protoFN.Nodes{
		n.Children[i],err = FetchNodeFromStorage(v,storage)
		if err != nil {
			return err
		}
		
	}
	return nil
}

func (n *FullNode) HasValue() bool {
	return len(n.Value) > 0
}
func (n *FullNode) SetValue(value []byte) {
	n.Value = value
	n.flag = NewFlag()
}
func (n *FullNode) Save(s storage.IStorage) {
	s.Put(n.Hash(), n.Serialize())
}

func (n *FullNode) copy() *FullNode { copy := *n; return &copy }

func (n *FullNode) String() string {
	hash := common.BytesToHash(n.flag.cacheHash)
	dirty := n.flag.dirty
	value := common.Bytes2Hex(n.Value)
	str := fmt.Sprintf("\u2514 FULLNODE hash:%v-dirty:%v-value:%v \n", hash, dirty, value)

	for _, v := range n.Children {
		str += fmt.Sprintf("	\u2514 %v\n", v)
	}

	return str
}

func (n *FullNode) StringWithPadding(padding string) string{
	// code like string function
	childPadding := padding +"\t"
	hash := common.BytesToHash(n.flag.cacheHash)
	dirty := n.flag.dirty
	value := common.Bytes2Hex(n.Value)
	str := fmt.Sprintf("\u2514 FULLNODE hash:%v-dirty:%v-value:%v \n", hash, dirty, value)

	for _, child := range n.Children {
		if child == nil{
			str += childPadding +"\u2514 <nil> \n"
		}else{
			str += child.StringWithPadding(childPadding)

		}
	}
	return fmt.Sprintf("%v%v", padding, str)
}


// Short Node

func (n *ShortNode) Flag() ( bool) { return n.flag.dirty }

func (n *ShortNode) Hash() []byte {
	if n.flag.dirty {
		hash := crypto.Keccak256(n.Serialize())
		n.flag.cacheHash = hash[:]
		n.flag.dirty = false
		return hash
	}
	return n.flag.cacheHash
}

func (n *ShortNode) CachedHash() []byte {
	return n.flag.cacheHash
}

func (n *ShortNode) Serialize() []byte {
	protoSN := &pb.MPTShortNode{
		Key:   n.Key,
		Value: n.Value.Hash(),
	}
	bSN, _ := proto.Marshal(protoSN)
	protoNode := &pb.MPTNode{
		Type: pb.MPTNODE_TYPE_SHORT,
		Data: bSN,
	}
	bN, _ := proto.Marshal(protoNode)
	return bN
}
func (n*ShortNode) Deserialize( bSN []byte,storage storage.IStorage)(err error) {
	protoSN := &pb.MPTShortNode{}
	err = proto.Unmarshal(bSN , protoSN)
	if err != nil {
		return err
	}
	childNode,err := FetchNodeFromStorage(protoSN.Value,storage)
	n.Value=childNode
	return nil

}
func (n *ShortNode) Save(s storage.IStorage) {
	s.Put(n.Hash(), n.Serialize())
}

func (n *ShortNode) copy() *ShortNode { copy := *n; return &copy }

func (n *ShortNode) String() string {
	hash := common.BytesToHash(n.flag.cacheHash)
	dirty := n.flag.dirty
	value := n.Value
	str := fmt.Sprintf("\u2514 SHORTNODE hash:%v-dirty:%v\r\n %v\n", hash, dirty, value)
	return str
}
func (n *ShortNode) StringWithPadding(padding string) string{
	hash := common.BytesToHash(n.flag.cacheHash)
	dirty := n.flag.dirty
	childPadding := padding + "\t"
	value := n.Value.StringWithPadding(childPadding)
	str := fmt.Sprintf("\u2514 SHORTNODE hash:%v-dirty:%v-value:\r\n %v\n", hash, dirty, value)
	
	return fmt.Sprintf("%v%v", padding, str)
}

// Value Node

func (n ValueNode) Flag() ( bool)  { return true }

func (n *ValueNode) Hash() []byte {
	if n.flag.dirty {
		hash := crypto.Keccak256(n.Value)
		n.flag.cacheHash = hash[:]
		n.flag.dirty = false

		return hash
	}

	return n.flag.cacheHash
}
func (n *ValueNode) CachedHash() []byte { return n.flag.cacheHash }
func (n *ValueNode) Serialize() []byte {
	protoNode := &pb.MPTNode{
		Type: pb.MPTNODE_TYPE_VALUE,
		Data: n.Value,
	}
	vN, _ := proto.Marshal(protoNode)
	return vN
}
func (n*ValueNode) Deserialize( value []byte,storage storage.IStorage)(err error) {
	n.Value = value	
	return nil
}
func (n *ValueNode) Save(s storage.IStorage) {
	s.Put(n.Hash(), n.Serialize())
}
func (n *ValueNode) String() string {
	hash := common.BytesToHash(n.flag.cacheHash)
	dirty := n.flag.dirty
	value := common.Bytes2Hex(n.Value)

	return fmt.Sprintf("\u2514 VALUENODE hash:%v-dirty:%v-value:%v\n", hash, dirty, value)
}
func (n *ValueNode) StringWithPadding(padding string) string{
	hash := common.BytesToHash(n.flag.cacheHash)
	dirty := n.flag.dirty
	value := common.Bytes2Hex(n.Value)
	str := fmt.Sprintf("\u2514 VALUENODE hash:%v-dirty:%v-value:%v\n", hash, dirty, value)
	
	return fmt.Sprintf("%v%v", padding, str)
}

