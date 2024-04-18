package simplesense

import (
	"time"

	"github.com/bianxiaojie/wrsn/common/state"
)

type Param0 struct {
	Timeunit time.Duration
}

type Action0 struct {
}

func (a *Action0) Action(sensible Sensible, param Param0) any {
	// 更新状态
	sensible.SetState(state.Sensing)

	// 更新能量
	sensingEnergyConsumed := sensible.ComputeSensingEnergyConsumed(param.Timeunit)
	sensible.SetEnergy(sensible.GetEnergy() - sensingEnergyConsumed)

	// 重置状态
	sensible.SetState(state.None)

	return nil
}
