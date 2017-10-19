package cellman

import (
	"gamelib/actor/plugin/timer"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/plugin"
	"github.com/AsynkronIT/protoactor-go/remote"
)

var (
	heartbeatInterval = time.Millisecond * 1000
)

// actor --------------------------------------------------
type cellActor struct {
}

func (c *cellActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Started:
	case *actor.Stopping, *actor.Restarting:
	case timer.TimerEvent:
	}
}

func init() {
	remote.Register("cell", actor.FromProducer(func() actor.Actor {
		return &cellActor{}
	}).WithMiddleware(
		//logger.MsgLogger,
		plugin.Use(timer.NewTimer(heartbeatInterval)),
	))
}
