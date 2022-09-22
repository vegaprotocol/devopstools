package hcvault

import (
	"encoding/json"
	"fmt"

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
