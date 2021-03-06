package keeper

import (
	"context"
	"fmt"
	"strings"

	"github.com/batphonghan/checkers/x/checkers/rules"
	"github.com/batphonghan/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) ForfeitExpiredGames(goCtx context.Context) {
	// TODO
	ctx := sdk.UnwrapSDKContext(goCtx)

	opponents := map[string]string{
		rules.BLACK_PLAYER.Color: rules.RED_PLAYER.Color,
		rules.RED_PLAYER.Color:   rules.BLACK_PLAYER.Color,
	}

	nextGame, found := k.GetNextGame(ctx)
	if !found {
		panic("NextGame not found")
	}
	storedGameId := nextGame.FifoHead
	var storedGame types.StoredGame

	for {
		if strings.Compare(storedGameId, types.NoFifoIdKey) == 0 {
			break
		}

		storedGame, found = k.GetStoredGame(ctx, storedGameId)
		if !found {
			panic("Fifo head game not found " + nextGame.FifoHead)
		}
		deadline, err := storedGame.GetDeadlineAsTime()
		if err != nil {
			panic(err)
		}

		if deadline.Before(ctx.BlockTime()) {
			// TODO
		} else {
			// All other games come after anyway
			break
		}

		k.RemoveFromFifo(ctx, &storedGame, &nextGame)

		if storedGame.MoveCount == 0 {
			storedGame.Winner = rules.NO_PLAYER.Color
			// No point in keeping a game that was never played
			k.RemoveStoredGame(ctx, storedGameId)
		} else {
			storedGame.Winner, found = opponents[storedGame.Turn]
			if !found {
				panic(fmt.Sprintf(types.ErrCannotFindWinnerByColor.Error(), storedGame.Turn))
			}
			k.SetStoredGame(ctx, storedGame)
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
				sdk.NewAttribute(sdk.AttributeKeyAction, types.ForfeitGameEventKey),
				sdk.NewAttribute(types.ForfeitGameEventIdValue, storedGameId),
				sdk.NewAttribute(types.ForfeitGameEventWinner, storedGame.Winner),
			),
		)

		storedGameId = nextGame.FifoHead
	}

	k.SetNextGame(ctx, nextGame)

}
