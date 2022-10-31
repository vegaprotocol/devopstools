package networktools

import (
	"fmt"

	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/types"
)

//
// Vega Core endpoints
//

func (network *NetworkTools) GetNetworkNodes() []string {
	switch network.Name {
	case types.NetworkMainnet:
		return []string{"mainnet-observer.ops.vega.xyz"}
	}
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

func (network *NetworkTools) GetNetworkDataNodes() []string {
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
				break
			} else {
				previousMissing = true
			}
		} else {
			hosts = append(hosts, host)
		}
	}
	return hosts
}

//
// GRPC
//

func (network *NetworkTools) GetNetworkGRPCDataNodes() []string {
	nodes := network.GetNetworkDataNodes()
	addresses := make([]string, len(nodes))
	for i, node := range nodes {
		addresses[i] = fmt.Sprintf("%s:3007", node)
	}
	return addresses
}

func (network *NetworkTools) GetNodeURL(nodeId string) string {
	return fmt.Sprintf("%s.%s", nodeId, network.DNSSuffix)
}
