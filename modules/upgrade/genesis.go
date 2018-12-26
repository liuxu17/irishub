package upgrade

import (
	"github.com/irisnet/irishub/app/protocol/keeper"
	"github.com/irisnet/irishub/modules/params"
	"github.com/irisnet/irishub/modules/upgrade/params"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/version"
)

// GenesisState - all upgrade state that must be provided at genesis
type GenesisState struct {
	GenesisVersion AppVersion           `json:genesis_version`
	UpgradeParams  upgradeparams.Params `json:upgrade_govparams`
}

// InitGenesis - build the genesis version For first Version
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	genesisVersion := data.GenesisVersion

	k.AddNewVersion(ctx, genesisVersion)
	k.pk.ClearUpgradeConfig(ctx)
	k.pk.SetCurrentProtocolVersion(ctx, genesisVersion.Protocol.Version)
	params.InitGenesisParameter(&upgradeparams.UpgradeParameter, ctx, data.UpgradeParams)
}

// WriteGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context) GenesisState {
	return GenesisState{
		GenesisVersion: NewVersion(
			keeper.UpgradeConfig{0,
				sdk.ProtocolDefinition{
					uint64(0),
					" ",
					uint64(1),
				}}, true),
		UpgradeParams: upgradeparams.NewUpgradeParams()}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		GenesisVersion: NewVersion(
			keeper.UpgradeConfig{0,
				sdk.ProtocolDefinition{
					uint64(0),
					"https://github.com/irisnet/irishub/releases/tag/v" + version.Version,
					uint64(1),
				}}, true),
		UpgradeParams: upgradeparams.NewUpgradeParams()}
}

// get raw genesis raw message for testing
func DefaultGenesisStateForTest() GenesisState {
	return GenesisState{
		GenesisVersion: NewVersion(
			keeper.UpgradeConfig{0,
				sdk.ProtocolDefinition{
					uint64(0),
					"https://github.com/irisnet/irishub/releases/tag/v" + version.Version,
					uint64(1),
				}}, true),
		UpgradeParams: upgradeparams.NewUpgradeParams()}
}
