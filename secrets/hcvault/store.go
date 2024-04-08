package hcvault

import (
	"fmt"
	"io"
	"os"
	"strings"

	vault "github.com/hashicorp/vault/api"
)

type SecretStore struct {
	Client *vault.Client
}

type LoginToken struct {
	VaultToken          string
	FileWithVaultToken  string
	GitHubToken         string
	FileWithGitHubToken string
}

func NewHCVaultSecretStore(
	vaultURL string,
	loginToken LoginToken,
) (*SecretStore, error) {
	var token string

	config := vault.DefaultConfig()
	config.Address = vaultURL

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create HashiCorp Vault client: %w", err)
	}

	if loginToken.VaultToken != "" {
		token = loginToken.VaultToken
	} else if loginToken.FileWithVaultToken != "" {
		token, err = readTokenFromFile(loginToken.FileWithVaultToken)
		if err != nil {
			return nil, err
		}
	} else {
		var gitHubToken string
		if loginToken.GitHubToken != "" {
			gitHubToken = loginToken.GitHubToken
		} else {
			gitHubToken, err = readTokenFromFile(loginToken.FileWithGitHubToken)
			if err != nil {
				return nil, err
			}
		}
		token, err = loginWithGitHubToken(client, gitHubToken)
		if err != nil {
			return nil, err
		}
	}

	client.SetToken(token)
	return &SecretStore{
		Client: client,
	}, nil
}

func readTokenFromFile(fileWithToken string) (string, error) {
	secretFile, err := os.Open(fileWithToken)
	if err != nil {
		return "", fmt.Errorf("unable to open file %s containing token: %w", fileWithToken, err)
	}
	defer secretFile.Close()

	limitReader := io.LimitReader(secretFile, 100)
	tokenBytes, err := io.ReadAll(limitReader)
	if err != nil {
		return "", fmt.Errorf("unable to read token from file %s: %w", fileWithToken, err)
	}
	token := strings.TrimSuffix(string(tokenBytes), "\n")
	return token, nil
}

func loginWithGitHubToken(client *vault.Client, gitHubToken string) (string, error) {
	loginData := map[string]interface{}{
		"token": gitHubToken,
	}

	resp, err := client.Logical().Write("auth/github/login", loginData)
	if err != nil {
		return "", fmt.Errorf("failed to login using GitHub token %w", err)
	}

	return resp.Auth.ClientToken, nil
}
