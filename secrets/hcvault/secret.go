package hcvault

import (
	"encoding/json"
	"fmt"
	"strings"
)

//
// Getters
//

func (c *HCVaultSecretStore) GetSecretAsByte(root string, path string) ([]byte, error) {
	secretData, err := c.GetSecret(root, path)
	if err != nil {
		return nil, err
	}
	secretDataByte, err := json.Marshal(secretData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private data for '%s'/'%s' secret; %w", root, path, err)
	}
	return secretDataByte, nil
}

func (c *HCVaultSecretStore) GetSecret(path ...string) (map[string]interface{}, error) {
	return c.GetSecretWithPath(strings.Join(path, "/data/"))
}

func (c *HCVaultSecretStore) GetSecretWithPath(path string) (map[string]interface{}, error) {
	resp, err := c.Client.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get '%s' secret from Vega Vault: %w", path, err)
	}
	if resp == nil || resp.Data == nil || resp.Data["data"] == nil {
		return nil, fmt.Errorf("secret '%s' from Vega Vault is empty", path)
	}

	data, conversionOk := resp.Data["data"].(map[string]interface{})
	if !conversionOk {
		return nil, fmt.Errorf("failed to convert secret %s", path)
	}

	if data == nil {
		return nil, fmt.Errorf("value for secret '%s' is empty", path)
	}

	return data, nil
}

//
// Check Existence
//

func (c *HCVaultSecretStore) DoesExist(root string, path string) (bool, error) {
	return c.DoesExistWithPath(fmt.Sprintf("%s/data/%s", root, path))
}

func (c *HCVaultSecretStore) DoesExistWithPath(path string) (bool, error) {
	resp, err := c.Client.Logical().Read(path)
	if err != nil {
		return false, fmt.Errorf("failed to get '%s' secret from Vega Vault %w", path, err)
	}
	if resp == nil || resp.Data["data"] == nil {
		return false, nil
	}
	return true, nil
}

//
// Setters
//

func (c *HCVaultSecretStore) UpsertSecretFromByte(root string, path string, secretDataByte []byte) error {
	var secretData map[string]interface{}
	if err := json.Unmarshal(secretDataByte, &secretData); err != nil {
		return fmt.Errorf("failed to parse private data for '%s'/'%s'; %w", root, path, err)
	}
	return c.UpsertSecret(root, path, secretData)
}

func (c *HCVaultSecretStore) UpsertSecret(root string, path string, secret map[string]interface{}) error {
	return c.UpsertSecretWithPath(fmt.Sprintf("%s/data/%s", root, path), secret)
}

func (c *HCVaultSecretStore) UpsertSecretWithPath(path string, secret map[string]interface{}) error {
	secretData := map[string]interface{}{
		"data": secret,
	}
	_, err := c.Client.Logical().Write(path, secretData)
	if err != nil {
		return fmt.Errorf("failed to upsert '%s' secret in Vega Vault %w", path, err)
	}
	return nil
}

//
// List
//

func (c *HCVaultSecretStore) GetSecretList(root string, path string) ([]string, error) {
	fullPath := fmt.Sprintf("%s/metadata/%s", root, path)
	resp, err := c.Client.Logical().List(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret list for '%s/%s' from Vega Vault %w", root, path, err)
	}
	if resp == nil {
		return nil, fmt.Errorf("empty response for get secret list for '%s/%s' from Vega Vault", root, path)
	}

	respList := resp.Data["keys"].([]interface{})
	if respList == nil {
		return nil, fmt.Errorf("list of secrets for '%s/%s' is empty", root, path)
	}
	secretNameList := make([]string, len(respList))
	for i, name := range respList {
		secretNameList[i] = name.(string)
	}

	return secretNameList, nil
}
