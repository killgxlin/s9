package gate

import (
	pnet "gamelib/actor/plugin/net"
	"gamelib/base/util"
	"s9/msg"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
)

// actor --------------------------------------------------
type connActor struct {
	id int32
}

func (c *connActor) Receive(ctx actor.Context) {
	switch m := ctx.Message().(type) {
	case *actor.Started:
		pid, e := cluster.Get("cell", "cell")
		util.PanicOnErr(e)
		ctx.Request(pid, &msg.Connected{})
	case *actor.Stopping, *actor.Restarting:
		pid, e := cluster.Get("cell", "cell")
		util.PanicOnErr(e)
		ctx.Tell(pid, &msg.Disconnect{Id: c.id})
	case *pnet.ConnectionEvent:
		ctx.Self().Stop()
	case *msg.SEnter:
		c.id = m.Self.Id
		pnet.SendMsg(ctx, m)
	case *msg.SAdd:
		pnet.SendMsg(ctx, m)
	case *msg.SRemove:
		pnet.SendMsg(ctx, m)
	case *msg.SMove:
		pnet.SendMsg(ctx, m)
	case *msg.CMove:
		pid, e := cluster.Get("cell", "cell")
		util.PanicOnErr(e)
		ctx.Tell(pid, m)
	}
}
