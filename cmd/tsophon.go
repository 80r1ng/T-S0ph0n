package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	Tsophon = &cobra.Command{
		Use:   "tsophon",
		Short: "                       ---This tool is used to add users to Windows systems in specific scenarios",
		Run:   send,
	}
	logo = `
  ______                     __              
 /_  __/   _________  ____  / /_  ____  ____ 
  / /_____/ ___/ __ \/ __ \/ __ \/ __ \/ __ \
 / /_____(__  ) /_/ / /_/ / / / / /_/ / / / /
/_/     /____/\____/ .___/_/ /_/\____/_/ /_/ 
                  /_/                        
`
)

func init() {
	Tsophon.AddCommand(add)
}

func hello() {
	fmt.Println(logo)
}

func send(cmd *cobra.Command, args []string) {
	hello()
	fmt.Println(cmd.Short)
}
