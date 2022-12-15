package merkle_patricia_trie

import (
	"bytes"
	"errors"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/meta-node/meta-node/pkg/storage"
)

var (
	// emptyRoot is the known root hash of an empty trie.
	emptyRoot = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

	// emptyState is the known hash of an empty state trie entry.
	emptyState = crypto.Keccak256Hash(nil)
)

type Trie struct {
	oldRoot []byte
	root    Node
	lock    *sync.RWMutex
}



func New(root Node, store storage.IStorage) *Trie {
	var oldRoot []byte = nil
	if root != nil {
		oldRoot = root.Hash()
	}
	return &Trie{
		oldRoot: oldRoot,
		root:    root,
		lock:    &sync.RWMutex{},
	}
}
func (t *Trie) String() string {
	root, hash, err := t.HashRoot()

	return fmt.Sprintf("Trie:\n root:%v\n hash:%v\n err:%v\n", root.StringWithPadding(""), hash, err)
}
func (t *Trie) Get(key []byte) ([]byte, error) {
	hashedKey := crypto.Keccak256(key)
	k := keybytesToHex(hashedKey)
	value, _, err := t.get(t.root, k, 0)
	return value, err
}

func (t *Trie) get(origNode Node, key []byte, pos int) (value []byte, newnode Node, err error) {
	switch n := (origNode).(type) {
	case nil:
		return nil, nil, nil
	case *ValueNode:
		return n.Value, n, nil
	case *ShortNode:
		if len(key)-pos < len(n.Key) || !bytes.Equal(n.Key, key[pos:pos+len(n.Key)]) {
			// key not found in trie
			return nil, n, nil
		}
		value, newnode, err = t.get(n.Value, key, pos+len(n.Key))
		if err == nil {
			n = n.copy()
			n.Value = newnode
		}
		return value, n, err
	case *FullNode:
		value, newnode, err = t.get(n.Children[key[pos]], key, pos+1)
		if err == nil {
			n = n.copy()
			n.Children[key[pos]] = newnode
		}
		return value, n, err
	default:
		panic(fmt.Sprintf("%T: invalid node: %v", origNode, origNode))
	}
}

func (t *Trie) Set(key, value []byte) error {
	hashedKey := crypto.Keccak256(key)
	k := keybytesToHex(hashedKey)
	if len(value) != 0 {
		_, n, err := t.insert(t.root, nil, k, &ValueNode{
			Value: value,
			flag : NewFlag(),
		})
		if err != nil {
			return err
		}	

		t.root = n
	}
	// TODO: if len == 0 then delete node
	return nil
}

func (t *Trie) insert(n Node, prefix, key []byte, value Node) (bool, Node, error) {
	if len(key) == 0 {
		if v, ok := n.(*ValueNode); ok {
			return !bytes.Equal(v.Value, value.(*ValueNode).Value), value, nil
		}
		return true, value, nil
	}
	switch n := n.(type) {
	case *ShortNode:
		matchlen := prefixLen(key, n.Key)
		// If the whole key matches, keep this short node as is
		// and only update the value.
		if matchlen == len(n.Key) {
			dirty, nn, err := t.insert(n.Value, append(prefix, key[:matchlen]...), key[matchlen:], value)
			if !dirty || err != nil {
				return false, n, err
			}
			return true, &ShortNode{
				n.Key,
				nn,
				NewFlag(),
			}, nil
		}
		// Otherwise branch out at the index where they differ.
		branch := &FullNode{flag: NewFlag()}
		var err error
		_, branch.Children[n.Key[matchlen]], err = t.insert(nil, append(prefix, n.Key[:matchlen+1]...), n.Key[matchlen+1:], n.Value)
		if err != nil {
			return false, nil, err
		}
		_, branch.Children[key[matchlen]], err = t.insert(nil, append(prefix, key[:matchlen+1]...), key[matchlen+1:], value)
		if err != nil {
			return false, nil, err
		}
		// Replace this shortNode with the branch if it occurs at index 0.
		if matchlen == 0 {
			return true, branch, nil
		}

		// Replace it with a short node leading up to the branch.
		return true, &ShortNode{key[:matchlen], branch, NewFlag()}, nil

	case *FullNode:
		dirty, nn, err := t.insert(n.Children[key[0]], append(prefix, key[0]), key[1:], value)
		if !dirty || err != nil {
			return false, n, err
		}
		n = n.copy()
		n.flag = NewFlag()
		n.Children[key[0]] = nn
		fmt.Printf("HỂ\n")
		return true, n, nil

	case nil:
		return true, &ShortNode{
			Key:   key,
			Value: value,
			flag:  NewFlag(),
		}, nil

	default:
		panic(fmt.Sprintf("%T: invalid node: %v", n, n))
	}
}
func (t *Trie) Delete(key []byte) (error) {
	hashedKey := crypto.Keccak256(key)
	k := keybytesToHex(hashedKey)
	_,n,err :=t.delete(t.root, nil, k)

	if err != nil{
		return err
	}
	t.root = n 
	return  nil
}

