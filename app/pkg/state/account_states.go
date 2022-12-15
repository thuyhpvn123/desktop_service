package state

import (
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/meta-node/meta-node/pkg/merkle_patricia_trie"
	"gitlab.com/meta-node/meta-node/pkg/storage"
)

type IAccountStates interface {
	SetAccountState(accountState IAccountState) error
	GetAccountState(address common.Address) (IAccountState, error)
	GetAccountStateRoot() (common.Hash, error)
	Copy() IAccountStates
	Commit(savePath string) error
	GetStorageIterator() storage.IIterator
}

type AccountStates struct {
	storage storage.IStorage
	trie    *merkle_patricia_trie.Trie
}

func NewAccountStates(storage storage.IStorage, trie *merkle_patricia_trie.Trie) IAccountStates {
	return &AccountStates{
		storage: storage,
		trie:    trie,
	}
}

func (as *AccountStates) SetAccountState(accountState IAccountState) error {
	bAddress := accountState.GetAddress().Bytes()
	bData, err := accountState.Marshal()
	if err != nil {
		return err
	}
	as.trie.Set(bAddress, bData)
	return nil
}

func (as *AccountStates) GetAccountState(address common.Address) (IAccountState, error) {
	bAddress := address.Bytes()
	b, err := as.trie.Get(bAddress)
	if err != nil {
		return nil, err
	}
	return Unmarshal(b)
}

func (as *AccountStates) GetAccountStateRoot() (common.Hash, error) {
	_, rootHash, err := as.trie.HashRoot()
	return rootHash, err
}

func (as *AccountStates) Copy() IAccountStates {
	return &AccountStates{
		trie: as.trie.Copy(),
	}
}

func (as *AccountStates) Commit(savePath string) error {
	err := as.trie.Commit(savePath)
	return err
}

func (as *AccountStates) GetStorageIterator() storage.IIterator {
	return as.storage.GetIterator()
}
