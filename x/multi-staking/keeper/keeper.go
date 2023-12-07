package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

type Keeper struct {
	storeKey         storetypes.StoreKey
	memKey           storetypes.StoreKey
	cdc              codec.BinaryCodec
	stakingKeeper    types.StakingKeeper
	stakingMsgServer stakingtypes.MsgServer
	distrMsgServer   distrtypes.MsgServer
	govMsgServer     govtypes.MsgServer
	bankKeeper       types.BankKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	stakingKeeper stakingkeeper.Keeper,
	distrKeeper distrkeeper.Keeper,
	govKeeper govkeeper.Keeper,
	bankKeeper types.BankKeeper,
	key storetypes.StoreKey,
	memKey storetypes.StoreKey,
) *Keeper {
	return &Keeper{
		cdc:              cdc,
		storeKey:         key,
		memKey:           memKey,
		stakingKeeper:    stakingKeeper,
		stakingMsgServer: stakingkeeper.NewMsgServerImpl(stakingKeeper),
		distrMsgServer:   distrkeeper.NewMsgServerImpl(distrKeeper),
		govMsgServer:     govkeeper.NewMsgServerImpl(govKeeper),
		bankKeeper:       bankKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k Keeper) BeforeUnbondedHandle(ctx sdk.Context) types.UnbondedMultiStakings {
	var unbondedStakingLists types.UnbondedMultiStakings
	matureUnbonds := k.stakingKeeper.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	for _, dvPair := range matureUnbonds {
		valAddr, err := sdk.ValAddressFromBech32(dvPair.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		delAddr := sdk.MustAccAddressFromBech32(dvPair.DelegatorAddress)

		ubdListByDelVal, err := k.GetUnbondedAmount(ctx, delAddr, valAddr)
		if err != nil {
			continue
		}

		unbondedStakingLists = append(unbondedStakingLists, ubdListByDelVal...)
	}

	return unbondedStakingLists
}

func (k Keeper) GetUnbondedAmount(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (types.UnbondedMultiStakings, error) {
	var unbondedStakingLists []types.UnbonedMultiStakingRecord

	ubd, found := k.stakingKeeper.GetUnbondingDelegation(ctx, delAddr, valAddr)
	if !found {
		return nil, stakingtypes.ErrNoUnbondingDelegation
	}

	ctxTime := ctx.BlockHeader().Time
	heightSumMap := make(map[int64]types.UnbonedMultiStakingRecord)

	// loop through all the entries and complete unbonding mature entries
	for i := 0; i < len(ubd.Entries); i++ {
		entry := ubd.Entries[i]
		if entry.IsMature(ctxTime) && !entry.Balance.IsZero()  {
			ubdr := types.NewUnbonedMultiStakingRecord(
				delAddr,
				valAddr,
				entry.CreationHeight,
				entry.CompletionTime,
				sdk.OneDec(),
				entry.Balance,
			)
			existingItem, exists := heightSumMap[entry.CreationHeight]
			if exists {
				existingItem.Balance = entry.Balance.Add(existingItem.Balance)
				heightSumMap[entry.CreationHeight] = existingItem
			} else {
				heightSumMap[entry.CreationHeight] = ubdr
			}

		}
	}

	for _, ubdr := range heightSumMap {
		unbondedStakingLists = append(unbondedStakingLists, ubdr)
	}
	return unbondedStakingLists, nil
}