package simplesense

import (
	"time"

	"github.com/bianxiaojie/rte/entity"
	"github.com/bianxiaojie/wrsn/common/state"
	"github.com/bianxiaojie/wrsn/common/value"
)

type Sensible interface {
	entity.Entity
	SetState(state.WRSNEntityState)                          // 在感知中实体状态更新为Sensing，感知后更新为None
	SetEnergy(value.Energy)                                  // 更新实体能量
	GetEnergy() value.Energy                                 // 获取实体能量
	ComputeSensingEnergyConsumed(time.Duration) value.Energy // 计算给定时间内实体的感知耗费的能量
}
