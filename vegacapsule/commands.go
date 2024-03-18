package vegacapsule

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vegaprotocol/devopstools/tools"

	"go.uber.org/zap"
)

type AddNodesBaseOn struct {
	Node  string
	Group string
}

func ListNodes(binary, customNetworkHomePath string) (VegacapsuleNodesListOut, error) {
	args := []string{"nodes", "ls"}
	if customNetworkHomePath != "" {
		args = append(args, "--home-path", customNetworkHomePath)
	}

	vegacapsuleNodesLsOut := VegacapsuleNodesListOut{}
	if _, err := tools.ExecuteBinary(binary, args, &vegacapsuleNodesLsOut); err != nil {
		return nil, fmt.Errorf("failed to execute vegacapsule nodes ls command: %w", err)
	}

	return vegacapsuleNodesLsOut, nil
}

func AddNodes(log *zap.Logger, binary string, filter AddNodesBaseOn, startNode bool, vegacapsuleNetworkHome string) (*NodeDetails, error) {
	if filter.Group == "" && filter.Node == "" {
		return nil, fmt.Errorf("group or node must be set in the filter")
	}
	if filter.Group != "" && filter.Node != "" {
		return nil, fmt.Errorf("set either group or node in the filter, not both")
	}

	log.Info("Adding new data node to existing network")

	tempDir, err := os.MkdirTemp("", "devopsscripts-temp")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary dir: %w", err)
	}

	defer os.RemoveAll(tempDir)

	addNodesArgs := []string{
		"nodes",
		"add",
		"--home-path", vegacapsuleNetworkHome,
		"--out-path", filepath.Join(tempDir, "create-new-node.json"),
	}

	if !startNode {
		addNodesArgs = append(addNodesArgs, "--start=false")
	}

	if filter.Node != "" {
		addNodesArgs = append(addNodesArgs, "--base-on", filter.Node)
	} else {
		addNodesArgs = append(addNodesArgs, "--base-on-group", filter.Group)
	}

	if _, err := tools.ExecuteBinary(binary, addNodesArgs, nil); err != nil {
		return nil, fmt.Errorf("failed to add new data-node to the running network: %s", err)
	}

	newNodeFileContent, err := os.ReadFile(filepath.Join(tempDir, "create-new-node.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to read details about new node: %w", err)
	}

	newNodeDetails := []NodeDetails{}

	if err := json.Unmarshal(newNodeFileContent, &newNodeDetails); err != nil {
		return nil, fmt.Errorf("failed to unmarshal new node details: %w", err)
	}

	if len(newNodeDetails) < 1 {
		return nil, fmt.Errorf("new node details are empty")
	}
	log.Info("Created new node in the vegacapsule network", zap.String("node-name", newNodeDetails[0].Name))
	log.Info("... done")

	return &newNodeDetails[0], nil
}

func StartNode(log *zap.Logger, nodeName, vegacapsuleBinary, vegacapsuleNetworkHome string) error {
	log.Info("Starting the new node", zap.String("node-name", nodeName))
	startNodeArgs := []string{
		"nodes",
		"start",
		"--name", nodeName,
		"--home-path", vegacapsuleNetworkHome,
	}

	if _, err := tools.ExecuteBinary(vegacapsuleBinary, startNodeArgs, nil); err != nil {
		return fmt.Errorf("failed to add new data-node to the running network: %s", err)
	}
	log.Info("...done")

	return nil
}
