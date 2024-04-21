package onetoonecharge

import (
	"time"

	"github.com/bianxiaojie/rte/entity"
	"github.com/bianxiaojie/wrsn/common/state"
	"github.com/bianxiaojie/wrsn/common/value"
)

// 一对一充电装置
type OneToOneCharger interface {
	entity.Entity
	SetState(state.WRSNEntityState) // 在充点中将实体状态更新为Charging，充电后更新为None
	GetPosition() value.Position    // 获取实体位置
	SetEnergy(value.Energy)         // 更新实体能量
	GetEnergy() value.Energy        // 获取实体能量
	// 计算给定时间内为充电目标充电耗费的能量和目标收到的能量
	ComputeChargingEnergyConsumedAndCharged(OneToOneChargingTarget, time.Duration) (value.Energy, value.Energy)
}

// 一对一充电目标
type OneToOneChargingTarget interface {
	entity.Entity
	SetState(state.WRSNEntityState) // 在充点中将实体状态更新为Charged，充电后更新为None
	GetPosition() value.Position    // 获取实体位置
	SetEnergy(value.Energy)         // 更新实体能量
	GetEnergy() value.Energy        // 获取实体能量
	GetMaxEnergy() value.Energy     // 获取实体最大能量
}
