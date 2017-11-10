package scene

import (
	"gamelib/actor/plugin/timer"
	"gamelib/base/util"
	"s9/imsg"
	"s9/msg"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/plugin"
	"github.com/AsynkronIT/protoactor-go/remote"
)

var (
	heartbeatInterval = time.Millisecond * 1000
)

// actor --------------------------------------------------
type sceneActor struct {
	playerDatas map[int32]*msg.PlayerData
}

func (s *sceneActor) Receive(ctx actor.Context) {

	switch m := ctx.Message().(type) {
	case *actor.Started:
		s.playerDatas = map[int32]*msg.PlayerData{}
	case *actor.Stopping, *actor.Restarting:
	case timer.TimerEvent:
	case *imsg.EnterSceneReq:
		m = util.Clone(m).(*imsg.EnterSceneReq)
		data, ok := s.playerDatas[m.Id]
		if !ok {
			data = &msg.PlayerData{
				Id:  m.Id,
				Pos: &msg.Vector2{0, 0},
				Vel: &msg.Vector2{0, 0},
			}
			s.playerDatas[data.Id] = data
		}
		pid := msg.GetCellPID(data.Pos)
		ctx.Tell(
			pid,
			&imsg.SwitchCellReq{
				Entity: &imsg.Entity{
					Data:     data,
					AgentPID: ctx.Sender(),
				},
			})
	case *imsg.ExitSceneReq:
		m = util.Clone(m).(*imsg.ExitSceneReq)
		s.playerDatas[m.Id] = m.Entity.Data
	}
}

func init() {
	remote.Register("scene", actor.FromProducer(func() actor.Actor {
		return &sceneActor{}
	}).WithMiddleware(
		//logger.MsgLogger,
		plugin.Use(timer.NewTimer(heartbeatInterval)),
	))
}
