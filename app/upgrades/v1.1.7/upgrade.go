// Package v1_1_7 is contains chain upgrade of the corresponding version.
package v1_1_7 //nolint:revive,stylecheck // app version

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

		onexAmount, _ := sdk.NewIntFromString("65975935131986921086784")
		ibc1Amount, _ := sdk.NewIntFromString("17216827458309103931141")
		ibc2Amount, _ := sdk.NewIntFromString("20100000")

		onexCoin := sdk.NewCoin("aonex", onexAmount)
		ibc1 := sdk.NewCoin("ibc/5BDD8875CC2AF7BC842BE44236ACD576EA4F53C36347F74903B852060D6BF29A", ibc1Amount)
		ibc2 := sdk.NewCoin("ibc/CCCBD7307FEB70B0CF7ADF8503F711F6741F41623D25BAD8CB736E03BE384264", ibc2Amount)

		returnCoins := sdk.NewCoins(onexCoin, ibc1, ibc2)
		returner, _ := sdk.AccAddressFromBech32("onomy1yc0lg97cy5e80jyajtkz0zke2rr4734anugf9g")

		keepers.BankKeeper.SendCoinsFromAccountToModule(ctx, returner, markettypes.ModuleName, returnCoins)

		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		ctx.Logger().Info("Upgrade complete")
		return vm, err
	}
}
