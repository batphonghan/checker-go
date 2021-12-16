package keeper

import (
	"context"
	"strconv"

	"github.com/batphonghan/checkers/x/checkers/rules"
	"github.com/batphonghan/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateGame(goCtx context.Context, msg *types.MsgCreateGame) (*types.MsgCreateGameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	nextGame, found := k.Keeper.GetNextGame(ctx)
	if !found {
		panic("NextGame not found")
	}
	newIndex := strconv.FormatUint(nextGame.IdValue, 10)

	storedGame := types.StoredGame{
		Creator: msg.Creator,
		Index:   newIndex,
		Game:    rules.New().String(),
		Red:     msg.Red,
		Black:   msg.Black,
	}

	// TODO
	// err := storedGame.Validate()
	// if err != nil {
	// 	return nil, err
	// }

	k.Keeper.SetStoredGame(ctx, storedGame)

	nextGame.IdValue++
	k.Keeper.SetNextGame(ctx, nextGame)

	return &types.MsgCreateGameResponse{
		IdValue: newIndex,
	}, nil
}
