package command

import "github.com/bianxiaojie/rte/entity"

type CommandSource interface {
	entity.Entity
}

type CommandTarget interface {
	entity.Entity
	SetCommand(any)
}
