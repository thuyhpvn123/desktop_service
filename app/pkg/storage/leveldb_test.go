package storage

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

var ldb *LevelDB

func initLevelDB() *LevelDB {
	path := "./ldb_test"
	ldb, _ = NewLevelDB(path)
	return ldb
}

func TestInitLevelDB(t *testing.T) {
	ldb = initLevelDB()
	assert.NotEmpty(t, ldb)
}

func TestPut(t *testing.T) {
	err1 := ldb.Put(common.FromHex("0x001"), common.FromHex("0xf1f1"))
	err2 := ldb.Put(common.FromHex("0x002"), common.FromHex("0xh2"))
	err3 := ldb.Put(common.FromHex("0x003"), common.FromHex("0xh3h3"))
	assert.Empty(t, err1)
	assert.Empty(t, err2)
	assert.Empty(t, err3)
}

func TestGet(t *testing.T) {
	db1, _ := ldb.Get(common.FromHex("0x001"))
	assert.Equal(t, db1, common.FromHex("0xf1f1"))
	db2, _ := ldb.Get(common.FromHex("0x002"))
	assert.Equal(t, db2, common.FromHex("0xh2"))
	db3, _ := ldb.Get(common.FromHex("0x003"))
	assert.Equal(t, db3, common.FromHex("0xh3h3"))
	db4, _ := ldb.Get(common.FromHex("0x004"))
	assert.Empty(t, db4)
}

func TestHas(t *testing.T) {
	chk1 := ldb.Has(common.FromHex("0x001"))
	chk2 := ldb.Has(common.FromHex("0x002"))
	chk3 := ldb.Has(common.FromHex("0x003"))
	chk4 := ldb.Has(common.FromHex("0x004"))
	assert.True(t, chk1)
	assert.True(t, chk2)
	assert.True(t, chk3)
	assert.False(t, chk4)
}

func TestDelete(t *testing.T) {
	ldb.Put(common.FromHex("0x004"), common.FromHex("0x00f0ff"))
	chk1 := ldb.Has(common.FromHex("0x004"))
	assert.True(t, chk1)
	ldb.Delete(common.FromHex("0x004"))
	chk2 := ldb.Has(common.FromHex("0x004"))
	assert.False(t, chk2)
}

func TestBatchPut(t *testing.T) {
	arr := [][2][]byte{
		{
			common.FromHex("0x001"),
			common.FromHex("0xh1111"),
		},
		{
			common.FromHex("0x002"),
			common.FromHex("0xf1f2"),
		},
		{
			common.FromHex("0x003"),
			common.FromHex("0x0000"),
		},
	}
	assert.NoError(t, ldb.BatchPut(arr))
	db1, _ := ldb.Get(common.FromHex("0x001"))
	assert.Equal(t, db1, common.FromHex("0xh1111"))
	db2, _ := ldb.Get(common.FromHex("0x002"))
	assert.Equal(t, db2, common.FromHex("0xf1f2"))
	db3, _ := ldb.Get(common.FromHex("0x003"))
	assert.Equal(t, db3, common.FromHex("0x0000"))
}

func TestGetIterator(t *testing.T) {
	iterator := ldb.GetIterator()
	assert.NotEmpty(t, iterator)
}
