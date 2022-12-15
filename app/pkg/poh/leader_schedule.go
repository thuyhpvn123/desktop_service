package poh

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"gitlab.com/meta-node/meta-node/pkg/pos"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const MAX_PERCENTAGE = 1000
const SLOT_PER_LEADER = 4

var (
	ErrTotalSlotNotDivisible = errors.New("total slot not divisible")
)

type LeaderSchedule struct {
	seed *uint256.Int

	fromSlot *uint256.Int
	toSlot   *uint256.Int

	slots map[uint256.Int]common.Address
}

func NewLeaderSchedule(
	seed *uint256.Int,
	stakePool *pos.StakePool,
	validatorStakeRequire *uint256.Int,
	fromSlot *uint256.Int,
	toSlot *uint256.Int,
) (*LeaderSchedule, error) {
	// from stake pool, filter out validators
	validatorsWithStake := stakePool.GetStakedAmountsAboveThreshHold(validatorStakeRequire)
	stakePercentages := calculateValidatorStakePercentage(validatorsWithStake)
	totalSlot := uint256.NewInt(0).Sub(toSlot, fromSlot)
	if totalSlot.Uint64()%SLOT_PER_LEADER != 0 {
		return nil, ErrTotalSlotNotDivisible
	}
	slots := calculateSlot(stakePercentages, seed, fromSlot, totalSlot.Uint64())
	// calculate slots
	ls := &LeaderSchedule{
		seed:     seed,
		fromSlot: fromSlot,
		toSlot:   toSlot,
		slots:    slots,
	}
	return ls, nil
}

func NewLeaderScheduleFromProto(protoLs *pb.LeaderSchedule) *LeaderSchedule {
	ls := &LeaderSchedule{
		seed:     uint256.NewInt(0).SetBytes(protoLs.Seed),
		fromSlot: uint256.NewInt(0).SetBytes(protoLs.FromSlot),
		toSlot:   uint256.NewInt(0).SetBytes(protoLs.ToSlot),
		slots:    make(map[uint256.Int]common.Address),
	}
	for i, v := range protoLs.Slots {
		ls.slots[*uint256.NewInt(0).SetBytes(common.FromHex(i))] = common.BytesToAddress(v)
	}
	return ls
}

func UnmarshalLeaderSchedule(b []byte) (*LeaderSchedule, error) {
	protoLs := &pb.LeaderSchedule{}
	err := proto.Unmarshal(b, protoLs)
	if err != nil {
		return nil, err
	}

	return NewLeaderScheduleFromProto(protoLs), nil
}

func (ls *LeaderSchedule) GetLeaderAtSlot(slot *uint256.Int) common.Address {
	return ls.slots[*slot]
}

func (ls *LeaderSchedule) GetToSlot() *uint256.Int {
	return ls.toSlot
}

func (ls *LeaderSchedule) SetSlots(slots map[uint256.Int]common.Address) {
	ls.slots = slots
}

func (ls *LeaderSchedule) SetToSlot(toSlot *uint256.Int) {
	ls.toSlot = toSlot
}

type Range struct {
	From uint64
	To   uint64
}

func calculateSlot(
	stakePercentage map[common.Address]uint64,
	seed *uint256.Int,
	fromSlot *uint256.Int,
	totalSlot uint64,
) map[uint256.Int]common.Address {
	rs := map[uint256.Int]common.Address{}
	// range from 0 to 1000
	rangeAddress := map[*Range]common.Address{}
	track := uint64(0)
	var lastRange *Range
	for address, percentage := range stakePercentage {
		rangeV := &Range{track, track + percentage}
		rangeAddress[rangeV] = address
		track += percentage
		lastRange = rangeV
	}
	if lastRange != nil && lastRange.To != MAX_PERCENTAGE {
		lastRange.To = MAX_PERCENTAGE
	}
	for rangeV, address := range rangeAddress {
		fmt.Printf("Range addresses from %v to %v address %v\n", rangeV.From, rangeV.To, address)
	}
	// for each slot use random from seed to get leader address, then update seed
	slotCount := fromSlot.Clone()
	toSlot := uint256.NewInt(0).AddUint64(slotCount, totalSlot)
	fmt.Printf("To slot %v\n", toSlot)
	for slotCount.Lt(toSlot) {
		rand.Seed(int64(seed.Uint64()))
		randValue := uint64(rand.Intn(1000))
		for rangeV, address := range rangeAddress {
			if randValue >= rangeV.From && randValue < rangeV.To {
				for u := uint64(0); u < SLOT_PER_LEADER; u++ {
					slotCount = slotCount.AddUint64(slotCount, 1)
					rs[*slotCount] = address
				}
				break
			}
		}
		seed = seed.AddUint64(seed, 1)
	}

	return rs
}

func calculateValidatorStakePercentage(validatorsWithStake map[common.Address]*uint256.Int) map[common.Address]uint64 {
	// calculate total staked of validators
	totalStake := uint256.NewInt(0)
	for _, v := range validatorsWithStake {
		totalStake = totalStake.Add(totalStake, v)
	}
	fmt.Printf("total stake %v\n", totalStake)
	// calculate /1000 staked of each validator
	stakePercentage := map[common.Address]uint64{}
	for k, v := range validatorsWithStake {
		percentage := uint256.NewInt(0).Mul(v, uint256.NewInt(MAX_PERCENTAGE))
		percentage = percentage.Div(percentage, totalStake)
		stakePercentage[k] = percentage.Uint64()
	}
	return stakePercentage
}

func (ls *LeaderSchedule) GetProto() protoreflect.ProtoMessage {
	protoLs := &pb.LeaderSchedule{
		Seed:     ls.seed.Bytes(),
		FromSlot: ls.fromSlot.Bytes(),
		ToSlot:   ls.toSlot.Bytes(),
	}
	return protoLs
}

func (ls *LeaderSchedule) Marshal() ([]byte, error) {
	return proto.Marshal(ls.GetProto())
}
