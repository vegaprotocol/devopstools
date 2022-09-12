package ops

import (
	"github.com/spf13/cobra"
)

// Root Command for OPS
var OpsCmd = &cobra.Command{
	Use:   "ops",
	Short: "General ops tasks",
	Long:  `Range of OPS tasks`,
}
