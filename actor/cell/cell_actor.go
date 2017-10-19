package cell

import (
	"gamelib/actor/plugin/timer"
	"log"
	"s9/msg"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/plugin"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/gogo/protobuf/proto"
)

var (
	heartbeatInterval = time.Millisecond * 100
)

type entityData struct {
	msg.PlayerData
	Pid *actor.PID
}

// actor --------------------------------------------------
type cellActor struct {
	context map[int32]*entityData
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
		c.context = map[int32]*entityData{}
	case *actor.Stopping, *actor.Restarting:
	case *msg.Connected:
		c.idGen++

		p := &entityData{
			PlayerData: msg.PlayerData{
				Id:  c.idGen,
				Pos: &msg.Vector3{0, 0, 0},
				Vel: &msg.Vector3{0, 0, 0},
			},
			Pid: ctx.Sender(),
		}

		c.broad(&msg.SAdd{Data: &p.PlayerData})

		ent := &msg.SEnter{Self: &p.PlayerData}
		for _, c := range c.context {
			ent.Other = append(ent.Other, &c.PlayerData)
		}
		ctx.Respond(ent)
		c.context[p.Id] = p
	case *msg.Disconnect:
		context, ok := c.context[m.Id]
		if !ok {
			return
		}
		delete(c.context, context.Id)
		c.broad(&msg.SRemove{Id: context.Id})
	case *msg.CMove:
		context, ok := c.context[m.Data.Id]
		if !ok {
			return
		}

		log.Println(context.PlayerData)
		context.PlayerData = *m.Data
		log.Println(context.PlayerData)
		c.broad(&msg.SMove{Data: m.Data})
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
		plugin.Use(timer.NewTimer(heartbeatInterval)),
	))
}
