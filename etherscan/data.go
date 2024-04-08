package etherscan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

type SmartContractData struct {
	SourceCode          map[string]string
	ABI                 string
	ContractName        string
	CreationHexByteCode string
	DownloadURL         string
}

func (c *Client) GetSmartContractData(
	ctx context.Context,
	hexAddress string,
) (*SmartContractData, error) {
	response, downloadURL, err := c.GetSourcecode(ctx, hexAddress)
	if err != nil {
		return nil, err
	}
	if len(response.Result) != 1 {
		return nil, fmt.Errorf("only one result is supported from Get Sourcecode response; but it contains %d results", len(response.Result))
	}

	// Format ABI to pretty json
	abi := response.Result[0].ABI
	var prettyJSON bytes.Buffer
	if err = json.Indent(&prettyJSON, []byte(abi), "", " "); err != nil {
		return nil, fmt.Errorf("failed to format ABI for Smart Contract: %s (%s). %w", hexAddress, downloadURL, err)
	}
	abi = prettyJSON.String()

	// Parse Source code
	sourceCode, err := parseSourceCode(
		response.Result[0].ContractName,
		response.Result[0].SourceCode,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Source Code for Smart Contract: %s. %w", hexAddress, err)
	}

	creationCodeResponse, err := c.GetTxlist(ctx, hexAddress)
	if err != nil {
		return nil, err
	}
	creationCode := creationCodeResponse.Result[0].Input

	return &SmartContractData{
		SourceCode:          sourceCode,
		ABI:                 abi,
		ContractName:        response.Result[0].ContractName,
		CreationHexByteCode: creationCode,
		DownloadURL:         downloadURL,
	}, nil
}

//
// Other
//

func parseSourceCode(ContractName string, sourceCode string) (map[string]string, error) {
	if strings.HasPrefix(sourceCode, "{{") {
		sourceCode = sourceCode[1 : len(sourceCode)-1]
	}
	if strings.HasPrefix(sourceCode, "{") {
		var payload map[string]struct {
			Content string `json:"content"`
		}
		if err := json.Unmarshal([]byte(sourceCode), &payload); err != nil {
			var payload2 struct {
				Sources map[string]struct {
					Content string `json:"content"`
				} `json:"sources"`
			}
			if err := json.Unmarshal([]byte(sourceCode), &payload2); err != nil {
				return nil, fmt.Errorf("failed to parse Smart Contract %s. %w", ContractName, err)
			}
			payload = payload2.Sources
		}
		result := map[string]string{}

		for name, data := range payload {
			name = filepath.Base(name)
			result[name] = data.Content
		}

		return result, nil
	}
	return map[string]string{
		ContractName + ".sol": sourceCode,
	}, nil
}
