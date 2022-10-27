package secrets

//
// Node secrets
//

// Private part stored in HashiCorp Vault
type VegaNodePrivate struct {
	// Metadata
	Name      string `json:"name"`
	Country   string `json:"country"`
	InfoURL   string `json:"info_url"`
	AvatarURL string `json:"avatar_url"`
	// Ethereum
	EthereumAddress    string `json:"ethereum_address"`
	EthereumPrivateKey string `json:"ethereum_private_key"`
	EthereumMnemonic   string `json:"ethereum_mnemonic"`
	// Vega
	VegaId             string  `json:"vega_id"`
	VegaPubKey         string  `json:"vega_public_key"`
	VegaPrivateKey     string  `json:"vega_private_key"`
	VegaRecoveryPhrase string  `json:"vega_recovery_phrase"`
	VegaPubKeyIndex    *uint64 `json:"vega_public_key_index"`
	// Tendermint
	TendermintNodeId              string `json:"tendermint_node_id"`
	TendermintNodePubKey          string `json:"tendermint_node_public_key"`
	TendermintNodePrivateKey      string `json:"tendermint_node_private_key"`
	TendermintValidatorAddress    string `json:"tendermint_validator_address"`
	TendermintValidatorPubKey     string `json:"tendermint_validator_public_key"`
	TendermintValidatorPrivateKey string `json:"tendermint_validator_private_key"`
	// Binary wallet file passphrase
	WalletBinaryPassphrase string         `json:"wallet_binary_passphrase"`
	BinaryWallets          *BinaryWallets `json:"binary_wallets"`
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
	GetVegaNodeList(network string) ([]string, error)
	GetAllVegaNode(network string) (map[string]*VegaNodePrivate, error)
}
