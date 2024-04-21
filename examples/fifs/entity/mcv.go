package entity

import (
	"fmt"
	"time"

	"github.com/bianxiaojie/rte/ctx"
	"github.com/bianxiaojie/rte/entity/action"
	"github.com/bianxiaojie/rte/utils/ref"
	"github.com/bianxiaojie/wrsn/action/charge/onetoonecharge"
	"github.com/bianxiaojie/wrsn/action/charge/switchbattery"
	"github.com/bianxiaojie/wrsn/action/move/movewithenergy"
	"github.com/bianxiaojie/wrsn/common/state"
	"github.com/bianxiaojie/wrsn/common/value"
)

type MCV struct {
	id                      string                // MCV的id
	state                   state.WRSNEntityState // MCV的当前状态，比如正在感知或被充电等
	position                value.Position        // MCV的位置
	moveSpeed               value.Length          // MCV每秒移动距离
	moveEnergyConsumedSpeed value.Energy          // MCV每米移动耗能
	energy                  value.Energy          // MCV的剩余能量
	maxEnergy               value.Energy          // MCV的最大能量
	chargingSpeed           value.Energy          // MCV每秒充电耗能
	chargingRate            float64               // 充电效率
	command                 any                   // 存储BS发送的命令
}

func MakeMCV(id string, position value.Position, moveSpeed value.Length, moveEnergyConsumedSpeed value.Energy, maxEnergy value.Energy, chargingSpeed value.Energy, chargingRate float64) *MCV {
	mcv := &MCV{}
	mcv.id = id
	mcv.position = position
	mcv.moveSpeed = moveSpeed
	mcv.moveEnergyConsumedSpeed = moveEnergyConsumedSpeed
	mcv.energy = maxEnergy
	mcv.maxEnergy = maxEnergy
	mcv.chargingSpeed = chargingSpeed
	mcv.chargingRate = chargingRate
	return mcv
}

// --- 以下是调用各种动作所需实现的接口方法，动作与接口的映射如下所示:
// onetoonecharge.OneToOneChargeStagedAction1 -> onetoonecharge.OneToOneCharger
// switchbattery.SwitchbatteryAction0 -> switchbattery.BatterySwitchable
// command.CommandAction1 -> command.CommandTarget
// movewithenergy.MoveWithEnergyStagedAction0 -> movewithenergy.MovableWithEnergy
func (mcv *MCV) Id() string {
	return mcv.id
}

func (mcv *MCV) SetState(state state.WRSNEntityState) {
	mcv.state = state
}

func (mcv *MCV) SetPosition(position value.Position) {
	mcv.position = position
}

func (mcv *MCV) GetPosition() value.Position {
	return mcv.position
}

func (mcv *MCV) ComputeMovingDistance(duration time.Duration) value.Length {
	// 计算duration时间内MCV移动距离
	return mcv.moveSpeed * value.Length(float64(duration)/float64(time.Second))
}

func (mcv *MCV) SetEnergy(energy value.Energy) {
	mcv.energy = energy
}

func (mcv *MCV) GetEnergy() value.Energy {
	return mcv.energy
}

func (mcv *MCV) GetMaxEnergy() value.Energy {
	return mcv.maxEnergy
}

func (mcv *MCV) ComputeMovingEnergyConsumed(distance value.Length) value.Energy {
	// 计算MCV移动distance距离的耗能
	return mcv.moveEnergyConsumedSpeed * value.Energy(distance/value.Meter)
}

func (mcv *MCV) ComputeChargingEnergyConsumedAndCharged(target onetoonecharge.OneToOneChargingTarget, duration time.Duration) (value.Energy, value.Energy) {
	// 计算duration时间内MCV充电耗能
	chargingEnergy := mcv.chargingSpeed * value.Energy(float64(duration)/float64(time.Second))

	// 计算节点接收到的能量 = MCV充电耗能 * 充电效率
	chargedEnergy := chargingEnergy * value.Energy(mcv.chargingRate)
	return chargingEnergy, chargedEnergy
}

func (mcv *MCV) SetCommand(command any) {
	mcv.command = command
}

// --- 以下是MCV的行为，每个时间单位按照方法名后缀的优先级依次执行
func (mcv *MCV) HandleCommand_1(context ctx.Context) {
	// 如果没有收到BS的充电命令，则返回
	if mcv.command == nil {
		return
	}

	if moveChargeCommand, ok := mcv.command.(*MoveChargeCommand); ok {
		// 如果待充电的节点不存在，则返回
		entity, ok := context.EntityManager().GetEntityById(moveChargeCommand.SensorId)
		if !ok {
			return
		}
		sensor := entity.(*Sensor)

		// 移动到节点位置
		moveParam := movewithenergy.MoveWithEnergyStagedParam0{
			Timeunit:       context.Timer().GetTimeunit(),
			TargetPosition: sensor.GetPosition(),
		}
		action.HandleNoneTargetStagedAction[*movewithenergy.MoveWithEnergyStagedAction0, movewithenergy.MoveWithEnergyStage0, movewithenergy.MovableWithEnergy](context.ActionHandler(), mcv, moveParam)

		// 为节点充电
		chargeParam := onetoonecharge.OneToOneChargeStagedParam1WithEnergyLevel{
			Timeunit:    context.Timer().GetTimeunit(),
			EnergyLevel: sensor.GetMaxEnergy(),
		}
		fmt.Printf("%v, 充电前%v的能量: %v\n", context.Timer().GetTime(), sensor.Id(), sensor.GetEnergy())
		action.HandleOneTargetStagedAction[*onetoonecharge.OneToOneChargeStagedAction1, onetoonecharge.OneToOneChargeStage1, onetoonecharge.OneToOneCharger, onetoonecharge.OneToOneChargingTarget, onetoonecharge.OneToOneChargeStagedParam1](context.ActionHandler(), mcv, sensor.Id(), chargeParam)
		fmt.Printf("%v, 充电后%v的能量: %v\n", context.Timer().GetTime(), sensor.Id(), sensor.GetEnergy())
	} else if _, ok := mcv.command.(*SwitchBatteryCommand); ok {
		// 移动到BS
		bs, _ := context.EntityManager().GetEntitiesByType(ref.ParseType[*BS]())[0].(*BS)
		moveParam := movewithenergy.MoveWithEnergyStagedParam0{
			Timeunit:       context.Timer().GetTimeunit(),
			TargetPosition: bs.GetPosition(),
		}
		action.HandleNoneTargetStagedAction[*movewithenergy.MoveWithEnergyStagedAction0, movewithenergy.MoveWithEnergyStage0, movewithenergy.MovableWithEnergy](context.ActionHandler(), mcv, moveParam)

		// 更换电池
		switchBatteryParam := switchbattery.SwitchbatteryParam0{}
		action.HandleNoneTargetAction[*switchbattery.SwitchbatteryAction0, switchbattery.BatterySwitchable](context.ActionHandler(), mcv, switchBatteryParam)
		fmt.Printf("%v, %v回到BS\n", context.Timer().GetTime(), mcv.id)
	}
	mcv.command = nil
}

func (mcv *MCV) Rest_2(context ctx.Context) {
	// 如果MCV耗尽能量则报错
	if mcv.energy <= 0 {
		panic("MCV耗尽能量")
	}

	// 如果MCV正在执行其他行为则休息，比如MCV正在为充电充电
	if mcv.state != state.None {
		return
	}
}
