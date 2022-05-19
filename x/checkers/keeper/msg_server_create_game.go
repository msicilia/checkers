package keeper

import (
	"context"
	"strconv"

	rules "github.com/alice/checkers/x/checkers/rules"
	"github.com/alice/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateGame(goCtx context.Context, msg *types.MsgCreateGame) (*types.MsgCreateGameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx
	// Get the existing GameID with the GetNextGame function scaffolded
	nextGame, found := k.Keeper.GetNextGame(ctx)
	if !found {
		panic("NextGame not found")
	}
	newIndex := strconv.FormatUint(nextGame.IdValue, 10) // is incremented below.

	// Create a new game and check it is ok (Validate checks addresses).
	newGame := rules.New()
	storedGame := types.StoredGame{
		Creator:   msg.Creator,
		Index:     newIndex,
		Game:      newGame.String(),
		Turn:      rules.PieceStrings[newGame.Turn],
		Red:       msg.Red,
		Black:     msg.Black,
		MoveCount: 0,
	}

	err := storedGame.Validate()
	if err != nil {
		return nil, err
	}
	// Now set the new values of the two objects.
	k.Keeper.SetStoredGame(ctx, storedGame)
	nextGame.IdValue++
	k.Keeper.SetNextGame(ctx, nextGame)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.StoredGameEventKey),
			sdk.NewAttribute(types.StoredGameEventCreator, msg.Creator),
			sdk.NewAttribute(types.StoredGameEventIndex, newIndex),
			sdk.NewAttribute(types.StoredGameEventRed, msg.Red),
			sdk.NewAttribute(types.StoredGameEventBlack, msg.Black),
		),
	)

	// Return the newly created ID
	return &types.MsgCreateGameResponse{
		IdValue: newIndex,
	}, nil
}
