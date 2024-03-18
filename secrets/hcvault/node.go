package hcvault

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/vegaprotocol/devopstools/secrets"
)

const (
	networkVaultRoot = "network"
)

func (c *HCVaultSecretStore) GetVegaNode(network string, node string) (*secrets.VegaNodePrivate, error) {
	path := fmt.Sprintf("%s/%s", network, node)
	secretDataByte, err := c.GetSecretAsByte(networkVaultRoot, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get private data for '%s' node in '%s' network; %w", node, network, err)
	}
	var result secrets.VegaNodePrivate
	if err = json.Unmarshal(secretDataByte, &result); err != nil {
		return nil, fmt.Errorf("failed to parse private data for '%s' node in '%s' network; %w", node, network, err)
	}
	return &result, nil
}

func (c *HCVaultSecretStore) StoreVegaNode(network string, node string, privateData *secrets.VegaNodePrivate) error {
	path := fmt.Sprintf("%s/%s", network, node)
	secretDataByte, err := json.Marshal(privateData)
	if err != nil {
		return fmt.Errorf("failed to parse private data for '%s' node in '%s' network; %w", node, network, err)
	}
	return c.UpsertSecretFromByte(networkVaultRoot, path, secretDataByte)
}

func (c *HCVaultSecretStore) DoesVegaNodeExist(network string, node string) (bool, error) {
	path := fmt.Sprintf("%s/%s", network, node)
	return c.DoesExist(networkVaultRoot, path)
}

func (c *HCVaultSecretStore) GetVegaNodeList(network string) ([]string, error) {
	secretList, err := c.GetSecretList(networkVaultRoot, network)
	if err != nil {
		return nil, fmt.Errorf("failed to get list of nodes for '%s' network; %w", network, err)
	}
	return secretList, nil
}

func (c *HCVaultSecretStore) GetAllVegaNode(network string) (map[string]*secrets.VegaNodePrivate, error) {
	nodeList, err := c.GetVegaNodeList(network)
	if err != nil {
		return nil, err
	}

	type Result struct {
		NodeId     string
		SecretData *secrets.VegaNodePrivate
		Err        error
	}

	resultsChannel := make(chan Result, len(nodeList))
	var wg sync.WaitGroup
	for _, node := range nodeList {
		wg.Add(1)
		go func(node string) {
			defer wg.Done()

			secret, err := c.GetVegaNode(network, node)
			resultsChannel <- Result{node, secret, err}
		}(node)
	}
	wg.Wait()
	close(resultsChannel)

	result := map[string]*secrets.VegaNodePrivate{}
	for nodeResult := range resultsChannel {
		if nodeResult.Err != nil {
			return nil, fmt.Errorf("failed to get secret for one node for network %s, %w", network, nodeResult.Err)
		}
		result[nodeResult.NodeId] = nodeResult.SecretData
	}
	return result, nil
}
