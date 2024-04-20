package movewithenergy

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
type action0StagedEntity struct {
	id         string
	shouldMove bool
	state      state.WRSNEntityState
	position   value.Position
	energy     value.Energy
}

func (e *action0StagedEntity) Id() string {
	return e.id
}

func (e *action0StagedEntity) SetState(state state.WRSNEntityState) {
	e.state = state
}

func (e *action0StagedEntity) SetPosition(position value.Position) {
	e.position = position
}

func (e *action0StagedEntity) GetPosition() value.Position {
	return e.position
}

func (e *action0StagedEntity) ComputeMovingDistance(duration time.Duration) value.Length {
	return value.Meter / value.Length(duration/time.Second)
}

func (e *action0StagedEntity) SetEnergy(energy value.Energy) {
	e.energy = energy
}

func (e *action0StagedEntity) GetEnergy() value.Energy {
	return e.energy
}

// 在本例中，实体每移动1米耗能1焦耳
func (e *action0StagedEntity) ComputeMovingEnergyConsumed(distance value.Length) value.Energy {
	return value.Joule / value.Energy(distance/value.Meter)
}

// 第二步，为自定义的实体添加行为，在行为中调用MoveWithEnergy这一动作
func (e *action0StagedEntity) Run_0(ctx context.Context) {
	// 判断是否应该移动
	if e.shouldMove {
		param := Param0Staged{
			Timeunit:       ctx.Timer().GetTimeunit(),
			TargetPosition: value.MakePosition(e.position.X()+value.Meter, e.position.Y()+value.Meter),
		}
		fmt.Printf("移动前的位置：%v, 能量： %v, 时间：%v\n", e.position, e.energy, ctx.Timer().GetTime())
		action.HandleNoneTargetStagedAction[*Action0Staged, Stage0, MovableWithEnergy](ctx.ActionHandler(), e, param)
		fmt.Printf("移动后的位置：%v, 能量： %v, 时间：%v\n", e.position, e.energy, ctx.Timer().GetTime())
	}
}

// 第三步，启动仿真
func TestAction0(t *testing.T) {
	// 创建仿真引擎，并将时间间隔设置为1秒，将停止时间都设置为2秒
	e := engine.MakeDefaultEngine(time.Second, 2*time.Second)
	// 添加实体行为
	e.EntityManager().AddBehaviorByType(ref.ParseType[*action0StagedEntity]())
	// 添加实体
	e.EntityManager().AddEntity(&action0StagedEntity{
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
