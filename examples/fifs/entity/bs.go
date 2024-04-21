package entity

import (
	"time"

	"github.com/bianxiaojie/rte/ctx"
	"github.com/bianxiaojie/rte/entity/action"
	"github.com/bianxiaojie/rte/utils/ref"
	"github.com/bianxiaojie/wrsn/action/communicate/command"
	"github.com/bianxiaojie/wrsn/common/utils/mat"
	"github.com/bianxiaojie/wrsn/common/value"
)

type BS struct {
	id                 string            // BS的id
	position           value.Position    // BS的位置
	requestCommandList []*RequestCommand // 节点请求的列表
}

func MakeBS(id string, position value.Position) *BS {
	bs := &BS{}
	bs.id = id
	bs.position = position
	bs.requestCommandList = make([]*RequestCommand, 0)
	return bs
}

// --- 以下是调用各种动作所需实现的接口方法，动作与接口的映射如下所示:
// command.CommandAction1 -> command.CommandSource
// command.CommandAction1 -> command.CommandTarget
func (bs *BS) Id() string {
	return bs.id
}

func (bs *BS) GetPosition() value.Position {
	return bs.position
}

func (bs *BS) SetCommand(command any) {
	if requestCommand, ok := command.(*RequestCommand); ok {
		bs.requestCommandList = append(bs.requestCommandList, requestCommand)
	}
}

// --- 以下是BS的行为，每个时间单位按照方法名后缀的优先级依次执行
func (bs *BS) Schedule_0(context ctx.Context) {
	// 如果MCV不存在则报错
	mcvs := context.EntityManager().GetEntitiesByType(ref.ParseType[*MCV]())
	if len(mcvs) == 0 {
		panic("MCV不存在")
	}
	mcv := mcvs[0].(*MCV)

	// 如果MCV正在执行其他命令则返回
	if mcv.command != nil {
		return
	}

	for len(bs.requestCommandList) > 0 {
		// 1. 如果请求节点已死亡，则跳过节点
		requestCommand := bs.requestCommandList[0]
		entity, ok := context.EntityManager().GetEntityById(requestCommand.SensorId)
		if !ok {
			bs.requestCommandList = bs.requestCommandList[1:]
			continue
		}
		sensor := entity.(*Sensor)

		// 2. 如果MCV无法在节点死亡前移动到节点，则跳过节点
		timeunit := context.Timer().GetTimeunit()
		lifetime := time.Duration(mat.Floor(sensor.GetEnergy()/sensor.ComputeSensingEnergyConsumed(timeunit))) * timeunit
		moveDuration := time.Duration(mat.Ceil(mcv.GetPosition().DistanceTo(sensor.GetPosition())/mcv.ComputeMovingDistance(timeunit))) * timeunit
		if lifetime <= moveDuration {
			bs.requestCommandList = bs.requestCommandList[1:]
			continue
		}

		// 3. 向MCV发送移动充电或更换电池命令
		// 计算移动到节点的耗能
		moveEnergy := mcv.ComputeMovingEnergyConsumed(mcv.ComputeMovingDistance(moveDuration))
		// 计算单位时间的移动速度
		_, chargedSpeed := mcv.ComputeChargingEnergyConsumedAndCharged(sensor, timeunit)
		// 计算节点需要被补给的总能量
		chargedEnergy := sensor.GetMaxEnergy() - sensor.GetEnergy() + sensor.ComputeSensingEnergyConsumed(lifetime)
		// 计算充电时长
		chargedDuration := time.Duration(mat.Ceil(chargedEnergy/chargedSpeed)) * timeunit
		// 计算充电MCV耗费能量
		chargeEnergy, _ := mcv.ComputeChargingEnergyConsumedAndCharged(sensor, chargedDuration)
		// 计算移动到BS的时长
		moveBackDuration := time.Duration(mat.Ceil(sensor.GetPosition().DistanceTo(bs.position))/mcv.ComputeMovingDistance(timeunit)) * timeunit
		// 计算移动到BS的耗能
		moveBackEnergy := mcv.ComputeMovingEnergyConsumed(mcv.ComputeMovingDistance(moveBackDuration))
		// 如果MCV剩余能量 > 移动到节点的耗能 + 为节点充电耗能 + 移动回BS耗能，则向MCV发送移动充电命令
		if mcv.GetEnergy() > moveEnergy+chargeEnergy+moveBackEnergy {
			bs.requestCommandList = bs.requestCommandList[1:]
			commandParam := command.CommandParam1{
				Command: &MoveChargeCommand{
					SensorId: sensor.Id(),
				},
			}
			action.HandleOneTargetAction[*command.CommandAction1, command.CommandSource](context.ActionHandler(), bs, mcv.Id(), commandParam)
			return
		}
		break
	}

	// 如果MCV无法充电且不在BS，则向MCV发送回基站更换电池命令
	if mcv.GetPosition() != bs.position {
		commandParam := command.CommandParam1{
			Command: &SwitchBatteryCommand{},
		}
		action.HandleOneTargetAction[*command.CommandAction1, command.CommandSource](context.ActionHandler(), bs, mcv.Id(), commandParam)
	}
}
