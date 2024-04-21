package entity

import (
	"time"

	"github.com/bianxiaojie/rte/ctx"
	"github.com/bianxiaojie/rte/entity/action"
	"github.com/bianxiaojie/wrsn/action/sense/simplesense"
	"github.com/bianxiaojie/wrsn/common/state"
	"github.com/bianxiaojie/wrsn/common/value"
)

type Sensor struct {
	id     string
	state  state.WRSNEntityState
	energy value.Energy
}

func MakeSensor(id string, energy value.Energy) *Sensor {
	sensor := &Sensor{}
	sensor.id = id
	sensor.state = state.None
	sensor.energy = energy
	return sensor
}

// --- 以下是实现simplesense.Sensible接口所需实现的方法
func (s *Sensor) Id() string {
	return s.id
}

func (s *Sensor) SetState(state state.WRSNEntityState) {
	s.state = state
}

func (s *Sensor) GetEnergy() value.Energy {
	return s.energy
}

func (s *Sensor) SetEnergy(energy value.Energy) {
	s.energy = energy
}

// 能耗速率为每秒1J
func (s *Sensor) ComputeSensingEnergyConsumed(duration time.Duration) value.Energy {
	return value.Joule * value.Energy(duration/time.Second)
}

// --- 以下是Sensor的行为，每个时间单位按照方法名后缀的优先级依次执行
// 节点耗能行为
func (s *Sensor) Sense_0(context ctx.Context) {
	param := simplesense.SimpleSenseParam0{
		Timeunit: context.Timer().GetTimeunit(),
	}
	action.HandleNoneTargetAction[*simplesense.SimpleSenseAction0, simplesense.Sensible](context.ActionHandler(), s, param)
}

// 节点状态判断行为，如果能量下降到0以下则将节点移除，节点被移除后不会再执行任何行为，包括Sense_0和RemoveIfDead_1
func (s *Sensor) RemoveIfDead_1(context ctx.Context) {
	if s.energy <= 0 {
		context.EntityManager().RemoveEntityById(s.id)
	}
}
