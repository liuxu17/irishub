package service

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/service/params"
	"github.com/irisnet/irishub/modules/params"
)

// GenesisState - all service state that must be provided at genesis
type GenesisState struct {
	MaxRequestTimeout  int64
	MinDepositMultiple int64
}

func NewGenesisState(maxRequestTimeout int64, minDepositMultiple int64) GenesisState {
	return GenesisState{
		MaxRequestTimeout:  maxRequestTimeout,
		MinDepositMultiple: minDepositMultiple,
	}
}

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, data GenesisState) {
	params.InitGenesisParameter(&serviceparams.MaxRequestTimeoutParameter, ctx, data.MaxRequestTimeout)
	params.InitGenesisParameter(&serviceparams.MinDepositMultipleParameter, ctx, data.MinDepositMultiple)
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context) GenesisState {
	maxRequestTimeout := serviceparams.GetMaxRequestTimeout(ctx)
	minDepositMultiple := serviceparams.GetMinDepositMultiple(ctx)

	return GenesisState{
		MaxRequestTimeout:  maxRequestTimeout,
		MinDepositMultiple: minDepositMultiple,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		MaxRequestTimeout:  100,
		MinDepositMultiple: 1000,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	return GenesisState{
		MaxRequestTimeout:  10,
		MinDepositMultiple: 10,
	}
}
