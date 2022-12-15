package test_structs

import "github.com/holiman/uint256"

type TestConsensusConfig struct {
	PacksPerEntry               uint64
	EntriesPerSlot              uint64
	HashesPerEntry              uint64
	EntriesPerSecond            uint64
	ValidatorMinStakeAmount     *uint256.Int
	ValidatorVoteApproveRateMul uint64
	ValidatorVoteApproveRateDiv uint64
	NodeVoteApproveRate         float64
}

func (c *TestConsensusConfig) GetPacksPerEntry() uint64 {
	return c.PacksPerEntry
}
func (c *TestConsensusConfig) GetEntriesPerSlot() uint64 {
	return c.EntriesPerSlot
}
func (c *TestConsensusConfig) GetHashesPerEntry() uint64 {
	return c.HashesPerEntry
}
func (c *TestConsensusConfig) GetEntriesPerSecond() uint64 {
	return c.EntriesPerSecond
}
func (c *TestConsensusConfig) GetValidatorMinStakeAmount() *uint256.Int {
	return c.ValidatorMinStakeAmount
}
func (c *TestConsensusConfig) GetValidatorVoteApproveRate() (uint64, uint64) {
	return c.ValidatorVoteApproveRateMul, c.ValidatorVoteApproveRateDiv
}
func (c *TestConsensusConfig) GetNodeVoteApproveRate() float64 {
	return c.NodeVoteApproveRate
}
