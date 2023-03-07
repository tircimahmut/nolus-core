package app

import (
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func (app *App) RegisterUpgradeHandlers() {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	app.registerUpgradeV1_43(upgradeInfo)
	app.registerUpgradeV1_44(upgradeInfo)
	app.registerUpgradeV2_0(upgradeInfo)
}

// performs upgrade from v0.1.39 -> v0.1.43.
func (app *App) registerUpgradeV1_43(_ storetypes.UpgradeInfo) {
	const UpgradeV1_43Plan = "v0.1.43"
	app.UpgradeKeeper.SetUpgradeHandler(UpgradeV1_43Plan, func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Upgrade handler execution", "name", UpgradeV1_43Plan)
		return fromVM, nil
	})
}

// performs upgrade from v0.1.43 -> v0.1.44.
func (app *App) registerUpgradeV1_44(_ storetypes.UpgradeInfo) {
	const UpgradeV1_44Plan = "v0.1.44"
	app.UpgradeKeeper.SetUpgradeHandler(UpgradeV1_44Plan, func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Upgrade handler execution", "name", UpgradeV1_44Plan)
		return fromVM, nil
	})
}

// performs upgrade from v0.1.43 -> v0.2.2.
func (app *App) registerUpgradeV2_0(_ storetypes.UpgradeInfo) {
	const UpgradeV2_0Plan = "v0.2.0"
	app.UpgradeKeeper.SetUpgradeHandler(UpgradeV2_0Plan, func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("Upgrade handler execution", "name", UpgradeV2_0Plan)
		return fromVM, nil
	})
}