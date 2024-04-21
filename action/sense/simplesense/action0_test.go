package simplesense

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
type simpleSenseAction0Entity struct {
	id          string
	shouldSense bool
	state       state.WRSNEntityState
	energy      value.Energy
}

// --- 以下是调用SimpleSenseAction0动作所需实现的Sensible接口方法
func (e *simpleSenseAction0Entity) Id() string {
	return e.id
}

func (e *simpleSenseAction0Entity) SetState(state state.WRSNEntityState) {
	e.state = state
}

func (e *simpleSenseAction0Entity) GetEnergy() value.Energy {
	return e.energy
}

func (e *simpleSenseAction0Entity) SetEnergy(energy value.Energy) {
	e.energy = energy
}

// 在本例中，实体每秒感知耗能1焦耳
func (e *simpleSenseAction0Entity) ComputeSensingEnergyConsumed(duration time.Duration) value.Energy {
	return value.Joule * value.Energy(duration/time.Second)
}

// 第二步，为自定义的实体添加行为，在行为中调用SimpleSense这一动作
// --- 以下是simpleSenseAction0Entity的行为，每个时间单位按照方法名后缀的优先级依次执行
func (e *simpleSenseAction0Entity) Run_0(context ctx.Context) {
	param := SimpleSenseParam0{
		Timeunit: context.Timer().GetTimeunit(),
	}
	// 判断是否应该感知
	if e.shouldSense {
		e.shouldSense = false
		fmt.Printf("感知前的能量：%v\n", e.energy)
		action.HandleNoneTargetAction[*SimpleSenseAction0, Sensible](context.ActionHandler(), e, param)
		fmt.Printf("感知后的能量：%v\n", e.energy)
	}
}

// 第三步，启动仿真
func TestAction0(t *testing.T) {
	// 创建仿真引擎，并将时间间隔和停止时间都设置为1秒
	e := engine.MakeDefaultEngine(time.Second, time.Second)
	// 添加实体行为
	e.EntityManager().AddBehaviorByType(ref.ParseType[*simpleSenseAction0Entity]())
	// 添加实体
	e.EntityManager().AddEntity(&simpleSenseAction0Entity{
		id:          "id",
		shouldSense: true,
		state:       state.None,
		energy:      value.Joule,
	})
	// 启动仿真引擎并等待仿真停止
	e.Start()
	e.WaitStopped()
}
