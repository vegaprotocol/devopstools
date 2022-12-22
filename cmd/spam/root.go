package spam

import (
	"github.com/spf13/cobra"
)

// Root Command for Spam
var SpamCmd = &cobra.Command{
	Use:   "spam",
	Short: "Set of tools for network stress-tests",
}

func init() {
	SpamCmd.AddCommand(ordersCmd)
}
