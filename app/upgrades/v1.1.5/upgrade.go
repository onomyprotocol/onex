// Package v1_1_5 is contains chain upgrade of the corresponding version.
package v1_1_5 //nolint:revive,stylecheck // app version

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

		onexAmount, _ := sdk.NewIntFromString("94784650277373001264452")

		ibc1Amount, _ := sdk.NewIntFromString("998915079")

		ibc2Amount, _ := sdk.NewIntFromString("28966246179579930912901")

		ibc3Amount, _ := sdk.NewIntFromString("2952795016")

		onexCoin := sdk.NewCoin("aonex", onexAmount)
		ibc1 := sdk.NewCoin("ibc/30EDC220372A2C3D0FC1D987E19062E35375DECD1001A5EFA44EB92FF59D1867", ibc1Amount)
		ibc2 := sdk.NewCoin("ibc/5BDD8875CC2AF7BC842BE44236ACD576EA4F53C36347F74903B852060D6BF29A", ibc2Amount)
		ibc3 := sdk.NewCoin("ibc/CCCBD7307FEB70B0CF7ADF8503F711F6741F41623D25BAD8CB736E03BE384264", ibc3Amount)

		reclaimCoins := sdk.NewCoins(onexCoin, ibc1, ibc2, ibc3)
		reclaimer, _ := sdk.AccAddressFromBech32("onomy1yc0lg97cy5e80jyajtkz0zke2rr4734anugf9g")

		keepers.BankKeeper.SendCoinsFromModuleToAccount(ctx, markettypes.ModuleName, reclaimer, reclaimCoins)

		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		ctx.Logger().Info("Upgrade complete")
		return vm, err
	}
}
