package command

import (
	"fmt"
	"testing"
	"time"

	"github.com/bianxiaojie/rte/context"
	"github.com/bianxiaojie/rte/engine"
	"github.com/bianxiaojie/rte/entity/action"
	"github.com/bianxiaojie/rte/utils/ref"
)

// 第一步，自定义负责发送命令的实体，并实现entity.go中定义的实体接口
// 比如基站就可以是一个负责发送命令的实体，用于命令MCV移动或充电
type action1SenderEntity struct {
	id string
}

func (e *action1SenderEntity) Id() string {
	return e.id
}

// 第二步，自定义发送命令的行为
func (e *action1SenderEntity) SendCommand_0(ctx context.Context) {
	param := Param1{}
	// 交替向id为receiver的实体发送move和default命令
	if ctx.Timer().GetTime()%(2*time.Second) == 0 {
		param.Command = "move"
	}
	action.HandleOneTargetAction[*Action1, CommandSource](ctx.ActionHandler(), e, "receiver", param)
}

// 第三步，自定义负责执行命令的实体类型，并实现entity.go中定义的实体接口
// 比如MCV就可以是一个负责接收命令的实体，用于根据命令执行移动或充电操作
type action1ReceiverEntity struct {
	id      string
	command any
}

func (e *action1ReceiverEntity) Id() string {
	return e.id
}

func (e *action1ReceiverEntity) SetCommand(command any) {
	e.command = command
}

// 第四步，自定义执行命令的行为
func (e *action1ReceiverEntity) Run_0(ctx context.Context) {
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
	e.EntityManager().AddBehaviorByType(ref.ParseType[*action1SenderEntity]())
	e.EntityManager().AddBehaviorByType(ref.ParseType[*action1ReceiverEntity]())
	// 添加实体
	e.EntityManager().AddEntity(&action1SenderEntity{
		id: "sender",
	})
	e.EntityManager().AddEntity(&action1ReceiverEntity{
		id: "receiver",
	})
	// 启动仿真引擎并等待仿真停止
	e.Start()
	e.WaitStopped()
}
