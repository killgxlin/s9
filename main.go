package main

import (
	"flag"
	"fmt"
	"gamelib/base/net/util"
	"log"
	"s9/actor/gate"
	"s9/imsg"
	"s9/msg"

	_ "s9/actor/auth"
	_ "s9/actor/cell"
	_ "s9/actor/scene"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/cluster/consul"
)

var (
	cport = flag.Int("cport", 8000, "cluster port")
	gport = flag.Int("gport", 9000, "gate port")
)

func testA() {
	for x := -10; x < 10; x++ {
		v := &msg.Vector2{float32(x), 0}
		name := imsg.GetCellName(v)
		c := imsg.GetCellByName(name)
		fmt.Println(v, name, c)
	}
	return
}
func testB() {
	for x := -10; x < 10; x++ {
		fmt.Println(x, x/3, x%3)
	}
	return
}

func main() {

	// consul
	cp, e := consul.New()
	if e != nil {
		log.Fatal(e)
	}
	defer cp.Shutdown()

	// cluster
	addr, e := util.FindLanAddr("tcp", *cport, *cport+1000)
	if e != nil {
		log.Panic(e)
	}
	cluster.Start("mycluster", addr, cp)

	// gate
	gate.Start(*gport, *gport+1000)

	for {
		_, e := console.ReadLine()
		if e != nil {
			break
		}
	}
}
