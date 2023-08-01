package networktools

import (
	"fmt"
	"sort"
	"strings"

	"github.com/vegaprotocol/devopstools/tools"
)

func (network *NetworkTools) GetNetworkGraph() (*NetworkGraph, error) {
	result := &NetworkGraph{Nodes: map[string]*NetworkGraphNode{}}

	tmClient := tools.NewTendermintRESTClient()
	endpoints := network.GetNetworkTendermintRESTEndpoints(false)
	knownEndpoints := map[string]bool{}
	for _, e := range endpoints {
		knownEndpoints[e] = true
	}

	i := 0
	for i < len(endpoints) {
		endpoint := endpoints[i]
		i += 1
		// Gather information
		status, err := tmClient.GetNodeStatus(endpoint)
		if err != nil {
			continue
		}
		netInfo, err := tmClient.GetNodeNetInfo(endpoint)
		if err != nil {
			return nil, err
		}
		peers := make([]string, len(netInfo.Peers))
		for i, peer := range netInfo.Peers {
			peers[i] = peer.NodeInfo.TendermintNodeId
		}
		// Add Node to graph
		if err = result.AddNode(
			status.NodeInfo.TendermintNodeId,
			status.NodeInfo.Moniker,
			status.ValidatorInfo.VotingPower > 0,
			peers,
		); err != nil {
			return nil, err
		}
		// Add Peers to graph (in case we can't connect to the Node)
		for _, peer := range netInfo.Peers {
			if err = result.AddNode(
				peer.NodeInfo.TendermintNodeId,
				peer.NodeInfo.Moniker,
				false,
				[]string{},
			); err != nil {
				return nil, err
			}
		}

		// discover endpoints
		for _, peer := range netInfo.Peers {
			if strings.HasSuffix(peer.NodeInfo.ListenAddr, "6") {
				newEndpoint := fmt.Sprintf("http://%s7", strings.TrimSuffix(peer.NodeInfo.ListenAddr, "6"))
				if _, ok := knownEndpoints[newEndpoint]; !ok {
					endpoints = append(endpoints, newEndpoint)
					knownEndpoints[newEndpoint] = true
				}
			}
		}
	}
	// sort Peers for each node
	for _, node := range result.Nodes {
		sort.Slice(node.Peers, func(i, j int) bool {
			return result.Nodes[node.Peers[i]].Moniker < result.Nodes[node.Peers[j]].Moniker
		})
	}
	return result, nil
}

type NetworkGraphNode struct {
	TendermintNodeId string
	Moniker          string
	Validator        bool
	Peers            []string
}

type NetworkGraph struct {
	Nodes map[string]*NetworkGraphNode
}

func (ng *NetworkGraph) AddNode(
	tendermintNodeId string,
	moniker string,
	validator bool,
	peers []string,
) error {
	if existingNodeData, ok := ng.Nodes[tendermintNodeId]; ok {
		if validator {
			existingNodeData.Validator = validator
		}
		if len(existingNodeData.Peers) == 0 {
			existingNodeData.Peers = peers
		}
		if existingNodeData.Moniker != moniker {
			return fmt.Errorf("inconsistent data for %s node: Moniker(%s=%s)", tendermintNodeId,
				existingNodeData.Moniker, moniker)
		}
	} else {
		ng.Nodes[tendermintNodeId] = &NetworkGraphNode{
			TendermintNodeId: tendermintNodeId,
			Moniker:          moniker,
			Validator:        validator,
			Peers:            peers,
		}
	}
	return nil
}

const (
	GreenText = "\033[1;32m%s\033[0m"
	RedText   = "\033[1;31m%s\033[0m"
	InfoText  = "\033[1;36m%s\033[0m"
)

func (ng *NetworkGraph) Print() {
	nodes := make([]*NetworkGraphNode, len(ng.Nodes))
	i := 0
	for _, node := range ng.Nodes {
		nodes[i] = node
		i += 1
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Moniker < nodes[j].Moniker
	})

	for _, nodeInfo := range nodes {
		hasPeerValidator := false
		for _, peer := range nodeInfo.Peers {
			if peerNode, ok := ng.Nodes[peer]; ok && peerNode.Validator {
				hasPeerValidator = true
				break
			}
		}

		txt := fmt.Sprintf("- %s (%s):\n", nodeInfo.Moniker, nodeInfo.TendermintNodeId)
		if nodeInfo.Validator {
			if hasPeerValidator {
				fmt.Printf(RedText, txt)
			} else {
				fmt.Printf(GreenText, txt)
			}
		} else {
			if len(nodeInfo.Peers) == 0 {
				fmt.Print(txt)
			} else {
				fmt.Printf(InfoText, txt)
			}
		}
		for _, peer := range nodeInfo.Peers {
			peerNode := ng.Nodes[peer]
			txt := fmt.Sprintf("\t- %s (%s)\n", peerNode.Moniker, peer)
			if peerNode.Validator {
				if nodeInfo.Validator {
					fmt.Printf(RedText, txt)
				} else {
					fmt.Printf(GreenText, txt)
				}
			} else {
				if len(peerNode.Peers) == 0 {
					fmt.Print(txt)
				} else {
					fmt.Printf(InfoText, txt)
				}
			}
		}
	}
}
