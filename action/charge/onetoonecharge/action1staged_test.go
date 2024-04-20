package onetoonecharge

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

// 第一步，自定义充电装置实体类型，并实现entity.go中定义的实体接口
type action1StagedChargerEntity struct {
	id           string
	shouldCharge bool
	state        state.WRSNEntityState
	position     value.Position
	energy       value.Energy
}

func (e *action1StagedChargerEntity) Id() string {
	return e.id
}

func (e *action1StagedChargerEntity) SetState(state state.WRSNEntityState) {
	e.state = state
}

func (e *action1StagedChargerEntity) GetPosition() value.Position {
	return e.position
}

func (e *action1StagedChargerEntity) SetEnergy(energy value.Energy) {
	e.energy = energy
}

func (e *action1StagedChargerEntity) GetEnergy() value.Energy {
	return e.energy
}

func (e *action1StagedChargerEntity) ComputeChargingEnergyConsumedAndCharged(target OneToOneChargingTarget,
	duration time.Duration) (value.Energy, value.Energy) {
	// 假设充电效率为70%
	ratio := 0.7
	// 假设Charger充电每秒耗电为1J
	chargingEnergy := value.Joule / value.Energy(duration/time.Second)
	// Target接收到的能量 = Charger耗能 * 充电效率
	chargedEnergy := chargingEnergy * value.Energy(ratio)
	return chargingEnergy, chargedEnergy
}

// 第二步，自定义充电行为
func (e *action1StagedChargerEntity) Charge_0(ctx context.Context) {
	// 判断是否应该充电
	if e.shouldCharge {
		e.shouldCharge = false
		// 实际在决定充电目标时应当是计算出来的，这里固定为"target"只是为了方便
		param := Param1StagedWithDuration{
			Timeunit: ctx.Timer().GetTimeunit(),
			Duration: 3 * time.Second,
		}
		fmt.Printf("充电前充电装置的能量: %v, 时间: %v\n", e.energy, ctx.Timer().GetTime())
		action.HandleOneTargetStagedAction[*Action1Staged, Stage1, OneToOneCharger, OneToOneChargingTarget, Param1Staged](ctx.ActionHandler(), e, "target", param)
		fmt.Printf("充电后充电装置的能量: %v, 时间: %v\n", e.energy, ctx.Timer().GetTime())
	}
}

func (e *action1StagedChargerEntity) Monitor_1(ctx context.Context) {
	// 在每个单位时间结束时监控充电装置的能量值
	fmt.Printf("充电装置的剩余能量为: %v, 时间: %v\n", e.energy, ctx.Timer().GetTime())
}

// 第三步，自定义充电目标实体类型，并实现entity.go中定义的实体接口
type action1StagedTargetEntity struct {
	id        string
	state     state.WRSNEntityState
	position  value.Position
	energy    value.Energy
	maxEnergy value.Energy
}

func (e *action1StagedTargetEntity) Id() string {
	return e.id
}

func (e *action1StagedTargetEntity) SetState(state state.WRSNEntityState) {
	e.state = state
}

func (e *action1StagedTargetEntity) GetPosition() value.Position {
	return e.position
}

func (e *action1StagedTargetEntity) SetEnergy(energy value.Energy) {
	e.energy = energy
}

func (e *action1StagedTargetEntity) GetEnergy() value.Energy {
	return e.energy
}

func (e *action1StagedTargetEntity) GetMaxEnergy() value.Energy {
	return e.maxEnergy
}

// 第四步，自定义充电目标的行为
func (e *action1StagedTargetEntity) Monitor_1(ctx context.Context) {
	// 在每个单位时间结束时监控充电目标的能量值
	fmt.Printf("充电目标的剩余能量为: %v, 时间: %v\n", e.energy, ctx.Timer().GetTime())
}

// 第五步，启动仿真
func TestAction0(t *testing.T) {
	// 创建仿真引擎，并将时间间隔设置为1秒，将停止时间都设置为3秒
	e := engine.MakeDefaultEngine(time.Second, 3*time.Second)
	// 添加实体行为
	e.EntityManager().AddBehaviorByType(ref.ParseType[*action1StagedChargerEntity]())
	e.EntityManager().AddBehaviorByType(ref.ParseType[*action1StagedTargetEntity]())
	// 添加实体
	e.EntityManager().AddEntity(&action1StagedChargerEntity{
		id:           "charger",
		shouldCharge: true,
		state:        state.None,
		energy:       1000 * value.Joule,
	})
	e.EntityManager().AddEntity(&action1StagedTargetEntity{
		id:        "target",
		state:     state.None,
		energy:    10 * value.Joule,
		maxEnergy: 100 * value.Joule,
	})
	// 启动仿真引擎并等待仿真停止
	e.Start()
	e.WaitStopped()
}
