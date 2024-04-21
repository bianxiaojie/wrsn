package switchbattery

import (
	"github.com/bianxiaojie/rte/entity"
	"github.com/bianxiaojie/wrsn/common/state"
	"github.com/bianxiaojie/wrsn/common/value"
)

type BatterySwitchable interface {
	entity.Entity
	SetState(state.WRSNEntityState) // 设置状态
	SetEnergy(value.Energy)         // 更新实体能量
	GetMaxEnergy() value.Energy     // 获取实体最大能量
}
