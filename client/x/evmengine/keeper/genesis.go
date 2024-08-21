package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/storyprotocol/iliad/client/x/evmengine/types"
	"github.com/storyprotocol/iliad/lib/errors"
)

func (k *Keeper) InitGenesis(ctx context.Context, gs *types.GenesisState) error {
	if err := k.ValidateGenesis(gs); err != nil {
		return err
	}
	if err := k.SetParams(ctx, gs.Params); err != nil {
		return err
	}

	if err := k.InsertGenesisHead(ctx, gs.Params.ExecutionBlockHash); err != nil {
		panic(errors.Wrap(err, "insert genesis head"))
	}

	return nil
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k *Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		Params: params,
	}
}

//nolint:revive // TODO: validate genesis
func (k *Keeper) ValidateGenesis(gs *types.GenesisState) error {
	return types.ValidateExecutionBlockHash(gs.Params.ExecutionBlockHash)
}
