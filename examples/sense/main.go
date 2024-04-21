package main

import (
	"fmt"
	"time"

	"github.com/bianxiaojie/rte/engine"
	"github.com/bianxiaojie/rte/utils/ref"
	"github.com/bianxiaojie/wrsn/common/value"
	"github.com/bianxiaojie/wrsn/examples/sense/entity"
)

func main() {
	// 获取节点的type
	sensorType := ref.ParseType[*entity.Sensor]()
	// 创建执行引擎
	e := engine.MakeDefaultEngine(time.Second, time.Hour)
	// 注册Sensor的行为
	e.EntityManager().AddBehaviorByType(sensorType)
	// 将节点添加到网络中
	for i := 0; i < 100; i++ {
		sensor := entity.MakeSensor(fmt.Sprintf("节点%d", i), 100*value.Joule)
		e.EntityManager().AddEntity(sensor)
	}
	// 启动网络
	e.Start()
	// 等待网络执行完成
	e.WaitStopped()

	fmt.Printf("网络结束时的节点数：%d\n", len(e.EntityManager().GetEntitiesByType(sensorType)))
}
