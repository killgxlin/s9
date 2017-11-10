package cell

import (
	"fmt"
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
	entities  map[int32]*imsg.Entity
	cell      *msg.Cell
	neighbors []*msg.Cell

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

	c.updateGhost(ctx, entity)
}
func (c *cellActor) remove(ctx actor.Context, entity *imsg.Entity, lpos *msg.Vector2) {
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

	c.removeGhost(ctx, entity.Data.Id, lpos)
}
func (c *cellActor) broad(ctx actor.Context, m proto.Message) {
	for _, e := range c.entities {
		if e.CellPID != nil {
			continue
		}
		ctx.Request(e.AgentPID, m)
	}

}
func (c *cellActor) onSwitchCell(ctx actor.Context, entity *imsg.Entity, lpos *msg.Vector2) {
	c.remove(ctx, entity, lpos)
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
		c.neighbors = msg.GenNeighbours(c.cell.Name)

	case *actor.Stopping, *actor.Restarting:
	case *imsg.ExitSceneReq:
		m = util.Clone(m).(*imsg.ExitSceneReq)
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
		c.remove(ctx, entity, entity.Data.Pos)

	case *msg.CUpdate:
		m = util.Clone(m).(*msg.CUpdate)
		entity, ok := c.entities[m.Data.Id]
		if !ok {
			return
		}

		fmt.Println("msg----------", m.Data)
		fmt.Println("upd----------", entity.Data)

		entity.Data.Vel = m.Data.Vel
		c.broad(ctx, &msg.SUpdate{Data: m.Data})
		c.onRelocate(ctx, entity, m.Data.Pos)
	case *imsg.SwitchCellReq:
		log.Println("switch---------", m.Entity.Data)
		m = util.Clone(m).(*imsg.SwitchCellReq)
		ctx.Request(
			m.Entity.AgentPID,
			&msg.SEnterCell{
				Self:     m.Entity.Data,
				Cell:     c.cell,
				Neighbor: c.neighbors,
			})
		c.add(ctx, m.Entity)
	case timer.TimerEvent:
		now := time.Now()
		delta := float32(now.Sub(c.lastEv).Seconds())
		updates := []*imsg.Entity{}
		for _, v := range c.entities {
			updates = append(updates, v)
		}

		for _, v := range updates {
			if v.Data.Vel.IsZero() {
				continue
			}
			npos := util.Clone(v.Data.Pos).(*msg.Vector2)
			npos.X = npos.X + v.Data.Vel.X*delta
			npos.Y = npos.Y + v.Data.Vel.Y*delta
			c.onRelocate(ctx, v, npos)
		}

		c.lastEv = now
	case *imsg.SyncGhost:
		m = util.Clone(m).(*imsg.SyncGhost)

		e := m.Entity
		e.CellPID = ctx.Sender()

		_, ok := c.entities[e.Data.Id]
		c.entities[e.Data.Id] = e

		if !ok {
			c.broad(ctx, &msg.SAdd{Data: []*msg.PlayerData{e.Data}})
			log.Println("broad add", e.Data)
			return
		}

		c.broad(ctx, &msg.SUpdate{Data: e.Data})
		log.Println("broad update", e.Data)
	case *imsg.RemoveGhost:
		m = util.Clone(m).(*imsg.RemoveGhost)

		delete(c.entities, m.Id)
		c.broad(ctx, &msg.SRemove{Id: []int32{m.Id}})
	}
}

func (c *cellActor) onRelocate(ctx actor.Context, e *imsg.Entity, npos *msg.Vector2) {
	if e.CellPID != nil {
		return
	}

	npos, e.Data.Pos = e.Data.Pos, npos

	if c.cell.OutOfSwitchBorder(e.Data.Pos) {
		c.onSwitchCell(ctx, e, npos)
		return
	}

	c.syncGhost(ctx, npos, e)
}

func (c *cellActor) syncGhost(ctx actor.Context, lpos *msg.Vector2, e *imsg.Entity) {
	for _, c := range c.neighbors {
		lastIn := c.InGhostBorder(lpos)
		nowIn := c.InGhostBorder(e.Data.Pos)
		switch {
		case /*!lastIn &&*/ nowIn:
			dstCellPID := msg.GetCellPIDByName(c.Name)
			ctx.Request(dstCellPID, &imsg.SyncGhost{e})
			fmt.Println("-------enter ghost zone", c.Name)
		case lastIn && !nowIn:
			dstCellPID := msg.GetCellPIDByName(c.Name)
			ctx.Request(dstCellPID, &imsg.RemoveGhost{e.Data.Id})
			fmt.Println("-------leave ghost zone", c.Name)
		}
	}
}
func (c *cellActor) updateGhost(ctx actor.Context, e *imsg.Entity) {
	for _, c := range c.neighbors {
		if !c.InGhostBorder(e.Data.Pos) {
			continue
		}

		dstCellPID := msg.GetCellPIDByName(c.Name)
		ctx.Request(dstCellPID, &imsg.SyncGhost{e})
	}
}
func (c *cellActor) removeGhost(ctx actor.Context, id int32, lpos *msg.Vector2) {
	for _, c := range c.neighbors {
		if !c.InGhostBorder(lpos) {
			continue
		}
		dstCellPID := msg.GetCellPIDByName(c.Name)
		ctx.Request(dstCellPID, &imsg.RemoveGhost{id})
	}
}

func init() {
	logger.Filter(
		timer.TimerEvent{},
		&actor.Started{},
	)
	remote.Register("cell", actor.FromProducer(func() actor.Actor {
		return &cellActor{}
	}).WithMiddleware(
		//logger.MsgLogger,
		plugin.Use(timer.NewTimer(heartbeatInterval)),
	))
}
