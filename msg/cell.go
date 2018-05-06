package msg

import (
	"fmt"
	"gamelib/base/util"
	"strconv"
	"strings"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/remote"
)

const (
	CellSize   = 20
	SwitchDist = 2
	MirrorDist = 4
)

func GetCellName(pos *Vector2) string {
	xIdx := int(pos.X / CellSize)
	if pos.X < 0 {
		xIdx = xIdx - 1
	}
	yIdx := int(pos.Y / CellSize)
	if pos.Y < 0 {
		yIdx = yIdx - 1
	}
	return fmt.Sprintf("cell_%d_%d", xIdx, yIdx)
}

func GetCellPIDByName(name string) *actor.PID {
	pid, e := cluster.Get(name, "cell")
	if e != remote.ResponseStatusCodeOK {
		panic(e)
	}
	return pid
}

func GetCellPID(pos *Vector2) *actor.PID {
	name := GetCellName(pos)
	pid, e := cluster.Get(name, "cell")
	if e != remote.ResponseStatusCodeOK {
		panic(e)
	}
	return pid
}

func getRangeByIdx(idx int) (float32, float32) {
	minx := float32(float64(idx) * CellSize)
	maxx := float32(float64(idx+1) * CellSize)
	return minx, maxx
}

func getIndexByName(name string) (int, int) {
	f := strings.Split(name, "_")

	xIdx, e := strconv.Atoi(f[1])
	util.PanicOnErr(e)
	yIdx, e := strconv.Atoi(f[2])
	util.PanicOnErr(e)

	return xIdx, yIdx
}

func getCellByIndex(xIdx, yIdx int) *Cell {
	border := &AABB{}
	border.Minx, border.Maxx = getRangeByIdx(xIdx)
	border.Miny, border.Maxy = getRangeByIdx(yIdx)

	sb := border.Clone()
	sb.Increase(SwitchDist)

	mb := border.Clone()
	mb.Increase(MirrorDist)

	name := fmt.Sprintf("cell_%d_%d", xIdx, yIdx)

	return &Cell{Name: name, Border: border, SwitchBorder: sb, MirrorBorder: mb}
}

func GetCellByName(name string) *Cell {
	return getCellByIndex(getIndexByName(name))
}

func GetCell(pid *actor.PID) *Cell {
	name := strings.Split(pid.Id, "$")[1]
	return GetCellByName(name)
}

func (c *Cell) OutOfSwitchBorder(pos *Vector2) bool {
	return !c.SwitchBorder.Include(pos)
}

func (c *Cell) InGhostBorder(pos *Vector2) bool {
	return c.MirrorBorder.Include(pos)
}

func GenNeighbours(name string) []*Cell {
	ret := []*Cell{}
	xIdx, yIdx := getIndexByName(name)
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			ret = append(ret, getCellByIndex(i+xIdx, j+yIdx))
		}
	}
	return ret
}
