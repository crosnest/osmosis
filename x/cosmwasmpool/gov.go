package cosmwasmpool

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	"github.com/osmosis-labs/osmosis/v15/x/cosmwasmpool/types"
)

const (
	poolCountLimit = 20
)

func NewCosmWasmPoolProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.UploadCosmWasmPoolCodeAndWhiteListProposal:
			return k.HandleUploadCosmWasmPoolCodeAndWhiteListProposal(ctx, c)
		default:
			return fmt.Errorf("unrecognized concentrated liquidity proposal content type: %T", c)
		}
	}
}

func (k Keeper) HandleUploadCosmWasmPoolCodeAndWhiteListProposal(ctx sdk.Context, p *types.UploadCosmWasmPoolCodeAndWhiteListProposal) error {
	cosmwasmPoolModuleAddress := k.accountKeeper.GetModuleAddress(types.ModuleName)

	// Only allow the x/cosmwasmpool module to instantiate this contract.
	instantiatePermissions := wasmtypes.AccessConfig{
		Permission: wasmtypes.AccessTypeAnyOfAddresses,
		Addresses:  []string{cosmwasmPoolModuleAddress.String()},
	}

	// Upload the code to the wasmvm.
	codeID, checksum, err := k.contractKeeper.Create(ctx, cosmwasmPoolModuleAddress, p.WASMByteCode, &instantiatePermissions)
	if err != nil {
		return err
	}

	// Add the code id to the whitelist.
	k.AddToWhitelist(ctx, codeID)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.TypeEvtUploadedCosmwasmPoolCode,
		sdk.NewAttribute(types.AttributeKeyCodeID, strconv.FormatUint(codeID, 10)),
		sdk.NewAttribute(types.AttributeKeyChecksum, string(checksum)),
	))

	return nil
}

func (k Keeper) HandleMigrateCosmWasmPoolCodeAndWhiteListProposal(ctx sdk.Context, poolIds []uint64, newCodeId uint64, uploadByteCode []byte, migrateMsg []byte) error {
	cosmwasmPoolModuleAddress := k.accountKeeper.GetModuleAddress(types.ModuleName)

	requestedPoolMigrationCount := len(poolIds)
	if requestedPoolMigrationCount > poolCountLimit {
		return fmt.Errorf("pool count (%d) exceeds limit (%d)", requestedPoolMigrationCount, poolCountLimit)
	}

	// Iterate requested pool ids to make sure that pool with such id exists.
	poolCount := k.poolmanagerKeeper.GetNextPoolId(ctx) - 1
	for _, poolId := range poolIds {
		if poolId > poolCount {
			return fmt.Errorf("pool id (%d) does not exist", poolId)
		}
	}

	// Iterate over pool ids and attempt to migrate each pool.
	for _, poolId := range poolIds {
		cwPool, err := k.GetPoolById(ctx, poolId)
		if err != nil {
			return err
		}

		k.wasmKeeper.GetCodeId()

		_, err = k.contractKeeper.Migrate(ctx, sdk.MustAccAddressFromBech32(cwPool.GetContractAddress()), cosmwasmPoolModuleAddress, newCodeId, migrateMsg)
		if err != nil {
			return err
		}
	}

	// Add the code id to the whitelist.
	k.AddToWhitelist(ctx, codeID)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.TypeEvtUploadedCosmwasmPoolCode,
		sdk.NewAttribute(types.AttributeKeyCodeID, strconv.FormatUint(codeID, 10)),
		sdk.NewAttribute(types.AttributeKeyChecksum, string(checksum)),
	))

	return nil
}
