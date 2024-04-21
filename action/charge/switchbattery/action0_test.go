package switchbattery

import (
	"fmt"
	"testing"
	"time"

	"github.com/bianxiaojie/rte/ctx"
	"github.com/bianxiaojie/rte/engine"
	"github.com/bianxiaojie/rte/entity/action"
	"github.com/bianxiaojie/rte/utils/ref"
	"github.com/bianxiaojie/wrsn/common/state"
	"github.com/bianxiaojie/wrsn/common/value"
)

// 第一步，自定义实体类型，并实现entity.go中定义的实体接口
type switchbatteryAction0Entity struct {
	id                  string
	shouldSwitchBattery bool
	state               state.WRSNEntityState
	energy              value.Energy
	maxEnergy           value.Energy
}

// --- 以下是调用SwitchbatteryAction0动作所需实现的BatterySwitchable接口方法
func (e *switchbatteryAction0Entity) Id() string {
	return e.id
}

func (e *switchbatteryAction0Entity) SetState(state state.WRSNEntityState) {
	e.state = state
}

func (e *switchbatteryAction0Entity) SetEnergy(energy value.Energy) {
	e.energy = energy
}

func (e *switchbatteryAction0Entity) GetMaxEnergy() value.Energy {
	return e.maxEnergy
}

// 第二步，为自定义的实体添加行为，在行为中调用SwitchBattery这一动作
// --- 以下是switchbatteryAction0Entity的行为，每个时间单位按照方法名后缀的优先级依次执行
func (e *switchbatteryAction0Entity) Run_0(context ctx.Context) {
	// 判断是否应该更换电池
	if e.shouldSwitchBattery {
		e.shouldSwitchBattery = false
		fmt.Printf("充电前的能量：%v\n", e.energy)
		param := SwitchbatteryParam0{}
		action.HandleNoneTargetAction[*SwitchbatteryAction0, BatterySwitchable](context.ActionHandler(), e, param)
		fmt.Printf("充电后的能量：%v\n", e.energy)
	}
}

// 第三步，启动仿真
func TestAction0(t *testing.T) {
	// 创建仿真引擎，并将时间间隔和停止时间都设置为1秒
	e := engine.MakeDefaultEngine(time.Second, time.Second)
	// 添加实体行为
	e.EntityManager().AddBehaviorByType(ref.ParseType[*switchbatteryAction0Entity]())
	// 添加实体
	e.EntityManager().AddEntity(&switchbatteryAction0Entity{
		id:                  "id",
		shouldSwitchBattery: true,
		state:               state.None,
		energy:              0,
		maxEnergy:           value.Joule,
	})
	// 启动仿真引擎并等待仿真停止
	e.Start()
	e.WaitStopped()
}
