package networktools

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/types"
)

//
// Vega Core endpoints
//

func (network *NetworkTools) GetNetworkNodes(healthyOnly bool) []string {
	hosts := []string{}
	previousMissing := false
	for i := 0; i < 100; i++ {
		host := fmt.Sprintf("n%02d.%s", i, network.DNSSuffix)
		if _, err := tools.GetIP(host); err != nil {
			if previousMissing {
				break
			} else {
				previousMissing = true
			}
			continue
		}

		// Return all nodes that resolves to IP
		if !healthyOnly {
			hosts = append(hosts, host)
			continue
		}

		// Check if the node really has statistics available for given DNS.
		if err := tools.RetryRun(3, 500*time.Millisecond, func() error {
			_, err := http.Get(fmt.Sprintf("https://%s/statistics", host))

			return err
		}); err != nil {
			network.logger.Sugar().Debugf("Node %s missing", host)
			continue
		} else {
			hosts = append(hosts, host)
		}
	}
	return hosts
}

func (network *NetworkTools) GetNetworkHealthyNodes() []string {
	hostStats := network.GetRunningStatisticsForAllHosts()
	nodenames := make([]string, 0, len(hostStats))
	for oneNodename := range hostStats {
		nodenames = append(nodenames, oneNodename)
	}
	return nodenames
}

//
// Data-Node endpoints
//

func (network *NetworkTools) GetNetworkDataNodes(healthyOnly bool) []string {
	switch network.Name {
	case types.NetworkMainnet:
		return []string{"api.vega.xyz"}
	}
	hosts := []string{}
	previousMissing := false
	for i := 0; i < 100; i++ {
		host := fmt.Sprintf("api.n%02d.%s", i, network.DNSSuffix)
		if _, err := tools.GetIP(host); err != nil {
			if previousMissing {
				break // There is no DNS for this and previous nodes, there is no reason to check other nodes
			} else {
				previousMissing = true
			}
			continue
		}

		// Check if data-node really has statistics available for given DNS.
		if err := tools.RetryRun(3, 500*time.Millisecond, func() error {
			_, err := http.Get(fmt.Sprintf("https://%s/statistics", host))

			return err
		}); err != nil {
			network.logger.Sugar().Debugf("Node %s missing", host)
			continue
		} else {
			hosts = append(hosts, host)
		}
	}
	return hosts
}

//
// GRPC
//

func (network *NetworkTools) GetNetworkGRPCVegaCore() []string {
	nodes := network.GetNetworkNodes(false)
	addresses := make([]string, len(nodes))
	for i, node := range nodes {
		addresses[i] = fmt.Sprintf("%s:3002", node)
	}
	return addresses
}

func (network *NetworkTools) GetNetworkGRPCDataNodes() []string {
	nodes := network.GetNetworkDataNodes(false)
	addresses := make([]string, len(nodes))
	for i, node := range nodes {
		addresses[i] = fmt.Sprintf("%s:3007", node)
	}
	return addresses
}

func (network *NetworkTools) GetNodeURL(nodeId string) string {
	return fmt.Sprintf("%s.%s", nodeId, network.DNSSuffix)
}

//
// Tendermint endpoints
//

func (network *NetworkTools) GetNetworkTendermintRESTEndpoints() []string {
	switch network.Name {
	case types.NetworkMainnet:
		return []string{"http://api2.vega.xyz:26657", "http://api3.vega.xyz:26657"}
	}
	hosts := []string{}
	previousMissing := false
	for i := 0; i < 100; i++ {
		host := fmt.Sprintf("tm.n%02d.%s", i, network.DNSSuffix)
		if _, err := tools.GetIP(host); err != nil {
			if previousMissing {
				break
			} else {
				previousMissing = true
			}
		} else {
			hosts = append(hosts, fmt.Sprintf("https://%s", host))
		}
	}
	return hosts
}
