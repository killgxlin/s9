package imsg

import "s9/msg"

func (ab *AABB) Include(pos *msg.Vector2) bool {
	return pos.X >= ab.Minx && pos.X < ab.Maxx && pos.Y >= ab.Miny && pos.Y < ab.Maxy
}

func (ab *AABB) Increase(d float32) {
	ab.Minx -= d
	ab.Maxx += d
	ab.Miny -= d
	ab.Maxy += d
}

func (ab *AABB) Clone() *AABB {
	n := &AABB{
		Minx: ab.Minx,
		Maxx: ab.Maxx,
		Miny: ab.Miny,
		Maxy: ab.Maxy,
	}
	return n
}
