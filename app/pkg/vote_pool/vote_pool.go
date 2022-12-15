package vote_pool

import (
	"errors"
	"math"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/meta-node/meta-node/pkg/bls"
	cm "gitlab.com/meta-node/meta-node/pkg/common"
	"gitlab.com/meta-node/meta-node/pkg/network"
	v "gitlab.com/meta-node/meta-node/pkg/vote"
)

var (
	ErrNotExistsInAddresses = errors.New("not exist in addresses")
	ErrAlreadyVoted         = errors.New("already voted")
	ErrInvalidSign          = errors.New("invalid sign")
)

// vote pool with count weight
type VotePool struct {
	approveRate    float64
	addresses      map[common.Address]network.IConnection   // use to track participators
	votes          map[common.Hash]map[cm.PublicKey]cm.Sign // vote hash => addresses
	mapAddressVote map[common.Address]common.Hash
	voteValues     map[common.Hash]interface{}
	result         *common.Hash

	voteMu sync.RWMutex
}

func NewVotePool(
	approveRate float64,
	addresses map[common.Address]network.IConnection,
) *VotePool {
	return &VotePool{
		approveRate:    approveRate,
		addresses:      addresses,
		votes:          make(map[common.Hash]map[cm.PublicKey]cm.Sign),
		mapAddressVote: make(map[common.Address]common.Hash),
		voteValues:     make(map[common.Hash]interface{}),
	}
}

func (p *VotePool) AddVote(vote v.IVote) error {
	p.voteMu.Lock()
	defer p.voteMu.Unlock()
	pk := vote.GetPublicKey()
	sign := vote.GetSign()
	hash := vote.GetHash()
	value := vote.GetValue()
	address := cm.GetAddressFromPubkey(pk)
	if !bls.VerifySign(pk, sign, hash.Bytes()) {
		return ErrInvalidSign
	}
	if v, ok := p.addresses[address]; !ok || v == nil {
		return ErrNotExistsInAddresses
	}

	if _, ok := p.mapAddressVote[address]; ok {
		return ErrAlreadyVoted
	}

	p.mapAddressVote[address] = hash
	if p.votes[hash] == nil {
		p.votes[hash] = make(map[cm.PublicKey]cm.Sign)
	}
	p.votes[hash][pk] = sign
	if value != nil {
		p.voteValues[hash] = value
	}
	p.checkVote(hash)
	return nil
}

func (p *VotePool) checkVote(voteHash common.Hash) bool {
	countVotes := len(p.votes[voteHash])
	//
	requireVotes := int(math.Ceil(float64(len(p.addresses)) * p.approveRate))
	//
	if countVotes >= requireVotes {
		p.result = &voteHash
		return true
	}
	return false
}

func (p *VotePool) Addresses() map[common.Address]network.IConnection {
	return p.addresses
}

func (p *VotePool) Result() *common.Hash {
	p.voteMu.RLock()
	defer p.voteMu.RUnlock()
	return p.result
}

func (p *VotePool) ResultValue() interface{} {
	p.voteMu.RLock()
	defer p.voteMu.RUnlock()
	return p.voteValues[*p.result]
}
