package helper

import (
	"github.com/spf13/cobra"
)

// Root Command for Network
var HelperCmd = &cobra.Command{
	Use:   "helper",
	Short: "Some helpers for devops",
}

func init() {
}
