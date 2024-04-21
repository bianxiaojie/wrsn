package switchbattery

import "github.com/bianxiaojie/wrsn/common/state"

type SwitchbatteryParam0 struct {
}

type SwitchbatteryAction0 struct {
}

func (a *SwitchbatteryAction0) Action(batterySwitchable BatterySwitchable, param SwitchbatteryParam0) any {
	// 更新状态
	batterySwitchable.SetState(state.SwitchingBattery)

	// 更新能量
	batterySwitchable.SetEnergy(batterySwitchable.GetMaxEnergy())

	// 重置状态
	batterySwitchable.SetState(state.None)

	return nil
}
