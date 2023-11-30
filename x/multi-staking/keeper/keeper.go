package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

type Keeper struct {
	storeKey       storetypes.StoreKey
	memKey         storetypes.StoreKey
	cdc            codec.BinaryCodec
	stakingKeeper  stakingkeeper.Keeper
	distrMsgServer distrtypes.MsgServer
	bankKeeper     types.BankKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	stakingKeeper stakingkeeper.Keeper,
	distrKeeper distrkeeper.Keeper,
	bankKeeper types.BankKeeper,
	key storetypes.StoreKey,
	memKey storetypes.StoreKey,
) *Keeper {
	return &Keeper{
		cdc:            cdc,
		storeKey:       key,
		memKey:         memKey,
		stakingKeeper:  stakingKeeper,
		distrMsgServer: distrkeeper.NewMsgServerImpl(distrKeeper),
		bankKeeper:     bankKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k Keeper) BeforeUnbondedHandle(ctx sdk.Context) types.UnbondedMultiStakings {
	var unbondedStakingLists []types.UnbonedMultiStaking
	matureUnbonds := k.stakingKeeper.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	for _, dvPair := range matureUnbonds {
		valAddr, err := sdk.ValAddressFromBech32(dvPair.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		delAddr := sdk.MustAccAddressFromBech32(dvPair.DelegatorAddress)

		balances, err := k.GetUnbondedAmount(ctx, delAddr, valAddr)
		if err != nil {
			continue
		}

		unbondedStaking := types.UnbonedMultiStaking{
			DelAddr: delAddr.String(),
			ValAddr: valAddr.String(),
			Amount:  balances,
		}

		unbondedStakingLists = append(unbondedStakingLists, unbondedStaking)
	}

	unbondedStakings := types.UnbondedMultiStakings{
		UnbondedMultiStakings: unbondedStakingLists,
	}

	return unbondedStakings
}

func (k Keeper) GetUnbondedAmount(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error) {
	ubd, found := k.stakingKeeper.GetUnbondingDelegation(ctx, delAddr, valAddr)
	if !found {
		return nil, stakingtypes.ErrNoUnbondingDelegation
	}

	bondDenom := k.stakingKeeper.GetParams(ctx).BondDenom
	balances := sdk.NewCoins()
	ctxTime := ctx.BlockHeader().Time

	// loop through all the entries and complete unbonding mature entries
	for i := 0; i < len(ubd.Entries); i++ {
		entry := ubd.Entries[i]
		if entry.IsMature(ctxTime) {
			ubd.RemoveEntry(int64(i))
			i--
			// track undelegation only when remaining or truncated shares are non-zero
			if !entry.Balance.IsZero() {
				amt := sdk.NewCoin(bondDenom, entry.Balance)
				balances = balances.Add(amt)
			}
		}
	}

	return balances, nil
}
