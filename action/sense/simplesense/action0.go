package simplesense

import (
	"time"

	"github.com/bianxiaojie/wrsn/common/state"
)

type SimpleSenseParam0 struct {
	Timeunit time.Duration
}

type SimpleSenseAction0 struct {
}

func (a *SimpleSenseAction0) Action(sensible Sensible, param SimpleSenseParam0) any {
	// 更新状态
	sensible.SetState(state.Sensing)

	// 更新能量
	sensingEnergyConsumed := sensible.ComputeSensingEnergyConsumed(param.Timeunit)
	sensible.SetEnergy(sensible.GetEnergy() - sensingEnergyConsumed)

	// 重置状态
	sensible.SetState(state.None)

	return nil
}
