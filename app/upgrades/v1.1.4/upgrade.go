// Package v1_1_4 is contains chain upgrade of the corresponding version.
package v1_1_4 //nolint:revive,stylecheck // app version

import (
	"github.com/onomyprotocol/onex/app/upgrades"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	markettypes "github.com/pendulum-labs/market/x/market/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *upgrades.UpgradeKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Starting module migrations...")

		// Deactivate all drops
		drops := keepers.MarketKeeper.GetAllDrop(ctx)
		for _, drop := range drops {
			drop.Active = false
			keepers.MarketKeeper.SetDrop(ctx, drop)
		}

		// Set pool drops to zero and wipe leaders
		pools := keepers.MarketKeeper.GetAllPool(ctx)
		for _, pool := range pools {
			pool.Drops = sdk.ZeroInt()
			pool.Leaders = []*markettypes.Leader{}

			keepers.MarketKeeper.SetPool(ctx, pool)
		}

		// Set member balances to zero
		members := keepers.MarketKeeper.GetAllMember(ctx)
		for _, member := range members {
			member.Balance = sdk.ZeroInt()
			keepers.MarketKeeper.SetMember(ctx, member)
		}

		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		ctx.Logger().Info("Upgrade complete")
		return vm, err
	}
}
