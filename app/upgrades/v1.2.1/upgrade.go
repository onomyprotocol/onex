// Package v1_2_1 is contains chain upgrade of the corresponding version.
package v1_2_1 //nolint:revive,stylecheck // app version

import (
	"strings"

	"github.com/onomyprotocol/onex/app/upgrades"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *upgrades.UpgradeKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Starting module migrations...")

		// Deactivate all drops and remove owners
		drops := keepers.MarketKeeper.GetAllDrop(ctx)
		for _, drop := range drops {
			if drop.Active {
				pair := strings.Split(drop.Pair, ",")

				denom1 := pair[0]
				denom2 := pair[1]

				pool, found := keepers.MarketKeeper.GetPool(ctx, drop.Pair)
				if !found {
					continue
				}

				member1, found := keepers.MarketKeeper.GetMember(ctx, denom2, denom1)
				if !found {
					continue
				}

				member2, found := keepers.MarketKeeper.GetMember(ctx, denom1, denom2)
				if !found {
					continue
				}

				// `total1 = (drop.Drops * member1.Balance) / pool.Drops`
				total1 := (drop.Drops.Mul(member1.Balance)).Quo(pool.Drops)
				total2 := (drop.Drops.Mul(member2.Balance)).Quo(pool.Drops)

				drop.Product = total1.Mul(total2)
				keepers.MarketKeeper.SetDrop(ctx, drop)
			}
		}

		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		ctx.Logger().Info("Upgrade complete")
		return vm, err
	}
}
