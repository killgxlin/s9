package imsg

import (
	"fmt"
	"gamelib/base/util"
	"math"
	"s9/msg"
	"strconv"
	"strings"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
)

const (
	CellSize   = 2
	SwitchDist = 0.5
)

func GetCellName(pos *msg.Vector2) string {
	xIdx := int(pos.X / CellSize)
	yIdx := int(pos.Y / CellSize)
	return fmt.Sprintf("cell_%d_%d", xIdx, yIdx)
}

func GetCellPID(pos *msg.Vector2) *actor.PID {
	name := GetCellName(pos)
	pid, e := cluster.Get(name, "cell")
	util.PanicOnErr(e)
	return pid
}

func getRangeByIdx(idx int) (float32, float32) {
	if idx == 0 {
		return -CellSize, CellSize
	}
	var minx, maxx float32
	sign := math.Abs(float64(idx)) / float64(idx)
	if sign > 0 {
		minx = float32(math.Abs(float64(idx)) * CellSize)
		maxx = float32((math.Abs(float64(idx)) + 1) * CellSize)
	} else {
		minx = -float32((math.Abs(float64(idx)) + 1) * CellSize)
		maxx = -float32(math.Abs(float64(idx)) * CellSize)
	}
	return minx, maxx
}

func GetCellByName(name string) *Cell {
	f := strings.Split(name, "_")

	xIdx, e := strconv.Atoi(f[1])
	util.PanicOnErr(e)
	yIdx, e := strconv.Atoi(f[2])
	util.PanicOnErr(e)

	border := &AABB{}
	border.Minx, border.Maxx = getRangeByIdx(xIdx)
	border.Miny, border.Maxy = getRangeByIdx(yIdx)

	sb := border.Clone()
	sb.Increase(SwitchDist)

	return &Cell{Border: border, SwitchBorder: sb}
}

func GetCell(pid *actor.PID) *Cell {
	return GetCellByName(pid.Id)
}

func (c *Cell) OutOfSwitchBorder(pos *msg.Vector2) bool {
	return !c.SwitchBorder.Include(pos)
}
