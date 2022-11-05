package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/networktools"
	"go.uber.org/zap"
)

type SentryStatsArgs struct {
	*NetworkArgs
}

var sentryStatsArgs SentryStatsArgs

// sentryStatsCmd represents the sentryStats command
var sentryStatsCmd = &cobra.Command{
	Use:   "sentry-stats",
	Short: "Gather information about network with sentry nodes",
	Long:  `Gather information about network with sentry nodes. Highlight all Validators that are not behind Sentry Nodes`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunSentryStats(sentryStatsArgs); err != nil {
			sentryStatsArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	sentryStatsArgs.NetworkArgs = &networkArgs

	NetworkCmd.AddCommand(sentryStatsCmd)
}

const (
	GreenText = "\033[1;32m%s\033[0m"
	RedText   = "\033[1;31m%s\033[0m"
	InfoText  = "\033[1;36m%s\033[0m"
)

func RunSentryStats(args SentryStatsArgs) error {
	network, err := networktools.NewNetworkTools(args.VegaNetworkName, args.Logger)
	if err != nil {
		return err
	}

	nodesById := map[string]*InfoAboutNode{}
	endpoints := network.GetNetworkTendermintRESTEndpoints()
	knownEndpoints := map[string]bool{}
	for _, e := range endpoints {
		knownEndpoints[e] = true
	}

	tmClient := NewTendermintRESTClient()
	i := 0
	for i < len(endpoints) {
		// visit endpoint
		endpoint := endpoints[i]
		i += 1
		node, err := tmClient.GetInfoAboutNode(endpoint)
		if err != nil {
			args.Logger.Debug("failed to get info about node", zap.String("node", endpoint), zap.Error(err))
			continue
		}
		nodesById[node.NodeInfo.TendermintNodeId] = node
		// add unknown peers to be visited
		for _, peer := range node.Peers {
			if strings.HasSuffix(peer.ListenAddr, "6") {
				newEndpoint := fmt.Sprintf("http://%s7", strings.TrimSuffix(peer.ListenAddr, "6"))
				if _, ok := knownEndpoints[newEndpoint]; !ok {
					endpoints = append(endpoints, newEndpoint)
					knownEndpoints[newEndpoint] = true
				}
			}
		}
	}

	nodes := make([]*InfoAboutNode, len(nodesById))
	i = 0
	for _, node := range nodesById {
		nodes[i] = node
		sort.Slice(node.Peers, func(i, j int) bool {
			return node.Peers[i].Moniker < node.Peers[j].Moniker
		})
		i += 1
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].NodeInfo.Moniker < nodes[j].NodeInfo.Moniker
	})

	for _, nodeInfo := range nodes {
		hasPeerValidator := false
		for _, peer := range nodeInfo.Peers {
			if info, ok := nodesById[peer.TendermintNodeId]; ok && info.Validator {
				hasPeerValidator = true
				break
			}
		}

		txt := fmt.Sprintf("- %s (%s):\n", nodeInfo.NodeInfo.Moniker, nodeInfo.NodeInfo.TendermintNodeId)
		if nodeInfo.Validator {
			if hasPeerValidator {
				fmt.Printf(RedText, txt)
			} else {
				fmt.Printf(GreenText, txt)
			}
		} else {
			fmt.Printf(InfoText, txt)
		}
		for _, peer := range nodeInfo.Peers {
			txt := fmt.Sprintf("\t- %s (%s)\n", peer.Moniker, peer.TendermintNodeId)
			if info, ok := nodesById[peer.TendermintNodeId]; !ok {
				fmt.Print(txt)
			} else if info.Validator {
				if nodeInfo.Validator {
					fmt.Printf(RedText, txt)
				} else {
					fmt.Printf(GreenText, txt)
				}
			} else {
				fmt.Printf(InfoText, txt)
			}

		}
	}
	return nil
}

type TendermintRESTClient struct {
	httpClient *http.Client
}

func NewTendermintRESTClient() *TendermintRESTClient {
	return &TendermintRESTClient{
		httpClient: &http.Client{
			Timeout: time.Second * 1,
		},
	}
}

type InfoAboutNode struct {
	NodeInfo  NodeInfo
	Validator bool
	Peers     []NodeInfo
}

func (c *TendermintRESTClient) GetInfoAboutNode(endpoint string) (*InfoAboutNode, error) {
	status, err := c.GetNodeStatus(endpoint)
	if err != nil {
		return nil, err
	}
	netInfo, err := c.GetNodeNetInfo(endpoint)
	if err != nil {
		return nil, err
	}
	peers := make([]NodeInfo, len(netInfo.Peers))
	for i, peer := range netInfo.Peers {
		peers[i] = peer.NodeInfo
	}
	return &InfoAboutNode{
		NodeInfo:  status.NodeInfo,
		Validator: status.ValidatorInfo.VotingPower > 0,
		Peers:     peers,
	}, nil
}

type NodeInfo struct {
	TendermintNodeId string `json:"id"`
	Moniker          string `json:"moniker"`
	ListenAddr       string `json:"listen_addr"`
}

type NetInfo struct {
	Peers []struct {
		NodeInfo NodeInfo `json:"node_info"`
		RemoteIp string   `json:"remote_ip"`
	} `json:"peers"`
}

func (c *TendermintRESTClient) GetNodeNetInfo(endpoint string) (*NetInfo, error) {
	var (
		errMsg = "failed to get Tendermint Node Net Info, %w"
	)
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	u = u.JoinPath("net_info")
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	defer response.Body.Close()

	var payload struct {
		Result NetInfo `json:"result"`
	}
	if err = json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	return &payload.Result, nil
}

type Status struct {
	NodeInfo      NodeInfo `json:"node_info"`
	ValidatorInfo struct {
		VotingPower int64 `json:"voting_power,string"`
	} `json:"validator_info"`
}

func (c *TendermintRESTClient) GetNodeStatus(endpoint string) (*Status, error) {
	var (
		errMsg = "failed to get Tendermint Node Status, %w"
	)
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	u = u.JoinPath("status")
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	defer response.Body.Close()

	var payload struct {
		Result Status `json:"result"`
	}
	if err = json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	return &payload.Result, nil
}
