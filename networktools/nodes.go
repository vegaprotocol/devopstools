package networktools

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/vegaprotocol/devopstools/tools"
	"github.com/vegaprotocol/devopstools/types"
	"golang.org/x/exp/slices"
)

type NodeType string

const (
	TypeValidator NodeType = "validator"
	TypeDataNode  NodeType = "data-node"
	TypeExplorer  NodeType = "block-explorer"
)

var AllKinds = []NodeType{TypeValidator, TypeDataNode, TypeExplorer}

func (network *NetworkTools) ListNodes(kind []NodeType) []string {
	if network.Name == types.NetworkMainnet {
		result := []string{}

		for i := 0; i < 10; i++ {
			if slices.Contains(kind, TypeExplorer) {
				result = append(result, fmt.Sprintf("be%d.%s", i, network.DNSSuffix))
			}
			if slices.Contains(kind, TypeDataNode) {
				result = append(result, fmt.Sprintf("api%d.%s", i, network.DNSSuffix))
			}
		}

		return result
	}

	result := []string{}
	if slices.Contains(kind, TypeExplorer) {
		result = append(result, fmt.Sprintf("be.%s", network.DNSSuffix))
	}

	for i := 0; i < 100; i++ {
		if slices.Contains(kind, TypeValidator) {
			result = append(result, fmt.Sprintf("n%02d.%s", i, network.DNSSuffix))
		}

		if slices.Contains(kind, TypeDataNode) {
			result = append(result, fmt.Sprintf("api.n%02d.%s", i, network.DNSSuffix))
		}
	}

	return result
}

//
// Vega Core endpoints
//

func (network *NetworkTools) checkNodes(nodes []string, healthyOnly bool) []string {
	hosts := []string{}
	previousMissing := false
	httpClient := http.Client{
		Timeout: network.restTimeout,
	}
	for _, host := range nodes {
		if _, err := tools.GetIP(host); err != nil {
			// We want to check all of the servers for mainnet
			if previousMissing && network.Name != types.NetworkMainnet {
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
			_, err := httpClient.Get(fmt.Sprintf("https://%s/statistics", host))

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

func (network *NetworkTools) GetNetworkNodes(healthyOnly bool) []string {
	return network.checkNodes(network.ListNodes([]NodeType{TypeValidator}), false)
}

func (network *NetworkTools) GetBlockExplorers(healthyOnly bool) []string {
	return network.checkNodes(network.ListNodes([]NodeType{TypeExplorer}), false)
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

	hosts := []string{}
	previousMissing := false
	httpClient := http.Client{
		Timeout: network.restTimeout,
	}
	for _, host := range network.ListNodes([]NodeType{TypeDataNode}) {
		if _, err := tools.GetIP(host); err != nil {
			if previousMissing && network.Name != types.NetworkMainnet {
				break // There is no DNS for this and previous nodes, there is no reason to check other nodes
			} else {
				previousMissing = true
			}
			continue
		}

		// Check if data-node really has statistics available for given DNS.
		if err := tools.RetryRun(3, 500*time.Millisecond, func() error {
			_, err := httpClient.Get(fmt.Sprintf("https://%s/statistics", host))

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
	addresses := []string{}
	for _, host := range network.ListNodes([]NodeType{TypeDataNode}) {
		address := net.JoinHostPort(host, "3007")
		conn, err := net.DialTimeout("tcp", address, 300*time.Millisecond)
		if err == nil && conn != nil {
			conn.Close()
			addresses = append(addresses, fmt.Sprintf("%s:3007", host))
		}
	}
	return addresses
}

func (network *NetworkTools) GetNodeURL(nodeId string) string {
	return fmt.Sprintf("%s.%s", nodeId, network.DNSSuffix)
}

//
// Tendermint endpoints
//

func (network *NetworkTools) GetNetworkTendermintRESTEndpoints(healthyOnly bool) []string {
	httpClient := http.Client{
		Timeout: network.restTimeout,
	}
	if network.Name == types.NetworkMainnet {
		result := []string{}
		for _, host := range network.ListNodes([]NodeType{TypeDataNode}) {
			url := fmt.Sprintf("http://%s:26657", host)

			if _, err := httpClient.Get(fmt.Sprintf("%s/abci_info", url)); err != nil {
				continue
			}

			result = append(result, url)
		}

		return result
	}

	hosts := []string{}
	previousMissing := false
	for _, host := range network.ListNodes([]NodeType{TypeValidator}) {
		host := fmt.Sprintf("tm.%s", host)
		if _, err := tools.GetIP(host); err != nil {
			if previousMissing {
				break
			} else {
				previousMissing = true
			}

			continue
		} else if !healthyOnly {
			hosts = append(hosts, fmt.Sprintf("https://%s", host))
			continue
		}

		network.logger.Sugar().Debugf("Checking /abci_info for %s", host)
		if _, err := httpClient.Get(fmt.Sprintf("https://%s/abci_info", host)); err != nil {
			continue
		}

		hosts = append(hosts, host)
	}
	return hosts
}
