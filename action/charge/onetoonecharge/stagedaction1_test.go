package onetoonecharge

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

// 第一步，自定义充电装置实体类型，并实现entity.go中定义的实体接口
type oneToOneChargeStagedAction1ChargerEntity struct {
	id           string
	shouldCharge bool
	state        state.WRSNEntityState
	position     value.Position
	energy       value.Energy
}

// --- 以下是调用OneToOneChargeStagedAction1动作所需实现的OneToOneCharger接口方法
func (e *oneToOneChargeStagedAction1ChargerEntity) Id() string {
	return e.id
}

func (e *oneToOneChargeStagedAction1ChargerEntity) SetState(state state.WRSNEntityState) {
	e.state = state
}

func (e *oneToOneChargeStagedAction1ChargerEntity) GetPosition() value.Position {
	return e.position
}

func (e *oneToOneChargeStagedAction1ChargerEntity) SetEnergy(energy value.Energy) {
	e.energy = energy
}

func (e *oneToOneChargeStagedAction1ChargerEntity) GetEnergy() value.Energy {
	return e.energy
}

func (e *oneToOneChargeStagedAction1ChargerEntity) ComputeChargingEnergyConsumedAndCharged(target OneToOneChargingTarget,
	duration time.Duration) (value.Energy, value.Energy) {
	// 假设充电效率为70%
	rate := 0.7
	// 假设Charger充电每秒耗电为1J
	chargingEnergy := value.Joule / value.Energy(duration/time.Second)
	// Target接收到的能量 = Charger耗能 * 充电效率
	chargedEnergy := chargingEnergy * value.Energy(rate)
	return chargingEnergy, chargedEnergy
}

// 第二步，自定义充电行为
// --- 以下是oneToOneChargeStagedAction1ChargerEntity的行为，每个时间单位按照方法名后缀的优先级依次执行
func (e *oneToOneChargeStagedAction1ChargerEntity) Charge_0(context ctx.Context) {
	// 判断是否应该充电
	if e.shouldCharge {
		e.shouldCharge = false
		// 实际在决定充电目标时应当是计算出来的，这里固定为"target"只是为了方便
		param := OneToOneChargeStagedParam1WithDuration{
			Timeunit: context.Timer().GetTimeunit(),
			Duration: 3 * time.Second,
		}
		fmt.Printf("充电前充电装置的能量: %v, 时间: %v\n", e.energy, context.Timer().GetTime())
		action.HandleOneTargetStagedAction[*OneToOneChargeStagedAction1, OneToOneChargeStage1, OneToOneCharger, OneToOneChargingTarget, OneToOneChargeStagedParam1](context.ActionHandler(), e, "target", param)
		fmt.Printf("充电后充电装置的能量: %v, 时间: %v\n", e.energy, context.Timer().GetTime())
	}
}

func (e *oneToOneChargeStagedAction1ChargerEntity) Monitor_1(context ctx.Context) {
	// 在每个单位时间结束时监控充电装置的能量值
	fmt.Printf("充电装置的剩余能量为: %v, 时间: %v\n", e.energy, context.Timer().GetTime())
}

// 第三步，自定义充电目标实体类型，并实现entity.go中定义的实体接口
type oneToOneChargeStagedAction1TargetEntity struct {
	id        string
	state     state.WRSNEntityState
	position  value.Position
	energy    value.Energy
	maxEnergy value.Energy
}

// --- 以下是调用OneToOneChargeStagedAction1动作所需实现的OneToOneChargingTarget接口方法
func (e *oneToOneChargeStagedAction1TargetEntity) Id() string {
	return e.id
}

func (e *oneToOneChargeStagedAction1TargetEntity) SetState(state state.WRSNEntityState) {
	e.state = state
}

func (e *oneToOneChargeStagedAction1TargetEntity) GetPosition() value.Position {
	return e.position
}

func (e *oneToOneChargeStagedAction1TargetEntity) SetEnergy(energy value.Energy) {
	e.energy = energy
}

func (e *oneToOneChargeStagedAction1TargetEntity) GetEnergy() value.Energy {
	return e.energy
}

func (e *oneToOneChargeStagedAction1TargetEntity) GetMaxEnergy() value.Energy {
	return e.maxEnergy
}

// 第四步，自定义充电目标的行为
// --- 以下是action1StagedTargetEntity的行为，每个时间单位按照方法名后缀的优先级依次执行
func (e *oneToOneChargeStagedAction1TargetEntity) Monitor_1(context ctx.Context) {
	// 在每个单位时间结束时监控充电目标的能量值
	fmt.Printf("充电目标的剩余能量为: %v, 时间: %v\n", e.energy, context.Timer().GetTime())
}

// 第五步，启动仿真
func TestAction0(t *testing.T) {
	// 创建仿真引擎，并将时间间隔设置为1秒，将停止时间都设置为3秒
	e := engine.MakeDefaultEngine(time.Second, 3*time.Second)
	// 添加实体行为
	e.EntityManager().AddBehaviorByType(ref.ParseType[*oneToOneChargeStagedAction1ChargerEntity]())
	e.EntityManager().AddBehaviorByType(ref.ParseType[*oneToOneChargeStagedAction1TargetEntity]())
	// 添加实体
	e.EntityManager().AddEntity(&oneToOneChargeStagedAction1ChargerEntity{
		id:           "charger",
		shouldCharge: true,
		state:        state.None,
		energy:       1000 * value.Joule,
	})
	e.EntityManager().AddEntity(&oneToOneChargeStagedAction1TargetEntity{
		id:        "target",
		state:     state.None,
		energy:    10 * value.Joule,
		maxEnergy: 100 * value.Joule,
	})
	// 启动仿真引擎并等待仿真停止
	e.Start()
	e.WaitStopped()
}
