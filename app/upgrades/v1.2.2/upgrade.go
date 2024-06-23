// Package v1_2_2 is contains chain upgrade of the corresponding version.
package v1_2_2 //nolint:revive,stylecheck // app version

import (
	"strings"

	"github.com/onomyprotocol/onex/app/upgrades"
	"github.com/pendulum-labs/market/x/market/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func addUid(s []uint64, r uint64) ([]uint64, bool) {
	for _, v := range s {
		if v == r {
			return s, false
		}
	}

	return append(s, r), true
}

func removeUid(s []uint64, r uint64) ([]uint64, bool) {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...), true
		}
	}
	return s, false
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *upgrades.UpgradeKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Starting module migrations...")

		// Wipe corrupted leaderboards
		pools := keepers.MarketKeeper.GetAllPool(ctx)
		for _, pool := range pools {
			pool.Leaders = []*types.Leader{}
			keepers.MarketKeeper.SetPool(ctx, pool)
		}

		var dropOwner types.Drops
		var dropper types.Drop

		// For each drop in database
		drops := keepers.MarketKeeper.GetAllDrop(ctx)
		for _, drop := range drops {
			// Get DropsOwner associated with this drop
			dropOwner, _ = keepers.MarketKeeper.GetDropsOwnerPair(ctx, drop.Owner, drop.Pair)
			// Reset dropOwner.Sum to Zero
			dropOwner.Sum = sdk.ZeroInt()
			// Add (active) or remove (inactive) uid
			if drop.Active {
				dropOwner.Uids, _ = addUid(dropOwner.Uids, drop.Uid)
			} else {
				dropOwner.Uids, _ = removeUid(dropOwner.Uids, drop.Uid)
			}
			// Recalculate
			for _, uid := range dropOwner.Uids {
				dropper, _ = keepers.MarketKeeper.GetDrop(ctx, uid)
				if dropper.Active {
					dropOwner.Sum = dropOwner.Sum.Add(drop.Drops)
				}
			}
			// Get pool associated with drop
			pool, _ := keepers.MarketKeeper.GetPool(ctx, drop.Pair)
			// Update leaders in pool based on recalculated sum
			pool = updateLeaders(ctx, pool, drop.Owner, dropOwner.Sum, keepers)
			// Set pool and drops
			keepers.MarketKeeper.SetPool(ctx, pool)
			keepers.MarketKeeper.SetDrops(ctx, dropOwner, drop.Owner, drop.Pair)
		}

		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		ctx.Logger().Info("Upgrade complete")
		return vm, err
	}
}

func updateLeaders(ctx sdk.Context, pool types.Pool, dropCreator string, dropCreatorSum sdk.Int, keepers *upgrades.UpgradeKeepers) types.Pool {
	maxLeaders := len(strings.Split(keepers.MarketKeeper.EarnRates(ctx), ","))

	for i := 0; i < len(pool.Leaders); i++ {
		if pool.Leaders[i].Address == dropCreator {
			pool.Leaders = pool.Leaders[:i+copy(pool.Leaders[i:], pool.Leaders[i+1:])]
		}
	}

	if dropCreatorSum.Equal(sdk.ZeroInt()) {
		return pool
	}

	if len(pool.Leaders) == 0 {
		pool.Leaders = append(pool.Leaders, &types.Leader{
			Address: dropCreator,
			Drops:   dropCreatorSum,
		})
	} else {
		for i := 0; i < len(pool.Leaders); i++ {
			if dropCreatorSum.GT(pool.Leaders[i].Drops) {
				if len(pool.Leaders) < maxLeaders {
					pool.Leaders = append(pool.Leaders, pool.Leaders[len(pool.Leaders)-1])
				}
				copy(pool.Leaders[i+1:], pool.Leaders[i:])
				pool.Leaders[i] = &types.Leader{
					Address: dropCreator,
					Drops:   dropCreatorSum,
				}
				break
			} else {
				if (i == len(pool.Leaders)-1) && len(pool.Leaders) < maxLeaders {
					pool.Leaders = append(pool.Leaders, &types.Leader{
						Address: dropCreator,
						Drops:   dropCreatorSum,
					})
					break
				}
			}
		}
	}
	return pool
}
