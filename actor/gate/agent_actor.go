package gate

import (
	pnet "gamelib/actor/plugin/net"
	"gamelib/base/util"
	"s9/imsg"
	"s9/msg"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
)

// actor --------------------------------------------------
type agentActor struct {
	id      int32
	cellPID *actor.PID
}

func (aa *agentActor) Receive(ctx actor.Context) {
	switch m := ctx.Message().(type) {
	case *actor.Started:
	case *actor.Stopping, *actor.Restarting:
		if aa.cellPID == nil {
			return
		}
		ctx.Tell(aa.cellPID, &imsg.ExitSceneReq{Id: aa.id})
	case *pnet.ConnectionEvent:
		ctx.Self().Stop()
	case *msg.CLogin:
		pid, e := cluster.Get("auth", "auth")
		util.PanicOnErr(e)
		rep, e := pid.RequestFuture(&imsg.AuthReq{Account: m.Account}, time.Second).Result()
		if e != nil {
			return
		}
		arep := rep.(*imsg.AuthRep)
		pid, e = cluster.Get("scene", "scene")
		util.PanicOnErr(e)
		ctx.Request(pid, &imsg.EnterSceneReq{Id: arep.Id})
	case *msg.CUpdate:
		if aa.cellPID != nil {
			ctx.Tell(aa.cellPID, m)
		}
	case *msg.SEnterCell:
		aa.id = m.Self.Id
		pnet.SendMsg(ctx, m)
		aa.cellPID = ctx.Sender()
	case *msg.SLeaveCell:
		pnet.SendMsg(ctx, m)
		aa.cellPID = nil
	case *msg.SAdd:
		pnet.SendMsg(ctx, m)
	case *msg.SRemove:
		pnet.SendMsg(ctx, m)
	case *msg.SUpdate:
		pnet.SendMsg(ctx, m)
	}
}
