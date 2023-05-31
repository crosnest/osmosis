package cosmwasmpool

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/v15/x/cosmwasmpool/types"
)

var (
	emptyPlaceholder = []byte{}
)

// AddToWhitelist adds the code id to the whitelist.
func (k Keeper) AddToWhitelist(ctx sdk.Context, codeId uint64) {
	ctx.KVStore(k.storeKey).Set(types.FormatCodeIdWhitelistPrefix(codeId), emptyPlaceholder)
}

// removeFromWhitelist removes the code id from the whitelist.
// nolint: unused
func (k Keeper) removeFromWhitelist(ctx sdk.Context, codeId uint64) {
	ctx.KVStore(k.storeKey).Delete(types.FormatCodeIdWhitelistPrefix(codeId))
}

// isWhitelisted returns true if the code id is in the whitelist.
func (k Keeper) isWhitelisted(ctx sdk.Context, codeId uint64) bool {
	return ctx.KVStore(k.storeKey).Has(types.FormatCodeIdWhitelistPrefix(codeId))
}
