package mint

import (
	"testing"
	"time"

	stakeTypes "github.com/irisnet/irishub/modules/stake/types"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

func TestNextInflation(t *testing.T) {
	minter := NewMinter(time.Now(), stakeTypes.StakeDenom, sdk.NewIntWithDecimal(100, 18))
	tests := []struct {
		params        Params
		blockInterval time.Duration
	}{
		{DefaultParams(), 5},
		{DefaultParams(), 10},
		{DefaultParams(), 20},
		{DefaultParams(), 100},
		{Params{sdk.NewDecWithPrec(20, 2)}, 5},
		{Params{sdk.NewDecWithPrec(10, 2)}, 5},
		{Params{sdk.NewDecWithPrec(5, 2)}, 5},
	}
	for _, tc := range tests {
		time.Sleep(tc.blockInterval * time.Millisecond)
		blockTime := time.Now()
		annualProvisions := minter.NextAnnualProvisions(tc.params)
		mintCoin := minter.BlockProvision(tc.params, annualProvisions, blockTime)

		blockDurationMili := tc.blockInterval.Nanoseconds() / int64(miliSecondPerYear)
		blockTimePercent := sdk.NewDec(blockDurationMili).Quo(sdk.NewDec(int64(60 * 60 * 8766 * 1000)))

		require.True(t, mintCoin.Amount.GT(blockTimePercent.MulInt(minter.InflationBase).TruncateInt()))
	}
}
