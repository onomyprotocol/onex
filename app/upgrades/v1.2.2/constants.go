package v1_2_2

import (
	"github.com/onomyprotocol/onex/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrades name.
	UpgradeName = "v1.2.2"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
}
