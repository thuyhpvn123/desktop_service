package config

import "github.com/holiman/uint256"

type IConsensusConfig interface {
	GetPacksPerEntry() uint64
	GetEntriesPerSlot() uint64
	GetHashesPerEntry() uint64
	GetEntriesPerSecond() uint64
	GetValidatorMinStakeAmount() *uint256.Int

	GetValidatorVoteApproveRate() (uint64, uint64)
	GetNodeVoteApproveRate() float64
}
