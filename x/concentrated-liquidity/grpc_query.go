package concentrated_liquidity

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/osmosis-labs/osmosis/osmoutils"
	"github.com/osmosis-labs/osmosis/v15/x/concentrated-liquidity/model"
	"github.com/osmosis-labs/osmosis/v15/x/concentrated-liquidity/types"
)

const (
	tickLiquidityInBatchQueryLimit = 10000
	liquidityDepthRangeQueryLimit  = 10000
)

var (
	liquidityDepthRangeQueryLimitInt = sdk.NewInt(liquidityDepthRangeQueryLimit)
)

var _ types.QueryServer = Querier{}

// Querier defines a wrapper around the x/concentrated-liquidity keeper providing gRPC method
// handlers.
type Querier struct {
	Keeper
}

func NewQuerier(k Keeper) Querier {
	return Querier{Keeper: k}
}

// Pool checks if a pool exists and returns the desired pool.
func (q Querier) Pool(
	ctx context.Context,
	req *types.QueryPoolRequest,
) (*types.QueryPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	pool, err := q.Keeper.GetPool(sdkCtx, req.PoolId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	any, err := codectypes.NewAnyWithValue(pool)
	if err != nil {
		return nil, err
	}

	return &types.QueryPoolResponse{Pool: any}, nil
}

// UserPositions returns positions of a specified address
func (q Querier) UserPositions(ctx context.Context, req *types.QueryUserPositionsRequest) (*types.QueryUserPositionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkAddr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	userPositions, err := q.Keeper.GetUserPositions(sdkCtx, sdkAddr)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryUserPositionsResponse{
		Positions: userPositions,
	}, nil
}

// Pools returns all concentrated pools in existence.
func (q Querier) Pools(
	ctx context.Context,
	req *types.QueryPoolsRequest,
) (*types.QueryPoolsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := sdkCtx.KVStore(q.Keeper.storeKey)
	poolStore := prefix.NewStore(store, types.PoolPrefix)

	var anys []*codectypes.Any
	pageRes, err := query.Paginate(poolStore, req.Pagination, func(key, _ []byte) error {
		pool := model.Pool{}
		// Get the next pool from the poolStore and pass it to the pool variable
		_, err := osmoutils.Get(poolStore, key, &pool)
		if err != nil {
			return err
		}

		// Retrieve the poolInterface from the respective pool
		poolI, err := q.Keeper.GetPool(sdkCtx, pool.GetId())
		if err != nil {
			return err
		}

		any, err := codectypes.NewAnyWithValue(poolI)
		if err != nil {
			return err
		}

		anys = append(anys, any)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPoolsResponse{
		Pools:      anys,
		Pagination: pageRes,
	}, nil
}

// Params returns module params
func (q Querier) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.QueryParamsResponse{Params: q.Keeper.GetParams(ctx)}, nil
}

// TickLiquidityInBatches returns array of liquidity depths in the given range of lower tick and upper tick.
// Note that the space between the ticks in the returned array would always be guaranteed spacing greater than given batch unit.
func (q Querier) TickLiquidityInBatches(goCtx context.Context, req *types.QueryTickLiquidityInBatchesRequest) (*types.QueryTickLiquidityInBatchesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	deltaTick := req.UpperTick - req.LowerTick
	if deltaTick < 0 {
		deltaTick = -deltaTick
	}
	requestedAmtofTicks := deltaTick / int64(req.BatchUnit)
	if tickLiquidityInBatchQueryLimit > requestedAmtofTicks {
		return nil, types.QueryRangeUnsupportedError{RequestedRange: requestedAmtofTicks, MaxRange: tickLiquidityInBatchQueryLimit}
	}

	liquidityDepths, err := q.Keeper.GetTickLiquidityForRangeInBatches(
		ctx,
		req.PoolId,
		req.LowerTick,
		req.UpperTick,
		req.BatchUnit,
	)
	if err != nil {
		return nil, err
	}

	return &types.QueryTickLiquidityInBatchesResponse{LiquidityDepths: liquidityDepths}, nil
}

// LiquidityDepthsForRange returns liquidity depths for the given range.
func (q Querier) LiquidityDepthsForRange(goCtx context.Context, req *types.QueryLiquidityDepthsForRangeRequest) (*types.QueryLiquidityDepthsForRangeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.LowerTick > req.UpperTick {
		return nil, types.InvalidLowerUpperTickError{LowerTick: req.LowerTick, UpperTick: req.UpperTick}
	}

	requestedRange := req.UpperTick - req.LowerTick
	// use constant pre-defined to limit range and check if reuested range does not exceed max range
	if requestedRange > liquidityDepthRangeQueryLimit {
		return nil, types.QueryRangeUnsupportedError{RequestedRange: requestedRange, MaxRange: liquidityDepthRangeQueryLimit}
	}

	liquidityDepths, err := q.Keeper.GetPerTickLiquidityDepthFromRange(ctx, req.PoolId, req.LowerTick, req.UpperTick)
	if err != nil {
		return nil, err
	}

	return &types.QueryLiquidityDepthsForRangeResponse{
		LiquidityDepths: liquidityDepths,
	}, nil
}
