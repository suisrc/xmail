package shell

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// tool package
func GetCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "version",
			Usage: "查看版本",
			Action: func(c *cli.Context) error {
				fmt.Println(c.App.Version)
				return nil
			},
		},
	}
}
