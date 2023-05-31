package cosmwasmpool

import sdk "github.com/cosmos/cosmos-sdk/types"

func (k Keeper) RemoveFromWhitelist(ctx sdk.Context, codeId uint64) {
	k.removeFromWhitelist(ctx, codeId)
}

func (k Keeper) IsWhitelisted(ctx sdk.Context, codeId uint64) bool {
	return k.isWhitelisted(ctx, codeId)
}
