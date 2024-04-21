package onetoonecharge

import (
	"fmt"
	"time"

	"github.com/bianxiaojie/wrsn/common/state"
	"github.com/bianxiaojie/wrsn/common/value"
)

type OneToOneChargeStagedParam1 interface {
	GetTimeunit() time.Duration // 仿真的时间单位
}

// 充电参数，在需要指定充电时长时使用
type OneToOneChargeStagedParam1WithDuration struct {
	Timeunit time.Duration // 仿真的时间单位
	Duration time.Duration // 充电总时长
}

func (p OneToOneChargeStagedParam1WithDuration) GetTimeunit() time.Duration {
	return p.Timeunit
}

// 充电参数，在需要指定充电后充电方剩余能量时使用
type OneToOneChargeStagedParam1WithEnergyLevel struct {
	Timeunit    time.Duration // 仿真的时间单位
	EnergyLevel value.Energy  // 充电后充电方剩余能量，比如为最大能量
}

func (p OneToOneChargeStagedParam1WithEnergyLevel) GetTimeunit() time.Duration {
	return p.Timeunit
}

// 充电参数，在需要指定充电方消耗能量时使用
type OneToOneChargeStagedParam1WithEnergyConsumed struct {
	Timeunit       time.Duration // 仿真的时间单位
	energyConsumed value.Energy  // 充电方需要消耗的能量
}

func (p OneToOneChargeStagedParam1WithEnergyConsumed) GetTimeunit() time.Duration {
	return p.Timeunit
}

// 充电参数，在需要指定被充电方收到能量时使用
type OneToOneChargeStagedParam1WithEnergyCharged struct {
	Timeunit      time.Duration // 仿真的时间单位
	energyCharged value.Energy  // 被充电方需要收到的能量
}

func (p OneToOneChargeStagedParam1WithEnergyCharged) GetTimeunit() time.Duration {
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

type OneToOneChargeStage1 struct {
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

func (s OneToOneChargeStage1) IsLastStage() bool {
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

func (s OneToOneChargeStage1) GetReturnedValue() any {
	return nil
}

// 一对一充电Action
type OneToOneChargeStagedAction1 struct {
}

func (a *OneToOneChargeStagedAction1) MakeStage(charger OneToOneCharger, target OneToOneChargingTarget, param OneToOneChargeStagedParam1) OneToOneChargeStage1 {
	// 目标不存在则直接返回
	if target == nil {
		return OneToOneChargeStage1{targetNil: true}
	}

	if p, ok := param.(OneToOneChargeStagedParam1WithDuration); ok {
		return OneToOneChargeStage1{
			mode:         durationMode,
			durationLeft: p.Duration,
		}
	} else if p, ok := param.(OneToOneChargeStagedParam1WithEnergyLevel); ok {
		return OneToOneChargeStage1{
			mode:         energyLevelMode,
			targetEnergy: target.GetEnergy(),
			energyLevel:  min(p.EnergyLevel, target.GetMaxEnergy()),
		}
	} else if p, ok := param.(OneToOneChargeStagedParam1WithEnergyConsumed); ok {
		return OneToOneChargeStage1{
			mode:           energyConsumedMode,
			energyConsumed: p.energyConsumed,
		}
	} else if p, ok := param.(OneToOneChargeStagedParam1WithEnergyCharged); ok {
		return OneToOneChargeStage1{
			mode:          energyChargedMode,
			energyCharged: p.energyCharged,
		}
	} else {
		panic(fmt.Sprintf("未知的参数类型：%T", param))
	}
}

func (a *OneToOneChargeStagedAction1) ActionStage(charger OneToOneCharger, target OneToOneChargingTarget, param OneToOneChargeStagedParam1, stage OneToOneChargeStage1) OneToOneChargeStage1 {
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
