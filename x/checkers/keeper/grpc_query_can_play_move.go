package keeper

import (
	"context"
	"strings"

	"github.com/batphonghan/checkers/x/checkers/rules"
	"github.com/batphonghan/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) CanPlayMove(goCtx context.Context, req *types.QueryCanPlayMoveRequest) (*types.QueryCanPlayMoveResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Game not found
	storedGame, found := k.GetStoredGame(ctx, req.IdValue)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrGameNotFound, types.ErrGameNotFound.Error(), req.IdValue)
	}

	if storedGame.Winner != rules.NO_PLAYER.Color {
		return &types.QueryCanPlayMoveResponse{
			Possible: false,
			Reason:   types.ErrGameFinished.Error(),
		}, nil
	}
	// game already been won
	if storedGame.Winner != rules.NO_PLAYER.Color {
		return &types.QueryCanPlayMoveResponse{
			Possible: false,
			Reason:   types.ErrGameFinished.Error(),
		}, nil
	}

	// Is valid player
	var player rules.Player
	if strings.Compare(rules.RED_PLAYER.Color, req.Player) == 0 {
		player = rules.RED_PLAYER
	} else if strings.Compare(rules.BLACK_PLAYER.Color, req.Player) == 0 {
		player = rules.BLACK_PLAYER
	} else {
		return &types.QueryCanPlayMoveResponse{
			Possible: false,
			Reason:   types.ErrCreatorNotPlayer.Error(),
		}, nil
	}

	// Is player turn
	game, err := storedGame.ParseGame()
	if err != nil {
		return nil, err
	}
	if !game.TurnIs(player) {
		return &types.QueryCanPlayMoveResponse{
			Possible: false,
			Reason:   types.ErrNotPlayerTurn.Error(),
		}, nil
	}

	// Attemp to move and report back
	game, err = storedGame.ParseGame()
	if err != nil {
		return nil, err
	}
	if !game.TurnIs(player) {
		return &types.QueryCanPlayMoveResponse{
			Possible: false,
			Reason:   types.ErrNotPlayerTurn.Error(),
		}, nil
	}

	return &types.QueryCanPlayMoveResponse{
		Possible: true,
		Reason:   "ok",
	}, nil
}
