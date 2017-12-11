package main

import (
	"fmt"
	"os"

	"github.com/everfore/exc/walkexc/pkg"
	"github.com/spf13/cobra"
	pull "github.com/toukii/pull/pullCmd"
)

var (
	RootCommand = &cobra.Command{
		Use:   "e",
		Short: "excute cmd",
		Long:  ``,
	}

	root = new(Root)
)

func init() {
	root.CmdMap = make(map[string]*cobra.Command)
	root.CmdMap[RootCommand.Name()] = RootCommand
	root.CmdMap[pkg.Command.Name()] = pkg.Command
	root.CmdMap[pull.Command.Name()] = pull.Command

	RootCommand.AddCommand(pkg.Command)
	RootCommand.AddCommand(pull.Command)

}

type Root struct {
	args   []string
	subs   []*Sub
	CmdMap map[string]*cobra.Command
}

type Sub struct {
	cmd  string
	args []string
}

func (r *Root) Parse(args []string) {
	r.args = args
	r.subs = make([]*Sub, 0, len(r.CmdMap))
	size := len(args)
	idxes := make([]int, 0, size)
	for i, arg := range args {
		if _, ex := r.CmdMap[arg]; ex {
			idxes = append(idxes, i)
		}
	}
	idxLen := len(idxes)
	if idxLen <= 0 {
		return
	}
	preIdx := idxes[0]
	for i := 1; i <= idxLen; i++ {
		if i == idxLen {
			r.subs = append(r.subs, &Sub{
				cmd:  args[preIdx],
				args: args[preIdx+1:],
			})
			break
		}
		r.subs = append(r.subs, &Sub{
			cmd:  args[idxes[preIdx]],
			args: args[preIdx+1 : idxes[i]],
		})
		preIdx = idxes[i]
	}

	r.Excute()
}

func (r *Root) Excute() {
	for _, sub := range r.subs {
		if sub.cmd == RootCommand.Name() {
			// defer func(sub *Sub) {
			// 	r.CmdMap[sub.cmd].SetArgs(sub.args)
			// 	r.CmdMap[sub.cmd].Execute()
			// }(sub)
			continue
		}
		fmt.Println(sub.cmd, sub.args)
		r.CmdMap[sub.cmd].SetArgs(sub.args)
		// if sub.cmd == pull.Command.Name() {
		// 	pull.Command.SetArgs(sub.args)
		// 	pull.Excute(sub.args)
		// 	continue
		// }
		r.CmdMap[sub.cmd].Execute()
	}
}

func main() {
	root.Parse(os.Args)
}
