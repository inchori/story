package types

import (
	"fmt"
)

// Staking params default values.
const (
	DefaultMaxWithdrawalPerBlock uint32 = 4

	DefaultMaxSweepPerBlock uint32 = 64

	DefaultMinPartialWithdrawalAmount uint64 = 600_000

	DefaultSingularityHeight uint64 = 1209600 // 42 days with 35 seconds block time
)

// NewParams creates a new Params instance.
func NewParams(maxWithdrawalPerBlock uint32, maxSweepPerBlock uint32, minPartialWithdrawalAmount, singularityHeight uint64) Params {
	return Params{
		MaxWithdrawalPerBlock:      maxWithdrawalPerBlock,
		MaxSweepPerBlock:           maxSweepPerBlock,
		MinPartialWithdrawalAmount: minPartialWithdrawalAmount,
		SingularityHeight:          singularityHeight,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultMaxWithdrawalPerBlock,
		DefaultMaxSweepPerBlock,
		DefaultMinPartialWithdrawalAmount,
		DefaultSingularityHeight,
	)
}

func ValidateMaxWithdrawalPerBlock(v uint32) error {
	if v == 0 {
		return fmt.Errorf("max withdrawal per block must be positive: %d", v)
	}

	return nil
}

func ValidateMaxSweepPerBlock(maxSweepPerBlock uint32, maxWithdrawalPerBlock uint32) error {
	if maxSweepPerBlock == 0 {
		return fmt.Errorf("max sweep per block must be positive: %d", maxSweepPerBlock)
	}

	if maxSweepPerBlock < maxWithdrawalPerBlock {
		return fmt.Errorf("max sweep per block must be greater than or equal to max withdrawal per block: %d < %d", maxSweepPerBlock, maxWithdrawalPerBlock)
	}

	return nil
}

func ValidateMinPartialWithdrawalAmount(v uint64) error {
	if v == 0 {
		return fmt.Errorf("min partial withdrawal amount must be positive: %d", v)
	}

	return nil
}

func ValidateSingularityHeight(v uint64) error {
	if v == 0 {
		return fmt.Errorf("singularity height must be positive: %d", v)
	}

	return nil
}
