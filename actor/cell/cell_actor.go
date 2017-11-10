package cell

import (
	"fmt"
	"gamelib/actor/plugin/logger"
	"gamelib/actor/plugin/timer"
	"gamelib/base/util"
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

	c.addGhost(ctx, entity)
}
func (c *cellActor) remove(ctx actor.Context, entity *imsg.Entity, spos *msg.Vector2) {
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

	c.removeGhost(ctx, entity, spos)
}
func (c *cellActor) broad(ctx actor.Context, m proto.Message) {
	for _, e := range c.entities {
		if e.CellPID != nil {
			continue
		}
		ctx.Request(e.AgentPID, m)
	}

}
func (c *cellActor) switchCell(ctx actor.Context, entity *imsg.Entity, spos *msg.Vector2) {
	c.remove(ctx, entity, spos)
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

	case *msg.CMove:
		m = util.Clone(m).(*msg.CMove)
		entity, ok := c.entities[m.Id]
		if !ok {
			return
		}

		fmt.Println("msg----------", m)
		fmt.Println("upd----------", entity.Data)

		spos := entity.Data.Pos

		entity.Data.Vel = m.Vel
		entity.Data.Pos = m.Pos

		c.broad(ctx, &msg.SUpdate{Data: entity.Data})

		c.onPosChange(ctx, entity, spos)
		c.onVelChange(ctx, entity)
	case *imsg.SwitchCellReq:
		m = util.Clone(m).(*imsg.SwitchCellReq)
		fmt.Println("switch---------", m.Entity.Data)
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

			spos := util.Clone(v.Data.Pos).(*msg.Vector2)
			v.Data.Pos.X = v.Data.Pos.X + v.Data.Vel.X*delta
			v.Data.Pos.Y = v.Data.Pos.Y + v.Data.Vel.Y*delta

			c.onPosChange(ctx, v, spos)
		}

		c.lastEv = now
	case *imsg.AddGhost:
		m = util.Clone(m).(*imsg.AddGhost)

		e := m.Entity
		e.CellPID = ctx.Sender()
		c.entities[e.Data.Id] = e

		c.broad(ctx, &msg.SAdd{Data: []*msg.PlayerData{e.Data}})
		fmt.Println("broad add", e.Data)
	case *imsg.SyncGhost:
		m = util.Clone(m).(*imsg.SyncGhost)

		e := m.Entity
		e.CellPID = ctx.Sender()
		c.entities[e.Data.Id] = e

		c.broad(ctx, &msg.SUpdate{Data: e.Data})
		fmt.Println("broad update", e.Data)
	case *imsg.RemoveGhost:
		m = util.Clone(m).(*imsg.RemoveGhost)

		delete(c.entities, m.Id)
		c.broad(ctx, &msg.SRemove{Id: []int32{m.Id}})
	}
}

func (c *cellActor) onPosChange(ctx actor.Context, entity *imsg.Entity, spos *msg.Vector2) {
	if entity.CellPID != nil {
		return
	}

	if c.cell.OutOfSwitchBorder(entity.Data.Pos) {
		c.switchCell(ctx, entity, spos)
		return
	}

	c.updateGhost(ctx, spos, entity)
}

func (c *cellActor) onVelChange(ctx actor.Context, entity *imsg.Entity) {
	if entity.CellPID != nil {
		return
	}

	c.syncGhost(ctx, entity)
}

func (c *cellActor) addGhost(ctx actor.Context, e *imsg.Entity) {
	for _, c := range c.neighbors {
		if !c.InGhostBorder(e.Data.Pos) {
			continue
		}

		dstCellPID := msg.GetCellPIDByName(c.Name)
		ctx.Request(dstCellPID, &imsg.AddGhost{e})
	}
}
func (c *cellActor) updateGhost(ctx actor.Context, spos *msg.Vector2, e *imsg.Entity) {
	for _, c := range c.neighbors {
		lastIn := c.InGhostBorder(spos)
		nowIn := c.InGhostBorder(e.Data.Pos)
		switch {
		case !lastIn && nowIn:
			dstCellPID := msg.GetCellPIDByName(c.Name)
			ctx.Request(dstCellPID, &imsg.AddGhost{e})
			fmt.Println("-------enter ghost zone", c.Name)
		case lastIn && !nowIn:
			dstCellPID := msg.GetCellPIDByName(c.Name)
			ctx.Request(dstCellPID, &imsg.RemoveGhost{e.Data.Id})
			fmt.Println("-------leave ghost zone", c.Name)
		}
	}
}
func (c *cellActor) syncGhost(ctx actor.Context, entity *imsg.Entity) {
	for _, c := range c.neighbors {
		if c.InGhostBorder(entity.Data.Pos) {
			dstCellPID := msg.GetCellPIDByName(c.Name)
			ctx.Request(dstCellPID, &imsg.SyncGhost{entity})
		}
	}
}
func (c *cellActor) removeGhost(ctx actor.Context, entity *imsg.Entity, spos *msg.Vector2) {
	for _, c := range c.neighbors {
		if c.InGhostBorder(spos) && !c.InGhostBorder(entity.Data.Pos) {
			dstCellPID := msg.GetCellPIDByName(c.Name)
			ctx.Request(dstCellPID, &imsg.RemoveGhost{entity.Data.Id})
		}
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
