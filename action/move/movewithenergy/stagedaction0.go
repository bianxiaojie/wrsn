package movewithenergy

import (
	"time"

	"github.com/bianxiaojie/wrsn/common/state"
	"github.com/bianxiaojie/wrsn/common/value"
)

// 移动参数
type MoveWithEnergyStagedParam0 struct {
	Timeunit       time.Duration  // 仿真的时间单位
	TargetPosition value.Position // 目标位置
}

// 移动Stage
type MoveWithEnergyStage0 struct {
	currentPosition value.Position // 实体当前位置
	targetPosition  value.Position // 目标位置
}

func (s MoveWithEnergyStage0) IsLastStage() bool {
	return s.currentPosition == s.targetPosition
}

func (s MoveWithEnergyStage0) GetReturnedValue() any {
	return nil
}

// 移动Action
type MoveWithEnergyStagedAction0 struct {
}

func (a *MoveWithEnergyStagedAction0) MakeStage(m MovableWithEnergy, p MoveWithEnergyStagedParam0) MoveWithEnergyStage0 {
	return MoveWithEnergyStage0{
		currentPosition: m.GetPosition(),
		targetPosition:  p.TargetPosition,
	}
}

func (a *MoveWithEnergyStagedAction0) ActionStage(movable MovableWithEnergy, param MoveWithEnergyStagedParam0, stage MoveWithEnergyStage0) MoveWithEnergyStage0 {
	if !stage.IsLastStage() {
		// 更新状态
		movable.SetState(state.Moving)

		// 计算位置
		moveDistance := movable.ComputeMovingDistance(param.Timeunit)
		distance := movable.GetPosition().DistanceTo(param.TargetPosition)
		var newPosition value.Position
		// 与目标剩余距离不超过移动速率，则直接将实体位置更新为目标位置，否则根据极角和距离计算新位置
		if moveDistance >= distance {
			newPosition = param.TargetPosition
		} else {
			polarAngle := movable.GetPosition().PolarAngleTo(param.TargetPosition)
			newPosition = movable.GetPosition().PositionTo(polarAngle, moveDistance)
		}

		// 更新位置
		movable.SetPosition(newPosition)

		// 更新能量，由于最小时间单位不可分割，即使实际距离小于移动速率，也要耗费单位时间内移动所需的能量
		movingEnergyConsumed := movable.ComputeMovingEnergyConsumed(moveDistance)
		movable.SetEnergy(movable.GetEnergy() - movingEnergyConsumed)

		stage.currentPosition = movable.GetPosition()
	}

	if stage.IsLastStage() {
		// 重置状态
		movable.SetState(state.None)
	}
	return stage
}
