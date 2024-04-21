package command

import (
	"fmt"
	"testing"
	"time"

	"github.com/bianxiaojie/rte/ctx"
	"github.com/bianxiaojie/rte/engine"
	"github.com/bianxiaojie/rte/entity/action"
	"github.com/bianxiaojie/rte/utils/ref"
)

// 第一步，自定义负责发送命令的实体，并实现entity.go中定义的实体接口
// 比如基站就可以是一个负责发送命令的实体，用于命令MCV移动或充电
type commandAction1SenderEntity struct {
	id string
}

// --- 以下是调用CommandAction1动作所需实现的CommandSource接口方法
func (e *commandAction1SenderEntity) Id() string {
	return e.id
}

// 第二步，自定义发送命令的行为
// --- 以下是commandAction1SenderEntity的行为，每个时间单位按照方法名后缀的优先级依次执行
func (e *commandAction1SenderEntity) SendCommand_0(context ctx.Context) {
	param := CommandParam1{}
	// 交替向id为receiver的实体发送move和default命令
	if context.Timer().GetTime()%(2*time.Second) == 0 {
		param.Command = "move"
	}
	action.HandleOneTargetAction[*CommandAction1, CommandSource](context.ActionHandler(), e, "receiver", param)
}

// 第三步，自定义负责执行命令的实体类型，并实现entity.go中定义的实体接口
// 比如MCV就可以是一个负责接收命令的实体，用于根据命令执行移动或充电操作
type commandAction1ReceiverEntity struct {
	id      string
	command any
}

// --- 以下是调用CommandAction1动作所需实现的CommandTarget接口方法
func (e *commandAction1ReceiverEntity) Id() string {
	return e.id
}

func (e *commandAction1ReceiverEntity) SetCommand(command any) {
	e.command = command
}

// 第四步，自定义执行命令的行为
// --- 以下是commandAction1ReceiverEntity的行为，每个时间单位按照方法名后缀的优先级依次执行
func (e *commandAction1ReceiverEntity) Run_0(context ctx.Context) {
	// 根据接收命令的不同执行不同的操作
	if e.command == "move" {
		fmt.Println("可以在这里调用action.HandXXXAction来执行MoveAction")
	} else {
		fmt.Println("执行默认的行为")
	}
}

// 第五步，启动仿真
func TestAction0(t *testing.T) {
	// 创建仿真引擎，并将时间间隔设置为1秒，将停止时间设置为3秒
	e := engine.MakeDefaultEngine(time.Second, 3*time.Second)
	// 添加实体行为
	e.EntityManager().AddBehaviorByType(ref.ParseType[*commandAction1SenderEntity]())
	e.EntityManager().AddBehaviorByType(ref.ParseType[*commandAction1ReceiverEntity]())
	// 添加实体
	e.EntityManager().AddEntity(&commandAction1SenderEntity{
		id: "sender",
	})
	e.EntityManager().AddEntity(&commandAction1ReceiverEntity{
		id: "receiver",
	})
	// 启动仿真引擎并等待仿真停止
	e.Start()
	e.WaitStopped()
}
