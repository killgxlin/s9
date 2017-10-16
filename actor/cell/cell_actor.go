package cell

import (
	"log"
	"s7/share/middleware/msglogger"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

type contextData struct {
	PlayerData
	Pid *actor.PID
}

// actor --------------------------------------------------
type cellActor struct {
	context map[int32]*contextData
	idGen   int32
}

func (c *cellActor) Receive(ctx actor.Context) {
	switch m := ctx.Message().(type) {
	case *actor.Started:
		c.context = map[int32]*contextData{}
	case *actor.Stopping, *actor.Restarting:
	case *Connected:
		c.idGen++
		p := &contextData{
			PlayerData: PlayerData{
				Id:  c.idGen,
				Pos: &Vector3{0, 0, 0},
				Vel: &Vector3{0, 0, 0},
			},
			Pid: ctx.Sender(),
		}
		add := &SAdd{Data: &p.PlayerData}
		ent := &SEnter{Self: &p.PlayerData}
		for _, c := range c.context {
			ent.Other = append(ent.Other, &c.PlayerData)
			c.Pid.Tell(add)
		}
		ctx.Respond(ent)
		c.context[p.Id] = p
	case *Disconnect:
		context, ok := c.context[m.Id]
		if !ok {
			return
		}
		delete(c.context, context.Id)
		rem := &SRemove{Id: context.Id}
		for _, c := range c.context {
			c.Pid.Tell(rem)
		}
		log.Println(c)
	case *CMove:
		context, ok := c.context[m.Data.Id]
		if !ok {
			return
		}

		mov := &SMove{Data: m.Data}

		context.PlayerData = *m.Data
		for _, c := range c.context {
			c.Pid.Tell(mov)
		}
	}
}

func init() {
	remote.Register("cell", actor.FromProducer(func() actor.Actor {
		return &cellActor{}
	}).WithMiddleware(msglogger.MsgLogger))
}
