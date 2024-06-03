package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/onomyprotocol/onex/app/upgrades"
)

// BeginBlockForks is intended to be ran in a chain upgrade.
func BeginBlockForks(ctx sdk.Context, app *App) {
	for _, fork := range Forks {
		if ctx.BlockHeight() == fork.UpgradeHeight {
			fork.BeginForkLogic(ctx, &upgrades.ForkKeepers{
				BankKeeper:   app.BankKeeper,
				MarketKeeper: app.MarketKeeper,
			})
			return
		}
	}
}
