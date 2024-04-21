package command

type CommandParam1 struct {
	Command any
}

type CommandAction1 struct {
}

func (a *CommandAction1) Action(source CommandSource, target CommandTarget, param CommandParam1) any {
	// 目标不存在则直接返回
	if target == nil {
		return nil
	}

	// 发送命令
	target.SetCommand(param.Command)

	return nil
}
