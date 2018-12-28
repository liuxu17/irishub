package gov

import (
	"github.com/irisnet/irishub/modules/params"
	sdk "github.com/irisnet/irishub/types"

	"github.com/irisnet/irishub/modules/gov/params"
	"time"
)

const StartingProposalID = 1

// GenesisState - all gov state that must be provided at genesis
type GenesisState struct {
	SystemHaltPeriod  int64                       `json:"terminator_period"`
	DepositProcedure  govparams.DepositProcedure  `json:"deposit_period"`
	VotingProcedure   govparams.VotingProcedure   `json:"voting_period"`
	TallyingProcedure govparams.TallyingProcedure `json:"tallying_procedure"`
}

func NewGenesisState(dp govparams.DepositProcedure, vp govparams.VotingProcedure, tp govparams.TallyingProcedure) GenesisState {
	return GenesisState{
		DepositProcedure:  dp,
		VotingProcedure:   vp,
		TallyingProcedure: tp,
	}
}

// InitGenesis - store genesis parameters
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {

	err := k.setInitialProposalID(ctx, StartingProposalID)
	if err != nil {
		// TODO: Handle this with #870
		panic(err)
	}

	k.SetSystemHaltPeriod(ctx, data.SystemHaltPeriod)
	k.SetSystemHaltHeight(ctx, -1)

	params.InitGenesisParameter(&govparams.DepositProcedureParameter, ctx, data.DepositProcedure)
	params.InitGenesisParameter(&govparams.VotingProcedureParameter, ctx, data.VotingProcedure)
	params.InitGenesisParameter(&govparams.TallyingProcedureParameter, ctx, data.TallyingProcedure)
}

// ExportGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	depositProcedure := GetDepositProcedure(ctx)
	votingProcedure := GetVotingProcedure(ctx)
	tallyingProcedure := GetTallyingProcedure(ctx)

	return GenesisState{
		DepositProcedure:  depositProcedure,
		VotingProcedure:   votingProcedure,
		TallyingProcedure: tallyingProcedure,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		SystemHaltPeriod:  20000,
		DepositProcedure:  govparams.NewDepositProcedure(),
		VotingProcedure:   govparams.NewVotingProcedure(),
		TallyingProcedure: govparams.NewTallyingProcedure(),
	}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForCliTest() GenesisState {

	depositProcedure := govparams.NewDepositProcedure()
	depositProcedure.MaxDepositPeriod = time.Duration(60) * time.Second
	return GenesisState{
		SystemHaltPeriod:  20,
		DepositProcedure:  depositProcedure,
		VotingProcedure:   govparams.NewVotingProcedure(),
		TallyingProcedure: govparams.NewTallyingProcedure(),
	}
}

func PrepForZeroHeightGenesis(ctx sdk.Context, k Keeper) {
	proposals := k.GetProposalsFiltered(ctx, nil, nil, StatusDepositPeriod, 0)
	for _, proposal := range proposals {
		proposalID := proposal.GetProposalID()
		k.RefundDeposits(ctx, proposalID)
	}

	proposals = k.GetProposalsFiltered(ctx, nil, nil, StatusVotingPeriod, 0)
	for _, proposal := range proposals {
		proposalID := proposal.GetProposalID()
		k.RefundDeposits(ctx, proposalID)
	}
}
