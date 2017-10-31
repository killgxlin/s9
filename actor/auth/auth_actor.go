package auth

import (
	"s9/imsg"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

var (
	heartbeatInterval = time.Millisecond * 100
)

// actor --------------------------------------------------
type authActor struct {
	idGen int32

	accounts map[string]int32
}

func (c *authActor) Receive(ctx actor.Context) {
	switch m := ctx.Message().(type) {
	case *actor.Started:
		c.accounts = map[string]int32{}
	case *actor.Stopping, *actor.Restarting:
	case *imsg.AuthReq:
		id, ok := c.accounts[m.Account]
		if !ok {
			c.idGen++
			id = c.idGen
			c.accounts[m.Account] = id
		}
		ctx.Respond(&imsg.AuthRep{Id: id})
	}
}

func init() {
	remote.Register("auth", actor.FromProducer(func() actor.Actor {
		return &authActor{}
	}).WithMiddleware(
	//logger.MsgLogger,
	))
}
