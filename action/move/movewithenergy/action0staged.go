package movewithenergy

import (
	"time"

	"github.com/bianxiaojie/wrsn/common/state"
	"github.com/bianxiaojie/wrsn/common/value"
)

// 移动参数
type Param0Staged struct {
	Timeunit       time.Duration  // 仿真的时间单位
	TargetPosition value.Position // 目标位置
}

// 移动Stage
type Stage0 struct {
	currentPosition value.Position // 实体当前位置
	targetPosition  value.Position // 目标位置
}

func (s Stage0) IsLastStage() bool {
	return s.currentPosition == s.targetPosition
}

func (s Stage0) GetReturnedValue() any {
	return nil
}

// 移动Action
type Action0Staged struct {
}

func (a *Action0Staged) MakeStage(m MovableWithEnergy, p Param0Staged) Stage0 {
	return Stage0{
		currentPosition: m.GetPosition(),
		targetPosition:  p.TargetPosition,
	}
}

func (a *Action0Staged) ActionStage(movable MovableWithEnergy, param Param0Staged, stage Stage0) Stage0 {
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
