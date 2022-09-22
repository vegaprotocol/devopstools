package secrets

//
// Node secrets
//

// Private part stored in HashiCorp Vault
type VegaNodePrivate struct {
	// Ethereum
	EthereumAddress    string `json:"ethereum_address"`
	EthereumPrivateKey string `json:"ethereum_private_key"`
	EthereumMnemonic   string `json:"ethereum_mnemonic"`
	// Vega
	VegaPubKey         string `json:"vega_public_key"`
	VegaPrivateKey     string `json:"vega_private_key"`
	VegaRecoveryPhrase string `json:"vega_recovery_phrase"`
	// Tendermint
	TendermintNodeId              string `json:"tendermint_node_id,omitempty"`
	TendermintNodePubKey          string `json:"tendermint_node_public_key,omitempty"`
	TendermintNodePrivateKey      string `json:"tendermint_node_private_key,omitempty"`
	TendermintValidatorAddress    string `json:"tendermint_validator_address,omitempty"`
	TendermintValidatorPubKey     string `json:"tendermint_validator_public_key,omitempty"`
	TendermintValidatorPrivateKey string `json:"tendermint_validator_private_key,omitempty"`
	// Binary wallet file passphrase
	WalletBinaryPassphrase string         `json:"wallet_binary_passphrase,omitempty"`
	BinaryWallets          *BinaryWallets `json:"binary_wallets,omitempty"`
}

type BinaryWallets struct {
	NodewalletPath       string `json:"nodewallet_path"`
	NodewalletBase64     string `json:"nodewallet"`
	VegaWalletPath       string `json:"vegawallet_path"`
	VegaWalletBase64     string `json:"vegawallet"`
	EthereumWalletPath   string `json:"ethereumwallet_path"`
	EthereumWalletBase64 string `json:"ethereumwallet"`
}

type NodeSecretStore interface {
	GetVegaNode(network string, node string) (*VegaNodePrivate, error)
	StoreVegaNode(network string, node string, privateData *VegaNodePrivate) error
	DoesVegaNodeExist(network string, node string) (bool, error)
}
