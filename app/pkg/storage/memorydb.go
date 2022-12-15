package storage

import (
	"encoding/hex"
	"errors"
	fmt "fmt"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

type MemoryDB struct {
	db   map[string][]byte
	lock *sync.RWMutex
}

type MemoryDbIterator struct {
	db   map[string][]byte
	keys []string
	idx  int
}

func NewMemoryDbIterator(db map[string][]byte) *MemoryDbIterator {
	var keys []string
	for key, _ := range db {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return &MemoryDbIterator{
		db:   db,
		keys: keys,
		idx:  0,
	}
}

func (mdb *MemoryDB) GetIterator() IIterator {
	return NewMemoryDbIterator(mdb.db)
}

func (mdb *MemoryDbIterator) Next() bool {
	if len(mdb.keys) == mdb.idx {
		return false
	}
	mdb.idx += 1
	return true
}

func (mdb *MemoryDbIterator) Key() []byte {
	return common.FromHex(mdb.keys[mdb.idx-1])
}

func (mdb *MemoryDbIterator) Value() []byte {
	return mdb.db[mdb.keys[mdb.idx-1]]
}

func (mdb *MemoryDbIterator) Release() {
	//ToDo:
}
func (mdb *MemoryDbIterator) Error() error {
	// ToDo:
	return nil
}

func NewMemoryDb() *MemoryDB {
	return &MemoryDB{
		db:   make(map[string][]byte),
		lock: &sync.RWMutex{},
	}
}

func (kv *MemoryDB) Get(key []byte) ([]byte, error) {
	kv.lock.RLock()
	defer kv.lock.RUnlock()
	keyHex := hex.EncodeToString(key)
	if v, ok := kv.db[keyHex]; ok {
		return v, nil
	} else {
		return nil, errors.New(fmt.Sprintf("[MemKV] key not found: %s", keyHex))
	}
}

func (kv *MemoryDB) Put(key, value []byte) error {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	keyHex := hex.EncodeToString(key)
	kv.db[keyHex] = value
	return nil
}

func (kv *MemoryDB) Has(key []byte) bool {
	kv.lock.RLock()
	defer kv.lock.RUnlock()
	keyHex := hex.EncodeToString(key)
	_, ok := kv.db[keyHex]
	return ok
}

func (kv *MemoryDB) Delete(key []byte) error {
	kv.lock.Lock()
	defer kv.lock.Unlock()
	keyHex := hex.EncodeToString(key)
	if _, ok := kv.db[keyHex]; ok {
		delete(kv.db, keyHex)
	} else {
		return errors.New(fmt.Sprintf("[MemKV] key not found: %s", keyHex))
	}
	return nil
}

func (kv *MemoryDB) BatchPut(kvs [][2][]byte) error {
	for i := range kvs {
		kv.Put(kvs[i][0], kvs[i][1])
	}
	return nil
}

func (kv *MemoryDB) Close() {}
