package switchbattery

import "github.com/bianxiaojie/wrsn/common/state"

type Param0 struct {
}

type Action0 struct {
}

func (a *Action0) Action(batterySwitchable BatterySwitchable, param Param0) any {
	// 更新状态
	batterySwitchable.SetState(state.SwitchingBattery)

	// 更新能量
	batterySwitchable.SetEnergy(batterySwitchable.GetMaxEnergy())

	// 重置状态
	batterySwitchable.SetState(state.None)

	return nil
}
