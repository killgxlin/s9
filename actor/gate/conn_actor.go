package gate

import (
	"s7/share/middleware/mnet"
	"s7/share/util"
	"s9/actor/cell"

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
		ctx.Request(pid, &cell.Connected{})
	case *actor.Stopping, *actor.Restarting:
		pid, e := cluster.Get("cell", "cell")
		util.PanicOnErr(e)
		ctx.Tell(pid, &cell.Disconnect{Id: c.id})
	case *mnet.ConnectionEvent:
		ctx.Self().Stop()
	case *cell.SEnter:
		c.id = m.Self.Id
		mnet.SendMsg(ctx, m)
	case *cell.SAdd:
		mnet.SendMsg(ctx, m)
	case *cell.SRemove:
		mnet.SendMsg(ctx, m)
	case *cell.SMove:
		mnet.SendMsg(ctx, m)
	case *cell.CMove:
		pid, e := cluster.Get("cell", "cell")
		util.PanicOnErr(e)
		ctx.Tell(pid, m)
	}
}
