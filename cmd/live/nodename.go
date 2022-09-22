package live

import (
	"fmt"
	"math/big"
	"os"
	"sort"
	"strings"

	"crypto/rand"

	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/veganetwork"
	"go.uber.org/zap"
)

type NodenameArgs struct {
	*LiveArgs
	Random bool
	All    bool
}

var nodenameArgs NodenameArgs

// nodenameCmd represents the nodename command
var nodenameCmd = &cobra.Command{
	Use:   "nodename",
	Short: "Get nodenames of running nodes",
	Long:  `Get nodenames of running nodes`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunNodename(nodenameArgs); err != nil {
			nodenameArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	nodenameArgs.LiveArgs = &liveArgs

	LiveCmd.AddCommand(nodenameCmd)
	nodenameCmd.PersistentFlags().BoolVar(&nodenameArgs.Random, "random", false, "Randomly selects one")
	nodenameCmd.PersistentFlags().BoolVar(&nodenameArgs.All, "all", false, "Include all nodes, even unhealthy. By default only healthy nodes are returned")
}

func RunNodename(args NodenameArgs) error {
	network, err := veganetwork.NewVegaNetwork(args.VegaNetworkName, args.Logger)
	if err != nil {
		return err
	}
	var nodenames []string
	if args.All {
		nodenames = network.GetNetworkNodes()
	} else {
		nodenames = network.GetNetworkHealthyNodes()
	}

	if args.Random {
		nodeNum := len(nodenames)
		if nodeNum == 0 {
			return fmt.Errorf("all nodes are down")
		}
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(nodeNum)))
		if err != nil {
			return err
		}
		fmt.Println(nodenames[idx.Int64()])
	} else {
		sort.Strings(nodenames)
		fmt.Println(strings.Join(nodenames, ","))
	}

	return nil
}
