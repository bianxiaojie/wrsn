package main

import (
	"fmt"
	"time"

	"github.com/bianxiaojie/rte/engine"
	"github.com/bianxiaojie/rte/utils/ref"
	"github.com/bianxiaojie/wrsn/common/value"
	"github.com/bianxiaojie/wrsn/examples/fifs/entity"
)

func main() {
	// 获取节点的type
	sensorType := ref.ParseType[*entity.Sensor]()
	mcvType := ref.ParseType[*entity.MCV]()
	bsType := ref.ParseType[*entity.BS]()
	// 创建执行引擎
	e := engine.MakeDefaultEngine(time.Second, time.Hour)
	// 注册Sensor的行为
	e.EntityManager().AddBehaviorByType(sensorType)
	e.EntityManager().AddBehaviorByType(mcvType)
	e.EntityManager().AddBehaviorByType(bsType)
	// 将节点添加到网络中
	for i := 0; i < 7; i++ {
		e.EntityManager().AddEntity(entity.MakeSensor(
			fmt.Sprintf("Sensor%d", i),
			value.MakePosition(100*value.Length(i-3)*value.Meter, 100*value.Length(i-3)*value.Meter),
			100*value.Joule,
			0.1*value.Joule,
			0.7,
		))
	}
	e.EntityManager().AddEntity(entity.MakeMCV(
		"MCV",
		value.MakePosition(0, 0),
		3*value.Meter,
		value.Joule,
		1000*value.Joule,
		value.Joule,
		1.0,
	))
	e.EntityManager().AddEntity(entity.MakeBS("BS", value.MakePosition(0, 0)))

	// 启动网络
	e.Start()
	// 等待网络执行完成
	e.WaitStopped()

	fmt.Printf("网络结束时的节点数：%d\n", len(e.EntityManager().GetEntitiesByType(sensorType)))
}
