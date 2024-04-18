package onetoonecharge

import (
	"fmt"
	"time"

	"github.com/bianxiaojie/wrsn/common/state"
	"github.com/bianxiaojie/wrsn/common/value"
)

type Param1Staged interface {
	GetTimeunit() time.Duration // 仿真的时间单位
}

// 充电参数，在需要指定充电时长时使用
type Param1StagedWithDuration struct {
	Timeunit time.Duration // 仿真的时间单位
	Duration time.Duration // 充电总时长
}

func (p Param1StagedWithDuration) GetTimeunit() time.Duration {
	return p.Timeunit
}

// 充电参数，在需要指定充电后充电方剩余能量时使用
type Param1StagedWithEnergyLevel struct {
	Timeunit    time.Duration // 仿真的时间单位
	EnergyLevel value.Energy  // 充电后充电方剩余能量，比如为最大能量
}

func (p Param1StagedWithEnergyLevel) GetTimeunit() time.Duration {
	return p.Timeunit
}

// 充电参数，在需要指定充电方消耗能量时使用
type Param1StagedWithEnergyConsumed struct {
	Timeunit       time.Duration // 仿真的时间单位
	energyConsumed value.Energy  // 充电方需要消耗的能量
}

func (p Param1StagedWithEnergyConsumed) GetTimeunit() time.Duration {
	return p.Timeunit
}

// 充电参数，在需要指定被充电方收到能量时使用
type Param1StagedWithEnergyCharged struct {
	Timeunit      time.Duration // 仿真的时间单位
	energyCharged value.Energy  // 被充电方需要收到的能量
}

func (p Param1StagedWithEnergyCharged) GetTimeunit() time.Duration {
	return p.Timeunit
}

// 参数模式
type mode int64

const (
	durationMode       mode = iota // 充电时长模式
	energyLevelMode                // 充电后充电方剩余能量模式
	energyConsumedMode             // 充电方消耗能量模式
	energyChargedMode              // 被充电方收到能量模式
)

type Stage1 struct {
	mode mode
	// 被充电方不存在
	targetNil bool
	// 充电时长模式
	durationLeft time.Duration
	// 充电后充电方剩余能量模式
	targetEnergy value.Energy
	energyLevel  value.Energy
	// 充电方消耗能量模式
	energyConsumed value.Energy
	// 被充电方收到能量模式
	energyCharged value.Energy
}

func (s Stage1) IsLastStage() bool {
	// 目标不存在则直接返回
	if s.targetNil {
		return true
	}

	switch s.mode {
	case durationMode:
		return s.durationLeft <= 0
	case energyLevelMode:
		return s.targetEnergy >= s.energyLevel
	case energyConsumedMode:
		return s.energyConsumed <= 0
	case energyChargedMode:
		return s.energyCharged <= 0
	default:
		panic(fmt.Sprintf("未知的参数类型：%v", s.mode))
	}
}

func (s Stage1) GetReturnedValue() any {
	return nil
}

// 一对一充电Action
type Action1Staged struct {
}

func (a *Action1Staged) MakeStage(charger OneToOneCharger, target OneToOneChargingTarget, param Param1Staged) Stage1 {
	// 目标不存在则直接返回
	if target == nil {
		return Stage1{targetNil: true}
	}

	if p, ok := param.(Param1StagedWithDuration); ok {
		return Stage1{
			mode:         durationMode,
			durationLeft: p.Duration,
		}
	} else if p, ok := param.(Param1StagedWithEnergyLevel); ok {
		return Stage1{
			mode:         energyLevelMode,
			targetEnergy: target.GetEnergy(),
			energyLevel:  min(p.EnergyLevel, target.GetMaxEnergy()),
		}
	} else if p, ok := param.(Param1StagedWithEnergyConsumed); ok {
		return Stage1{
			mode:           energyConsumedMode,
			energyConsumed: p.energyConsumed,
		}
	} else if p, ok := param.(Param1StagedWithEnergyCharged); ok {
		return Stage1{
			mode:          energyChargedMode,
			energyCharged: p.energyCharged,
		}
	} else {
		panic(fmt.Sprintf("未知的参数类型：%T", param))
	}
}

func (a *Action1Staged) ActionStage(charger OneToOneCharger, target OneToOneChargingTarget, param Param1Staged, stage Stage1) Stage1 {
	if !stage.IsLastStage() {
		// 更新状态
		charger.SetState(state.Charging)
		target.SetState(state.Charged)

		// 更新能量
		energyConsumed, energyCharged := charger.ComputeChargingEnergyConsumedAndCharged(target, param.GetTimeunit())
		charger.SetEnergy(charger.GetEnergy() - energyConsumed)
		target.SetEnergy(min(target.GetEnergy()+energyCharged, target.GetMaxEnergy()))

		switch stage.mode {
		case durationMode:
			stage.durationLeft -= param.GetTimeunit()
		case energyLevelMode:
			stage.targetEnergy = target.GetEnergy()
		case energyConsumedMode:
			stage.energyConsumed -= energyConsumed
		case energyChargedMode:
			stage.energyCharged -= energyCharged
		default:
			panic(fmt.Sprintf("未知的参数类型：%v", stage.mode))
		}
	}

	if stage.IsLastStage() {
		// 重置状态
		charger.SetState(state.None)
		if target != nil {
			target.SetState(state.None)
		}
	}
	return stage
}
