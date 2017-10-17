package cell

import (
	"gamelib/actor/plugin/timer"
	"log"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/plugin"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/gogo/protobuf/proto"
)

var (
	delta = time.Millisecond * 100
)

type contextData struct {
	PlayerData
	Pid *actor.PID
}

// actor --------------------------------------------------
type cellActor struct {
	context map[int32]*contextData
	idGen   int32

	lastEv time.Time
}

func (c *cellActor) broad(m proto.Message) {
	for _, c := range c.context {
		c.Pid.Tell(m)
	}
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

		c.broad(&SAdd{Data: &p.PlayerData})

		ent := &SEnter{Self: &p.PlayerData}
		for _, c := range c.context {
			ent.Other = append(ent.Other, &c.PlayerData)
		}
		ctx.Respond(ent)
		c.context[p.Id] = p
	case *Disconnect:
		context, ok := c.context[m.Id]
		if !ok {
			return
		}
		delete(c.context, context.Id)
		c.broad(&SRemove{Id: context.Id})
	case *CMove:
		context, ok := c.context[m.Data.Id]
		if !ok {
			return
		}

		log.Println(context.PlayerData)
		context.PlayerData = *m.Data
		log.Println(context.PlayerData)
		c.broad(&SMove{Data: m.Data})
	case timer.TimerEvent:
		now := time.Now()
		delta := float32(now.Sub(c.lastEv).Seconds())
		for _, v := range c.context {
			v.Pos.X = v.Pos.X + v.Vel.X*delta
			v.Pos.Y = v.Pos.Y + v.Vel.Y*delta
			v.Pos.Z = v.Pos.Z + v.Vel.Z*delta
		}

		c.lastEv = time.Now()
	}
}

func init() {
	remote.Register("cell", actor.FromProducer(func() actor.Actor {
		return &cellActor{}
	}).WithMiddleware(
		//logger.MsgLogger,
		plugin.Use(timer.NewTimer(delta)),
	))
}