// delete returns the new root of the trie with key deleted.
// It reduces the trie to minimal form by simplifying
// nodes on the way up after deleting recursively.
func (t *Trie) delete(origNode Node, prefix, key []byte) (b bool, newnode Node, err error) {
	switch n := (origNode).(type) {
	case *ShortNode:

		matchlen := prefixLen(key, n.Key)
		if matchlen < len(n.Key) {
			return false, n, nil // don't replace n on mismatch
		}
		if matchlen == len(key) {
			// The matched short node is deleted entirely and track
			// it in the deletion set. The same the valueNode doesn't
			// need to be tracked at all since it's always embedded.
	
			return true, nil, nil // remove n entirely for whole matches
		}
		// The key is longer than n.Key. Remove the remaining suffix
		// from the subtrie. Child can never be nil here since the
		// subtrie must contain at least two other values with keys
		// longer than n.Key.
		dirty, child, err := t.delete(n.Value, append(prefix, key[:len(n.Key)]...), key[len(n.Key):])

		if !dirty || err != nil {
			return false, n, err
		}
		switch child := child.(type) {
		case *ShortNode:
			// The child shortNode is merged into its parent, track
			// is deleted as well.

			// Deleting from the subtrie reduced it to another
			// short node. Merge the nodes to avoid creating a
			// shortNode{..., shortNode{...}}. Use concat (which
			// always creates a new slice) instead of append to
			// avoid modifying n.Key since it might be shared with
			// other nodes.
			return true, &ShortNode{concat(n.Key, child.Key...), child.Value, NewFlag()}, nil
		default:
			return true, &ShortNode{n.Key, child, NewFlag()}, nil
		}

	case *FullNode:
		dirty, nn, err := t.delete(n.Children[key[0]], append(prefix, key[0]), key[1:])
		if !dirty || err != nil {
			return false, n, err
		}
		n = n.copy()
		n.flag = NewFlag()
		n.Children[key[0]] = nn

		// Because n is a full node, it must've contained at least two children
		// before the delete operation. If the new child value is non-nil, n still
		// has at least two children after the deletion, and cannot be reduced to
		// a short node.
		if nn != nil {
			return true, n, nil
		}
		// Reduction:
		// Check how many non-nil entries are left after deleting and
		// reduce the full node to a short node if only one entry is
		// left. Since n must've contained at least two children
		// before deletion (otherwise it would not be a full node) n
		// can never be reduced to nil.
		//
		// When the loop is done, pos contains the index of the single
		// value that is left in n or -2 if n contains at least two
		// values.
		pos := -1
		for i, cld := range &n.Children {
			if cld != nil {
				if pos == -1 {
					pos = i
				} else {
					pos = -2
					break
				}
			}
		}
		if pos >= 0 {
			if pos != 16 {
				// If the remaining entry is a short node, it replaces
				// n and its key gets the missing nibble tacked to the
				// front. This avoids creating an invalid
				// shortNode{..., shortNode{...}}.  Since the entry
				// might not be loaded yet, resolve it just for this
				// check.
				cnode, err := t.resolve(n.Children[pos], append(prefix, byte(pos)))
				if err != nil {
					return false, nil, err
				}
				if cnode, ok := cnode.(*ShortNode); ok {
					// Replace the entire full node with the short node.
					// Mark the original short node as deleted since the
					// value is embedded into the parent now.
			
					k := append([]byte{byte(pos)}, cnode.Key...)
					return true, &ShortNode{k, cnode.Value, NewFlag()}, nil
				}
			}
			// Otherwise, n is replaced by a one-nibble short node
			// containing the child.
			return true, &ShortNode{[]byte{byte(pos)}, n.Children[pos], NewFlag()}, nil
		}
		// n still contains at least two values and cannot be reduced.
		return true, n, nil

	case *ValueNode:

		return true, nil, nil

	case nil:
		return false, nil, nil

	default:
		panic(fmt.Sprintf("%T: invalid node: %v (%v)", n, n, key))
	}
}
func concat(s1 []byte, s2 ...byte) []byte {
	r := make([]byte, len(s1)+len(s2))
	copy(r, s1)
	copy(r[len(s1):], s2)
	return r
}
func (t *Trie) resolve(n Node, prefix []byte) (Node, error) {
	// if n, ok := n.(hashNode); ok {
	// 	return t.resolveAndTrack(n, prefix)
	// }
	return n, nil
}

// Commit collects all dirty nodes in the trie and replaces them with the
// corresponding node hash. All collected nodes (including dirty leaves if
// collectLeaf is true) will be encapsulated into a nodeset for return.
// The returned nodeset can be nil if the trie is clean (nothing to commit).
// Once the trie is committed, it's not usable anymore. A new trie must
// be created with new root and updated trie database for following usage

func (t *Trie) Commit(storage storage.IStorage) {
	t.lock.Lock()
	defer t.lock.Unlock()
	if t.root == nil {
		return
	}
	t.commit(t.root,storage)
	t.oldRoot = t.root.CachedHash()
}


func (t *Trie) commit(node Node,storage storage.IStorage) {

	switch n := node.(type) {
	case *FullNode:
		for i := 0; i < len(n.Children); i++ {
			t.commit(n.Children[i],storage)
		}
		n.Save(storage)

	case *ShortNode:
		t.commit(n.Value,storage)
		n.Save(storage)

	case *ValueNode:
		n.Save(storage)

	}
}
// hashRoot calculates the root hash of the given trie
func (t *Trie) HashRoot() (Node, common.Hash, error) {
	if t.root == nil {
		return nil, common.Hash{}, errors.New("nil root")
	}
	// If the number of changes is below 100, we let one thread handle it
	hashedRoot := t.root.Hash()
	return t.root, common.BytesToHash(hashedRoot), nil
}
