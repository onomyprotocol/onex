// Package v1_2_2 is contains chain upgrade of the corresponding version.
package v1_2_2 //nolint:revive,stylecheck // app version

import (
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

		var dropOwner types.Drops
		var dropp types.Drop

		// For each drop in database
		drops := keepers.MarketKeeper.GetAllDrop(ctx)
		for _, drop := range drops {
			// Get DropsOwner associated with this drop
			dropOwner, _ = keepers.MarketKeeper.GetDropsOwnerPair(ctx, drop.Owner, drop.Pair)
			// Reset dropOwner.Sum to Zero
			dropOwner.Sum = sdk.ZeroInt()
			if drop.Active {
				dropOwner.Uids, _ = addUid(dropOwner.Uids, drop.Uid)
			} else {
				dropOwner.Uids, _ = removeUid(dropOwner.Uids, drop.Uid)
			}
			// Recaculate
			for _, uid := range dropOwner.Uids {
				dropp, _ = keepers.MarketKeeper.GetDrop(ctx, uid)
				if dropp.Active {
					dropOwner.Sum = dropOwner.Sum.Add(drop.Drops)
				}
			}
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
