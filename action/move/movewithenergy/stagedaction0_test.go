package movewithenergy

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
type moveWithEnergyStagedAction0Entity struct {
	id         string
	shouldMove bool
	state      state.WRSNEntityState
	position   value.Position
	energy     value.Energy
}

// --- 以下是调用MoveWithEnergyStagedAction0动作所需实现的MovableWithEnergy接口方法
func (e *moveWithEnergyStagedAction0Entity) Id() string {
	return e.id
}

func (e *moveWithEnergyStagedAction0Entity) SetState(state state.WRSNEntityState) {
	e.state = state
}

func (e *moveWithEnergyStagedAction0Entity) SetPosition(position value.Position) {
	e.position = position
}

func (e *moveWithEnergyStagedAction0Entity) GetPosition() value.Position {
	return e.position
}

func (e *moveWithEnergyStagedAction0Entity) ComputeMovingDistance(duration time.Duration) value.Length {
	return value.Meter * value.Length(duration/time.Second)
}

func (e *moveWithEnergyStagedAction0Entity) SetEnergy(energy value.Energy) {
	e.energy = energy
}

func (e *moveWithEnergyStagedAction0Entity) GetEnergy() value.Energy {
	return e.energy
}

// 在本例中，实体每移动1米耗能1焦耳
func (e *moveWithEnergyStagedAction0Entity) ComputeMovingEnergyConsumed(distance value.Length) value.Energy {
	return value.Joule * value.Energy(distance/value.Meter)
}

// 第二步，自定义实体的行为
// --- 以下是movewithenergyStagedAction0Entity的行为，每个时间单位按照方法名后缀的优先级依次执行
func (e *moveWithEnergyStagedAction0Entity) Run_0(context ctx.Context) {
	// 判断是否应该移动
	if e.shouldMove {
		param := MoveWithEnergyStagedParam0{
			Timeunit:       context.Timer().GetTimeunit(),
			TargetPosition: value.MakePosition(e.position.X()+value.Meter, e.position.Y()+value.Meter),
		}
		fmt.Printf("移动前的位置：%v, 能量： %v, 时间：%v\n", e.position, e.energy, context.Timer().GetTime())
		action.HandleNoneTargetStagedAction[*MoveWithEnergyStagedAction0, MoveWithEnergyStage0, MovableWithEnergy](context.ActionHandler(), e, param)
		fmt.Printf("移动后的位置：%v, 能量： %v, 时间：%v\n", e.position, e.energy, context.Timer().GetTime())
	}
}

// 第三步，启动仿真
func TestAction0(t *testing.T) {
	// 创建仿真引擎，并将时间间隔设置为1秒，将停止时间都设置为2秒
	e := engine.MakeDefaultEngine(time.Second, 2*time.Second)
	// 添加实体行为
	e.EntityManager().AddBehaviorByType(ref.ParseType[*moveWithEnergyStagedAction0Entity]())
	// 添加实体
	e.EntityManager().AddEntity(&moveWithEnergyStagedAction0Entity{
		id:         "id",
		shouldMove: true,
		state:      state.None,
		position:   value.MakePosition(0, 0),
		energy:     3 * value.Joule,
	})
	// 启动仿真引擎并等待仿真停止
	e.Start()
	e.WaitStopped()
}
