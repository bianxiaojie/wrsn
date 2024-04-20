package simplesense

import (
	"fmt"
	"testing"
	"time"

	"github.com/bianxiaojie/rte/context"
	"github.com/bianxiaojie/rte/engine"
	"github.com/bianxiaojie/rte/entity/action"
	"github.com/bianxiaojie/rte/utils/ref"
	"github.com/bianxiaojie/wrsn/common/state"
	"github.com/bianxiaojie/wrsn/common/value"
)

// 第一步，自定义实体类型，并实现entity.go中定义的实体接口
type action0Entity struct {
	id          string
	shouldSense bool
	state       state.WRSNEntityState
	energy      value.Energy
}

func (e *action0Entity) Id() string {
	return e.id
}

func (e *action0Entity) SetState(state state.WRSNEntityState) {
	e.state = state
}

func (e *action0Entity) GetEnergy() value.Energy {
	return e.energy
}

func (e *action0Entity) SetEnergy(energy value.Energy) {
	e.energy = energy
}

// 在本例中，实体每秒感知耗能1焦耳
func (e *action0Entity) ComputeSensingEnergyConsumed(duration time.Duration) value.Energy {
	return value.Joule * value.Energy(duration/time.Second)
}

// 第二步，为自定义的实体添加行为，在行为中调用SimpleSense这一动作
func (e *action0Entity) Run_0(ctx context.Context) {
	param := Param0{
		Timeunit: ctx.Timer().GetTimeunit(),
	}
	// 判断是否应该感知
	if e.shouldSense {
		e.shouldSense = false
		fmt.Printf("感知前的能量：%v\n", e.energy)
		action.HandleNoneTargetAction[*Action0, Sensible](ctx.ActionHandler(), e, param)
		fmt.Printf("感知后的能量：%v\n", e.energy)
	}
}

// 第三步，启动仿真
func TestAction0(t *testing.T) {
	// 创建仿真引擎，并将时间间隔和停止时间都设置为1秒
	e := engine.MakeDefaultEngine(time.Second, time.Second)
	// 添加实体行为
	e.EntityManager().AddBehaviorByType(ref.ParseType[*action0Entity]())
	// 添加实体
	e.EntityManager().AddEntity(&action0Entity{
		id:          "id",
		shouldSense: true,
		state:       state.None,
		energy:      value.Joule,
	})
	// 启动仿真引擎并等待仿真停止
	e.Start()
	e.WaitStopped()
}
