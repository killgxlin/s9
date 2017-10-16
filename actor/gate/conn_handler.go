package gate

import (
	"flag"
	"fmt"
	"gamelib/base/handler/cmdhandler"
	"io"
)

// handler --------------------------------------------------
var (
	handler = cmdhandler.NewCmdHandler()
)

func init() {
	handler.Register("echo", func(args []string, fs *flag.FlagSet, out io.Writer, ctx interface{}) {
		msg := fs.String("m", "", "message for echo")
		e := fs.Parse(args)
		if e != nil {
			fmt.Fprintln(out, e)
			return
		}

		fmt.Fprintln(out, *msg)
	})
}
