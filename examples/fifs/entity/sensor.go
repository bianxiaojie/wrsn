package entity

import (
	"fmt"
	"time"

	"github.com/bianxiaojie/rte/ctx"
	"github.com/bianxiaojie/rte/entity/action"
	"github.com/bianxiaojie/rte/utils/ref"
	"github.com/bianxiaojie/wrsn/action/communicate/command"
	"github.com/bianxiaojie/wrsn/action/sense/simplesense"
	"github.com/bianxiaojie/wrsn/common/state"
	"github.com/bianxiaojie/wrsn/common/value"
)

type Sensor struct {
	id                       string                // 节点的id
	state                    state.WRSNEntityState // 节点的当前状态，比如正在感知或被充电等
	position                 value.Position        // 节点的位置
	energy                   value.Energy          // 节点的剩余能量
	maxEnergy                value.Energy          // 节点的最大能量
	senseEnergyConsumedSpeed value.Energy          // 节点每秒耗能速度
	threshold                float64               // 节点发起请求的阈值百分比
	hasSentRequest           bool                  // 标记节点是否已经发起请求，避免重复发请求，被充电后需要重置为false
}

func MakeSensor(id string, position value.Position, maxEnergy value.Energy, senseEnergyConsumedSpeed value.Energy, threshold float64) *Sensor {
	sensor := &Sensor{}
	sensor.id = id
	sensor.position = position
	sensor.energy = maxEnergy
	sensor.maxEnergy = maxEnergy
	sensor.senseEnergyConsumedSpeed = senseEnergyConsumedSpeed
	sensor.threshold = threshold
	return sensor
}

// --- 以下是调用各种动作所需实现的接口方法，动作与接口的映射如下所示:
// simplesense.SimpleSenseAction0 -> simplesense.Sensible
// onetoonecharge.OneToOneChargeStagedAction1 -> onetoonecharge.OneToOneChargingTarget
// command.CommandAction1 -> command.CommandSource
func (s *Sensor) Id() string {
	return s.id
}

func (s *Sensor) SetState(state state.WRSNEntityState) {
	s.state = state
}

func (s *Sensor) GetPosition() value.Position {
	return s.position
}

func (s *Sensor) SetEnergy(energy value.Energy) {
	s.energy = energy
}

func (s *Sensor) GetEnergy() value.Energy {
	return s.energy
}

func (s *Sensor) GetMaxEnergy() value.Energy {
	return s.maxEnergy
}

func (e *Sensor) ComputeSensingEnergyConsumed(duration time.Duration) value.Energy {
	return e.senseEnergyConsumedSpeed * value.Energy(float64(duration)/float64(time.Second))
}

// --- 以下是Sensor的行为，每个时间单位按照方法名后缀的优先级依次执行
func (s *Sensor) Sense_2(context ctx.Context) {
	// 如果节点耗尽能量则将节点移除
	if s.energy <= 0 {
		context.EntityManager().RemoveEntityById(s.id)
	}

	// 如果节点正在执行其他行为则不感知，比如节点正在被充电
	if s.state != state.None {
		return
	}

	// 感知耗能
	param := simplesense.SimpleSenseParam0{
		Timeunit: context.Timer().GetTimeunit(),
	}
	action.HandleNoneTargetAction[*simplesense.SimpleSenseAction0, simplesense.Sensible](context.ActionHandler(), s, param)
}

func (s *Sensor) SendRequest_3(context ctx.Context) {
	// 如果节点正在被充电，则将hasSentRequest重置为false，这样节点在被充电后又可以重新发起请求
	if s.state == state.Charged {
		s.hasSentRequest = false
		return
	}

	// 如果节点已发起请求，或能量高于请求阈值，则不发请求
	if s.hasSentRequest || s.energy > s.maxEnergy*value.Energy(s.threshold) {
		return
	}

	fmt.Printf("%v, %v发起请求的能量: %v\n", context.Timer().GetTime(), s.id, s.energy)

	// 将hasSentRequest标记为true，避免重复发送请求
	s.hasSentRequest = true
	param := command.CommandParam1{
		Command: &RequestCommand{
			SensorId: s.id,
		},
	}
	bs := context.EntityManager().GetEntitiesByType(ref.ParseType[*BS]())[0]
	action.HandleOneTargetAction[*command.CommandAction1, command.CommandSource](context.ActionHandler(), s, bs.Id(), param)
}
