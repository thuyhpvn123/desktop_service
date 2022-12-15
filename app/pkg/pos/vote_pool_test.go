package pos

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"gitlab.com/meta-node/meta-node/pkg/vote"
)

var (
	testAddress1 = common.HexToAddress("067d82035bacafcf39258296bcbaa96ddf8678f6")
	testAddress2 = common.HexToAddress("d381f9537af2a152aa2fb2d721fe9b285f1e87f8")
	testAddress3 = common.HexToAddress("09a73b6984bf0d37f602fd31703b64433fa5f4d4")

	testPubkey1 = common.FromHex("a2702ce6bbfb2e013935781bac50a0e168732bd957861e6fbf185d688c82ade34c9f33fead179decb5953b3382b061df")
	testPubkey2 = common.FromHex("84260cdfc6e5ffe9c463212514193f2d5bd112d8e4868698d72cfb4cd003d749a01c0ec35fe729cc61e6301d816c4b30")
	testPubkey3 = common.FromHex("8557a68f0640d4fcaea396973d5647dde58aec8a414e6f56ad594e9994eb05cd124646ff87f07c151e47dbf30395b396")

	invalidSign = common.FromHex("111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111")
	testSign1   = common.FromHex("a507c03ab7ebb69a4b3adc22a0347bb2466788e6a3baa174a62bd74cdff60dfd6d6ba9ec6237098f1ceef6013bfeff1d0c8be716266710e1493c422293a676e7f168007324a23435d4590896f97f8e3686cf0c280240b9406800c1cec6bafb5d")
	testSign2   = common.FromHex("a56fe79b7619b15433e65bfcf85586ec997449f50e985ea2c59de9832a2ea37914c0b8b4b68e923a6d1aa4fd3d511819119e64479b9f30ee21488f0306a0b38e7b476ee956ec9cdab1006a57a3c272417cd387427f0a6f229992030760177801")
	testSign3   = common.FromHex("8b2a796d86a59788a8e568607c750c5e32ab843686316d085f7b9448b175d54f9f4ef82fb39d14dd7186966e7c4a81a10a73f4c9613b000c8257e9e0d9ed6bba257b19446a328b26d2daf00f43a3ea3c11073846cea0e47d62b2b96bfe22b3ec")

	testHash = common.HexToHash("1111111111111111111111111111111111111111111111111111111111111111")
)

func TestAddVote(t *testing.T) {
	stakePool := NewStakePool()

	stakePool.Stake(testAddress1, uint256.NewInt(10))
	stakePool.Stake(testAddress2, uint256.NewInt(5))
	stakePool.Stake(testAddress3, uint256.NewInt(10))

	votePool := NewVotePool(stakePool, uint256.NewInt(10), 2, 3)

	err := votePool.AddVote(vote.NewBlockVote(
		&pb.BlockVote{
			Hash:   testHash.Bytes(),
			Count:  []byte{1},
			Pubkey: testPubkey1,
			Sign:   invalidSign,
		},
	))
	assert.Equal(t, err, ErrInvalidSign)

	err = votePool.AddVote(vote.NewBlockVote(
		&pb.BlockVote{
			Hash:   testHash.Bytes(),
			Count:  []byte{1},
			Pubkey: testPubkey1,
			Sign:   testSign1,
		},
	))
	assert.Nil(t, err)
	assert.Nil(t, votePool.Result(), "Result must be nil")

	err = votePool.AddVote(vote.NewBlockVote(
		&pb.BlockVote{
			Hash:   testHash.Bytes(),
			Count:  []byte{1},
			Pubkey: testPubkey2,
			Sign:   testSign2,
		},
	))
	assert.Equal(t, err, ErrNotExistsInAddresses)
	assert.Nil(t, votePool.Result(), "Result must be nil")

	votePool.AddVote(vote.NewBlockVote(
		&pb.BlockVote{
			Hash:   testHash.Bytes(),
			Count:  []byte{1},
			Pubkey: testPubkey3,
			Sign:   testSign3,
		},
	))
	assert.Equal(t, *votePool.Result(), testHash)
}
