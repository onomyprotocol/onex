package v1_1_6

import (
	"github.com/onomyprotocol/onex/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrades name.
	UpgradeName   = "v1_1_6"
	UpgradeHeight = 389310
)

var Fork = upgrades.Fork{
	UpgradeName:    UpgradeName,
	UpgradeHeight:  UpgradeHeight,
	BeginForkLogic: RunForkLogic,
}
