// Package v1_2_0 is contains chain upgrade of the corresponding version.
package v1_2_0 //nolint:revive,stylecheck // app version

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

		// Deactivate all drops and remove owners
		drops := keepers.MarketKeeper.GetAllDrop(ctx)
		for _, drop := range drops {
			drop.Active = false
			keepers.MarketKeeper.SetDrop(ctx, drop)
			keepers.MarketKeeper.RemoveDropOwner(ctx, drop)
		}

		// Set pool drops to zero and wipe leaders
		pools := keepers.MarketKeeper.GetAllPool(ctx)
		for _, pool := range pools {
			pool.Drops = sdk.ZeroInt()
			pool.History = 0
			pool.Leaders = []*markettypes.Leader{}
			pool.Volume1.Amount = sdk.ZeroInt()
			pool.Volume2.Amount = sdk.ZeroInt()
			keepers.MarketKeeper.SetPool(ctx, pool)
		}

		// Set member balances to zero
		members := keepers.MarketKeeper.GetAllMember(ctx)
		for _, member := range members {
			member.Balance = sdk.ZeroInt()
			member.Limit = 0
			member.Stop = 0
			member.Previous = sdk.ZeroInt()
			keepers.MarketKeeper.SetMember(ctx, member)
		}

		// Set order status to balances to zero
		orders := keepers.MarketKeeper.GetAllOrder(ctx)
		for _, order := range orders {
			if order.Status == "active" {
				order.Status = "canceled"
			}
			keepers.MarketKeeper.SetOrder(ctx, order)
		}

		marketAccount := keepers.AccountKeeper.GetModuleAccount(ctx, markettypes.ModuleName)

		marketCoins := keepers.BankKeeper.GetAllBalances(ctx, marketAccount.GetAddress())

		reclaimer, _ := sdk.AccAddressFromBech32("onomy1yc0lg97cy5e80jyajtkz0zke2rr4734anugf9g")

		keepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, markettypes.ModuleName, reclaimer, marketCoins)

		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		ctx.Logger().Info("Upgrade complete")
		return vm, err
	}
}
