package v1_1_7

import (
	"github.com/onomyprotocol/onex/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrades name.
	UpgradeName = "v1.1.7"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
}
