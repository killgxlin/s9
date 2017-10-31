package cell

import (
	"gamelib/actor/plugin/logger"
	"gamelib/actor/plugin/timer"
	"gamelib/base/util"
	"log"
	"s9/imsg"
	"s9/msg"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/plugin"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/gogo/protobuf/proto"
)

var (
	heartbeatInterval = time.Millisecond * 100
)

// actor --------------------------------------------------
type cellActor struct {
	entities map[int32]*imsg.Entity
	cell     *msg.Cell

	lastEv time.Time
}

func (c *cellActor) add(ctx actor.Context, entity *imsg.Entity) {
	_, ok := c.entities[entity.Data.Id]
	if ok {
		return
	}
	c.broad(ctx, &msg.SAdd{Data: []*msg.PlayerData{entity.Data}})
	sadd := &msg.SAdd{}
	for _, e := range c.entities {
		sadd.Data = append(sadd.Data, e.Data)
	}
	ctx.Request(entity.AgentPID, sadd)
	c.entities[entity.Data.Id] = entity
}
func (c *cellActor) remove(ctx actor.Context, entity *imsg.Entity) {
	_, ok := c.entities[entity.Data.Id]
	if !ok {
		return
	}
	delete(c.entities, entity.Data.Id)
	c.broad(ctx, &msg.SRemove{Id: []int32{entity.Data.Id}})
	srem := &msg.SRemove{}
	for _, e := range c.entities {
		srem.Id = append(srem.Id, e.Data.Id)
	}
	ctx.Request(entity.AgentPID, srem)
}
func (c *cellActor) broad(ctx actor.Context, m proto.Message) {
	for _, e := range c.entities {
		ctx.Request(e.AgentPID, m)
	}
}
func (c *cellActor) onSwitchCell(ctx actor.Context, entity *imsg.Entity) {
	c.remove(ctx, entity)
	ctx.Request(entity.AgentPID, &msg.SLeaveCell{CellName: c.cell.Name})
	dstCellPID := msg.GetCellPID(entity.Data.Pos)
	ctx.Request(dstCellPID, &imsg.SwitchCellReq{Entity: entity})
}

func (c *cellActor) Receive(ctx actor.Context) {
	switch m := ctx.Message().(type) {
	case *actor.Started:
		c.entities = map[int32]*imsg.Entity{}
		c.cell = msg.GetCell(ctx.Self())
		c.lastEv = time.Now()
	case *actor.Stopping, *actor.Restarting:
	case *imsg.ExitSceneReq:
		entity, ok := c.entities[m.Id]
		if !ok {
			return
		}
		pid, e := cluster.Get("scene", "scene")
		util.PanicOnErr(e)
		if pid == nil {
			return
		}
		m.Entity = entity
		ctx.Request(pid, m)
		c.remove(ctx, entity)

	case *msg.CUpdate:
		entity, ok := c.entities[m.Data.Id]
		if !ok {
			return
		}

		*entity.Data = *m.Data
		c.broad(ctx, &msg.SUpdate{Data: m.Data})
	case *imsg.SwitchCellReq:
		ctx.Request(
			m.Entity.AgentPID,
			&msg.SEnterCell{
				Self: m.Entity.Data,
				Cell: c.cell,
			})
		c.add(ctx, m.Entity)
	case timer.TimerEvent:
		now := time.Now()
		delta := float32(now.Sub(c.lastEv).Seconds())
		leaves := []*imsg.Entity{}
		for _, v := range c.entities {
			lpos := proto.Clone(v.Data)
			v.Data.Pos.X = v.Data.Pos.X + v.Data.Vel.X*delta
			v.Data.Pos.Y = v.Data.Pos.Y + v.Data.Vel.Y*delta
			if c.cell.OutOfSwitchBorder(v.Data.Pos) {
				log.Println(lpos, v.Data)
				leaves = append(leaves, v)
			}
		}
		for _, v := range leaves {
			c.onSwitchCell(ctx, v)
		}

		c.lastEv = time.Now()
	}
}

func init() {
	logger.Filter(
		timer.TimerEvent{},
		&actor.Started{})
	remote.Register("cell", actor.FromProducer(func() actor.Actor {
		return &cellActor{}
	}).WithMiddleware(
		logger.MsgLogger,
		plugin.Use(timer.NewTimer(heartbeatInterval)),
	))
}
