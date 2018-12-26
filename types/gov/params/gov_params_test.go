package govparams

import (
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/store"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/params"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
	"fmt"
	"time"
)

func defaultContext(key sdk.StoreKey, tkeyParams *sdk.TransientStoreKey) sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)

	cms.LoadLatestVersion()
	ctx := sdk.NewContext(cms, abci.Header{}, false, log.NewNopLogger())
	return ctx
}

func TestInitGenesisParameter(t *testing.T) {
	skey := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	ctx := defaultContext(skey, tkeyParams)
	cdc := codec.New()

	paramKeeper := params.NewKeeper(
		cdc,
		skey, tkeyParams,
	)

	p1 := NewDepositProcedure()
	p2 := NewDepositProcedure()

	subspace := paramKeeper.Subspace("Gov").WithTypeTable(
		params.NewTypeTable(
			DepositProcedureParameter.GetStoreKey(), DepositProcedure{},
			VotingProcedureParameter.GetStoreKey(), VotingProcedure{},
			TallyingProcedureParameter.GetStoreKey(), TallyingProcedure{},
		))
	params.SetParamReadWriter(subspace, &DepositProcedureParameter, &DepositProcedureParameter)
	params.InitGenesisParameter(&DepositProcedureParameter, ctx, nil)

	fmt.Println(DepositProcedureParameter.ToJson(""))
	require.Equal(t, p1, DepositProcedureParameter.Value)
	require.Equal(t, DepositProcedureParameter.ToJson(""), `{"critical_min_deposit":[{"denom":"iris-atto","amount":"4000000000000000000000"}],"important_min_deposit":[{"denom":"iris-atto","amount":"2000000000000000000000"}],"normal_min_deposit":[{"denom":"iris-atto","amount":"1000000000000000000000"}],"max_deposit_period":86400000000000}`)

	params.InitGenesisParameter(&DepositProcedureParameter, ctx, p2)
	require.Equal(t, p1, DepositProcedureParameter.Value)
}

func TestRegisterParamMapping(t *testing.T) {
	skey := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	ctx := defaultContext(skey, tkeyParams)
	cdc := codec.New()

	paramKeeper := params.NewKeeper(
		cdc,
		skey, tkeyParams,
	)

	p1 := NewDepositProcedure()
	p2 := NewDepositProcedure()
    p2.MaxDepositPeriod = time.Duration(THREE_DAYS) * time.Second
	subspace := paramKeeper.Subspace("Gov").WithTypeTable(
		params.NewTypeTable(
			DepositProcedureParameter.GetStoreKey(), DepositProcedure{},
			VotingProcedureParameter.GetStoreKey(), VotingProcedure{},
			TallyingProcedureParameter.GetStoreKey(), TallyingProcedure{},
		))
	params.SetParamReadWriter(subspace, &DepositProcedureParameter, &DepositProcedureParameter)
	params.RegisterGovParamMapping(&DepositProcedureParameter)
	params.InitGenesisParameter(&DepositProcedureParameter, ctx, nil)

	require.Equal(t, params.ParamMapping["Gov/"+string(DepositProcedureParameter.GetStoreKey())].ToJson(""),`{"critical_min_deposit":[{"denom":"iris-atto","amount":"4000000000000000000000"}],"important_min_deposit":[{"denom":"iris-atto","amount":"2000000000000000000000"}],"normal_min_deposit":[{"denom":"iris-atto","amount":"1000000000000000000000"}],"max_deposit_period":86400000000000}`)
	require.Equal(t, p1, DepositProcedureParameter.Value)

	params.ParamMapping["Gov/"+string(DepositProcedureParameter.GetStoreKey())].Update(ctx, `{"critical_min_deposit":[{"denom":"iris-atto","amount":"4000000000000000000000"}],"important_min_deposit":[{"denom":"iris-atto","amount":"2000000000000000000000"}],"normal_min_deposit":[{"denom":"iris-atto","amount":"1000000000000000000000"}],"max_deposit_period":259200000000000}`)
	DepositProcedureParameter.LoadValue(ctx)
	require.Equal(t, p2, DepositProcedureParameter.Value)
}

