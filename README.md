# WRSN组件库

## 简介

该库基于实时仿真引擎[rte](https://github.com/bianxiaojie/rte)，提供了开箱即用的WRSN **Action**组件，这些动作包括感知耗能、移动、充电等。目前可用的组件数目不多，如果有新的需求，可以提issues，也欢迎贡献新的**Action**组件。

## 快速开始

首先，在项目目录执行go get命令获取本依赖库：

```shell
# go get github.com/bianxiaojie/wrsn@v0.1.0
$ go get github.com/bianxiaojie/wrsn@最新版本
```

然后，找到需要使用的动作组件，比如SimpleSenseAction0，并实现SimpleSenseAction0同文件夹下entity.go中定义的实体接口方法，比如Sensible。以下是一个示例：

```go
// entity.go
type Sensible interface {
	entity.Entity
	SetState(state.WRSNEntityState)                          // 在感知中实体状态更新为Sensing，感知后更新为None
	SetEnergy(value.Energy)                                  // 更新实体能量
	GetEnergy() value.Energy                                 // 获取实体能量
	ComputeSensingEnergyConsumed(time.Duration) value.Energy // 计算给定时间内实体的感知耗费的能量
}

// action0_test.go
// 自定义实体
type simpleSenseAction0Entity struct {
	id          string
	shouldSense bool
	state       state.WRSNEntityState
	energy      value.Energy
}

// 以下是调用SimpleSenseAction0动作所需实现的Sensible接口方法
func (e *simpleSenseAction0Entity) Id() string {
	return e.id
}

func (e *simpleSenseAction0Entity) SetState(state state.WRSNEntityState) {
	e.state = state
}

func (e *simpleSenseAction0Entity) GetEnergy() value.Energy {
	return e.energy
}

func (e *simpleSenseAction0Entity) SetEnergy(energy value.Energy) {
	e.energy = energy
}

func (e *simpleSenseAction0Entity) ComputeSensingEnergyConsumed(duration time.Duration) value.Energy {
	return value.Joule * value.Energy(duration/time.Second)
}
```

然后，就可以在自定义**Behavior**函数中调用**Action**组件：

```go
// Run Behavior
func (e *simpleSenseAction0Entity) Run_0(context ctx.Context) {
    // 创建Action参数
	param := SimpleSenseParam0{
		Timeunit: context.Timer().GetTimeunit(),
	}
    fmt.Printf("感知前的能量：%v\n", e.energy)
    // 调用感知Action
    action.HandleNoneTargetAction[*SimpleSenseAction0, Sensible](context.ActionHandler(), e, param)
    fmt.Printf("感知后的能量：%v\n", e.energy)
}
```

最后，启动仿真引擎，就可以验证执行结果了。当感知Action被调用时，组件库会处理感知的逻辑，包括状态的变更、能量的消耗等。

```go
// 创建仿真引擎，并将时间间隔和停止时间都设置为1秒
e := engine.MakeDefaultEngine(time.Second, time.Second)
// 添加自定义实体Behavior
e.EntityManager().AddBehaviorByType(ref.ParseType[*simpleSenseAction0Entity]())
// 添加自定义实体
e.EntityManager().AddEntity(&simpleSenseAction0Entity{
    id:          "id",
    shouldSense: true,
    state:       state.None,
    energy:      value.Joule,
})
// 启动仿真引擎并等待仿真停止
e.Start()
e.WaitStopped()
```

执行结果如下：

```go
// 感知前的能量：1.0J
// 感知后的能量：0.0J
```

## 概念

在上面的例子中，提到两个概念，**Behavior**函数和**Action**组件。

### Behavior函数

形如这样的函数会被识别为**Behavior**函数：`func (e *Entity) XXX_1(context ctx.Context)`。特点如下：

- 函数名由一个下划线分隔，下划线前面是函数名，后面是优先级值。仿真引擎会将所有注册的**Behavior**函数按照该优先级值排序，值越小，优先级越高，在同一轮中会被优先调用。
- 参数包含一个ctx.Context。借助该上下文，我们可以调用动作，获取仿真时间，获取网络中的实体等。

在仿真过程中，所有**Behavior**函数会按照顺序循环执行，直到达到仿真结束。

比如，用Sense_0和Charge_1两个**Behavior**函数，那么仿真引擎将按照Sense_0 -> Charge_1 -> Sense_0 -> Charge_1 ...的顺序执行，每执行一次循环，时间都会增加一个单位。

### Action组件

**Action**组件就是可复用的实体动作，比如充电、移动。我们可以通过多种方式对**Action**组件进行分类。

- 根据Action执行目标的多少，可以分为：NoneTargetAction、OneTargetAction和MultipleTargetAction，分别表示无目标动作、单目标动作和多目标动作，比如更换电池是无目标动作（动作源不算动作目标），一对一充电是单目标动作，一对多充电是多目标动作。
- 根据Action执行时间的长短，可以分为：Action和StagedAction，分别表示短时动作和长时动作。前者表示该动作只需要一个仿真单位时间就可以完成，比如更换电池就是短时动作。后者表示该动作需要多个仿真时间单位才能完成，比如移动、充电就是长时动作，调用这种类型的动作后，函数会在若干个时间单位后返回。

根据上述分类，在该组件库的action包中，也会按照该分类方式对Action组件所在的go文件命名，命名规则是[staged]action0，前面的staged用于区分短时和长时动作，后面的数字表示执行目标的数量。

## 目录说明

action：该目录定义了各种**Action**组件，每个包下包含三种类型的go文件，action、entity和test。action中定义了动作和动作参数，entity中定义了调用该动作的实体需要实现的接口，test中包含使用该动作的例子，包括如何自定义实体并调用动作组件。

common：组件库的公用代码。

examples：包含一些更复杂的使用**Action**组件的仿真示例。