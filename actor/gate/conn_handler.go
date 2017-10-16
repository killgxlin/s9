package gate

import (
	"flag"
	"fmt"
	"io"
	"s8/util"
)

// handler --------------------------------------------------
var (
	handler = util.NewCmdHandler()
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
