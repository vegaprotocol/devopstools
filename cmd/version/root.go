package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "version",
	Short: "Print software version",
	Long:  `Print software version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("devopstools %s (%s)\n", cliVersion, cliVersionHash)
	},
}
