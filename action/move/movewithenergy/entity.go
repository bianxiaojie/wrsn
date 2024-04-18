package movewithenergy

import (
	"time"

	"github.com/bianxiaojie/rte/entity"
	"github.com/bianxiaojie/wrsn/common/state"
	"github.com/bianxiaojie/wrsn/common/value"
)

// 可移动且有能量的实体
type MovableWithEnergy interface {
	entity.Entity
	SetState(state.WRSNEntityState)                        // 在移动中将实体状态更新为Moving，移动后更新为None
	SetPosition(value.Position)                            // 更新实体位置
	GetPosition() value.Position                           // 获取实体位置
	ComputeMovingDistance(time.Duration) value.Length      // 计算给定时间内实体的移动距离
	SetEnergy(value.Energy)                                // 更新实体能量
	GetEnergy() value.Energy                               // 获取实体能量
	ComputeMovingEnergyConsumed(value.Length) value.Energy // 计算实体移动能耗
}
