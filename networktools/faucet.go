package networktools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"time"
)

func (network *NetworkTools) GetFaucetBaseURL() (string, error) {
	switch network.Name {
	case "devnet":
		return "https://faucet.d.vega.xyz", nil
	case "stagnet3":
		return "https://faucet.stagnet3.vega.xyz", nil
	case "fairground":
		return "https://faucet.tmp.vega.xyz", nil
	default:
		return fmt.Sprintf("https://faucet.%s.vega.xyz", network.Name), nil
	}
}

func (network *NetworkTools) MintFakeTokens(
	vegaPubKey string,
	vegaAssetId string,
	amount *big.Int,
) error {
	errMsg := fmt.Sprintf("failed to mint %s of %s for %s fake tokens %s", amount.String(), vegaAssetId, vegaPubKey, network.Name)
	baseURL, err := network.GetFaucetBaseURL()
	if err != nil {
		return fmt.Errorf("%s, %w", errMsg, err)
	}
	url := fmt.Sprintf("%s/api/v1/mint", baseURL)

	body := map[string]string{
		"party":  vegaPubKey,
		"asset":  vegaAssetId,
		"amount": amount.String(),
	}
	byteBody, err := json.MarshalIndent(body, "", "\t")
	if err != nil {
		return fmt.Errorf("%s, %w", errMsg, err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(byteBody))
	if err != nil {
		return fmt.Errorf("%s, %w", errMsg, err)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	httpClient := http.Client{
		Timeout: time.Second * 5,
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s, %w", errMsg, err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%s, faucet response %s", errMsg, string(bodyBytes))
	}
	return nil
}
