package gate

import (
	"gamelib/actor/plugin/logger"
	pnet "gamelib/actor/plugin/net"
	"gamelib/base/net/util"
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/plugin"
)

type gateActor struct {
}

func (g *gateActor) Receive(ctx actor.Context) {
	switch m := ctx.Message().(type) {
	case *actor.Started:
	case *pnet.AcceptorEvent:
		if m.E != nil {
			ctx.Self().Stop()
			return
		}
		log.Printf("gate port:%v", m.GetPort())
		if m.C != nil {
			ctx.SpawnPrefix(startAgent(m.C), "agent")
		}
	}
}

func Start(start, end int) {
	logger.Filter(&pnet.ConnOp{})

	addr, e := util.FindLanAddr("tcp", start, end)
	if e != nil {
		log.Panic(e)
	}

	prop := actor.FromProducer(func() actor.Actor {
		return &gateActor{}
	}).WithMiddleware(
		logger.MsgLogger,
		plugin.Use(pnet.NewAcceptor(addr, 100, 100)),
	)
	_, e = actor.SpawnNamed(prop, "gate")
	if e != nil {
		log.Panic(e)
	}
}