func TestDepositProcedureParam(t *testing.T) {
	t.SkipNow()
	skey := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	ctx := defaultContext(skey, tkeyParams)
	cdc := codec.New()

	paramKeeper := params.NewKeeper(
		cdc,
		skey, tkeyParams,
	)

	p1 := NewDepositProcedure()

	p2 := NewDepositProcedure()

	subspace := paramKeeper.Subspace("Gov").WithTypeTable(
		params.NewTypeTable(
			DepositProcedureParameter.GetStoreKey(), DepositProcedure{},
			VotingProcedureParameter.GetStoreKey(), VotingProcedure{},
			TallyingProcedureParameter.GetStoreKey(), TallyingProcedure{},
		))

	DepositProcedureParameter.SetReadWriter(subspace)
	find := DepositProcedureParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	DepositProcedureParameter.InitGenesis(nil)
	require.Equal(t, p1, DepositProcedureParameter.Value)

	require.Equal(t, DepositProcedureParameter.ToJson(""), "{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"10000000000000000000\"}],\"max_deposit_period\":172800000000000}")

	DepositProcedureParameter.Update(ctx, "{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"200000000000000000000\"}],\"max_deposit_period\":172800000000000}")

	require.NotEqual(t, p1, DepositProcedureParameter.Value)
	require.Equal(t, p2, DepositProcedureParameter.Value)

	result := DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"atom\",\"amount\":\"200000000000000000000\"}],\"max_deposit_period\":172800000000000}")
	require.Error(t, result)

	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"2000000000000000000\"}],\"max_deposit_period\":172800000000000}")
	require.Error(t, result)

	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"20000000000000000000000000\"}],\"max_deposit_period\":172800000000000}")
	require.Error(t, result)

	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"200000000000000000\"}],\"max_deposit_period\":172800000000000}")
	require.Error(t, result)

	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-att\",\"amount\":\"2000000000000000000\"}],\"max_deposit_period\":172800000000000}")
	require.Error(t, result)

	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"20000000000000000000\"}],\"max_deposit_period\":172800000000000}")
	require.NoError(t, result)

	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"2000000000000000000\"}],\"max_deposit_period\":1}")
	require.Error(t, result)

	result = DepositProcedureParameter.Valid("{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"2000000000000000000\"}],\"max_deposit_period\":172800000000000}")
	require.Error(t, result)

	DepositProcedureParameter.InitGenesis(p2)
	require.Equal(t, p2, DepositProcedureParameter.Value)
	DepositProcedureParameter.InitGenesis(p1)
	require.Equal(t, p1, DepositProcedureParameter.Value)

	DepositProcedureParameter.LoadValue(ctx)
	require.Equal(t, p2, DepositProcedureParameter.Value)

}

func TestVotingProcedureParam(t *testing.T) {
	t.SkipNow()
	skey := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	ctx := defaultContext(skey, tkeyParams)
	cdc := codec.New()

	paramKeeper := params.NewKeeper(
		cdc,
		skey, tkeyParams,
	)

	p1 := NewVotingProcedure()

	p2 := NewVotingProcedure()

	subspace := paramKeeper.Subspace("Gov").WithTypeTable(
		params.NewTypeTable(
			DepositProcedureParameter.GetStoreKey(), DepositProcedure{},
			VotingProcedureParameter.GetStoreKey(), VotingProcedure{},
			TallyingProcedureParameter.GetStoreKey(), TallyingProcedure{},
		))

	VotingProcedureParameter.SetReadWriter(subspace)
	find := VotingProcedureParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	VotingProcedureParameter.InitGenesis(nil)
	require.Equal(t, p1, VotingProcedureParameter.Value)

	require.Equal(t, VotingProcedureParameter.ToJson(""), "{\"voting_period\":172800000000000}")

	VotingProcedureParameter.Update(ctx, "{\"voting_period\":192800000000000}")

	require.NotEqual(t, p1, VotingProcedureParameter.Value)
	require.Equal(t, p2, VotingProcedureParameter.Value)

	result := VotingProcedureParameter.Valid("{\"voting_period\":400000}")
	require.Error(t, result)

	VotingProcedureParameter.InitGenesis(p2)
	require.Equal(t, p2, VotingProcedureParameter.Value)
	VotingProcedureParameter.InitGenesis(p1)
	require.Equal(t, p1, VotingProcedureParameter.Value)

	VotingProcedureParameter.LoadValue(ctx)
	require.Equal(t, p2, VotingProcedureParameter.Value)

}

func TestTallyingProcedureParam(t *testing.T) {
	t.SkipNow()
	skey := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	ctx := defaultContext(skey, tkeyParams)
	cdc := codec.New()

	paramKeeper := params.NewKeeper(
		cdc,
		skey, tkeyParams,
	)

	p1 := NewTallyingProcedure()

	p2 := NewTallyingProcedure()

	subspace := paramKeeper.Subspace("Gov").WithTypeTable(
		params.NewTypeTable(
			DepositProcedureParameter.GetStoreKey(), DepositProcedure{},
			VotingProcedureParameter.GetStoreKey(), VotingProcedure{},
			TallyingProcedureParameter.GetStoreKey(), TallyingProcedure{},
		))

	TallyingProcedureParameter.SetReadWriter(subspace)
	find := TallyingProcedureParameter.LoadValue(ctx)
	require.Equal(t, find, false)

	TallyingProcedureParameter.InitGenesis(nil)
	require.Equal(t, p1, TallyingProcedureParameter.Value)
	require.Equal(t, "{\"threshold\":\"0.5000000000\",\"veto\":\"0.3340000000\",\"participation\":\"0.6670000000\"}", TallyingProcedureParameter.ToJson(""))

	TallyingProcedureParameter.Update(ctx, "{\"threshold\":\"0.5\",\"veto\":\"0.3340000000\",\"participation\":\"0.0200000000\"}")

	require.NotEqual(t, p1, TallyingProcedureParameter.Value)
	require.Equal(t, p2, TallyingProcedureParameter.Value)

	result := TallyingProcedureParameter.Valid("{\"threshold\":\"1/1\",\"veto\":\"1/3\",\"participation\":\"1/100\"}")
	require.Error(t, result)

	result = TallyingProcedureParameter.Valid("{\"threshold\":\"abcd\",\"veto\":\"1/3\",\"participation\":\"1/100\"}")
	require.Error(t, result)

	TallyingProcedureParameter.InitGenesis(p2)
	require.Equal(t, p2, TallyingProcedureParameter.Value)
	TallyingProcedureParameter.InitGenesis(p1)
	require.Equal(t, p1, TallyingProcedureParameter.Value)

	TallyingProcedureParameter.LoadValue(ctx)
	require.Equal(t, p2, TallyingProcedureParameter.Value)

}
