package command

type Param1 struct {
	Command any
}

type Action1 struct {
}

func (a *Action1) Action(source CommandSource, target CommandTarget, param Param1) any {
	// 目标不存在则直接返回
	if target == nil {
		return nil
	}

	// 发送命令
	target.SetCommand(param.Command)

	return nil
}
