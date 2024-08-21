package config

const (
	TypeFull      = "full"
	TypeValidator = "validator"
)

type Config struct {
	// Name of the configuration file used. It's extracted from the filename.
	Name NetworkName

	EnvironmentName NetworkName `toml:"environment_name"`

	// Bridges lists the Ethereum bridges used by the Vega network.
	Bridges Bridges `toml:"bridges"`

	// Nodes lists all the nodes participating to the network.
	Nodes []Node `toml:"nodes"`

	Bots Bots `toml:"bots"`

	Network Network `toml:"network"`

	Explorer Explorer `toml:"explorer"`
}

type Explorer struct {
	RESTURL string `toml:"rest_url"`
}

// Bridges holds the configuration of all the Ethereum bridges used by the
// Vega network.
type Bridges struct {
	// Primary configures the primary bridge used by the Vega network, usually
	// Ethereum Mainnet (or equivalent), acting as the primary collateral, multisig,
	// staking, and vesting contracts holder.
	Primary PrimaryBridge `toml:"primary"`

	// EVM configures the EVM bridge used by the Vega network, acting
	// as secondary collateral, and multisig contracts holder.
	EVM EVMBridge `toml:"evm"`
}

type PrimaryBridge struct {
	ClientURL     string          `toml:"client_url"`
	BlockExplorer Etherscan       `toml:"block_explorer"`
	Wallets       BridgeWallets   `toml:"wallets"`
	Signers       []string        `toml:"signers"`
	Versions      BridgesVersions `toml:"versions"`
}

type EVMBridge struct {
	ClientURL     string          `toml:"client_url"`
	BlockExplorer Etherscan       `toml:"block_explorer"`
	Wallets       BridgeWallets   `toml:"wallets"`
	Signers       []string        `toml:"signers"`
	Versions      BridgesVersions `toml:"versions"`
}

type BridgesVersions struct {
	Staking string `toml:"staking"`
}

// Node describes a node on the network.
type Node struct {
	ID       string       `toml:"id"`
	Type     string       `toml:"type"`
	Metadata NodeMetadata `toml:"metadata"`
	Secrets  NodeSecrets  `toml:"secrets"`
	API      NodeAPI      `toml:"api"`
}

type Bots struct {
	Trading  BotsAPI `toml:"trading"`
	Research BotsAPI `toml:"research"`
}

// Etherscan describes the Etherscan block explorer API connection.
type Etherscan struct {
	// RESTURL defines the REST endpoint used to query the block explorer API.
	RESTURL string `toml:"rest_url"`

	// APIKey defines the authentication key used to query the block explorer API.
	APIKey string `toml:"api_key"`
}

type BotsAPI struct {
	// RESTURL defines the REST endpoint used to query the bots API.
	RESTURL string `toml:"rest_url"`

	// APIKey defines the authentication key used to query the bots API.
	APIKey string `toml:"api_key"`
}

type BridgeWallets struct {
	Minter EthereumWallet `toml:"minter"`
}

type EthereumWallet struct {
	Address    string `toml:"address"`
	Mnemonic   string `toml:"mnemonic"`
	PrivateKey string `toml:"private_key"`
	Seed       string `toml:"seed"`
}

type Network struct {
	Wallets NetworkWallets `toml:"wallets"`
}

type NetworkWallets struct {
	VegaTokenWhale VegaWallet `toml:"vega_token_whale"`
	Faucet         VegaWallet `toml:"faucet"`
}

type VegaWallet struct {
	Name           string `toml:"name"`
	PublicKey      string `toml:"public_key"`
	RecoveryPhrase string `toml:"recovery_phrase"`
}

type NodeAPI struct {
	BlockchainRESTURL string `toml:"blockchain_rest_url"`
	VegaRESTURL       string `toml:"vega_rest_url"`
	VegaGRPCURL       string `toml:"vega_grpc_url"`
	DataNodeRESTURL   string `toml:"datanode_rest_url"`
	DataNodeGRPCURL   string `toml:"datanode_grpc_url"`
}

type NodeMetadata struct {
	Name      string `toml:"name"`
	Country   string `toml:"country"`
	InfoURL   string `toml:"info_url"`
	AvatarURL string `toml:"avatar_url"`
}

type NodeSecrets struct {
	// Ethereum
	EthereumAddress    string `toml:"ethereum_address"`
	EthereumPrivateKey string `toml:"ethereum_private_key"`
	EthereumMnemonic   string `toml:"ethereum_mnemonic"`
	// Vega
	VegaId             string  `toml:"vega_id"`
	VegaPubKey         string  `toml:"vega_public_key"`
	VegaPrivateKey     string  `toml:"vega_private_key"`
	VegaRecoveryPhrase string  `toml:"vega_recovery_phrase"`
	VegaPubKeyIndex    *uint64 `toml:"vega_public_key_index"`
	// Data-Node DeHistory
	DeHistoryPeerId          string `toml:"de_history_peer_id"`
	DeHistoryPrivateKey      string `toml:"de_history_private_key"`
	NetworkHistoryPeerId     string `toml:"network_history_peer_id"`
	NetworkHistoryPrivateKey string `toml:"network_history_private_key"`
	// Tendermint
	TendermintNodeId              string `toml:"tendermint_node_id"`
	TendermintNodePubKey          string `toml:"tendermint_node_public_key"`
	TendermintNodePrivateKey      string `toml:"tendermint_node_private_key"`
	TendermintValidatorAddress    string `toml:"tendermint_validator_address"`
	TendermintValidatorPubKey     string `toml:"tendermint_validator_public_key"`
	TendermintValidatorPrivateKey string `toml:"tendermint_validator_private_key"`

	// Binary wallet file passphrase
	WalletBinaryPassphrase string         `toml:"wallet_binary_passphrase"`
	BinaryWallets          *BinaryWallets `toml:"binary_wallets"`
}

type BinaryWallets struct {
	NodewalletPath       string `toml:"nodewallet_path"`
	NodewalletBase64     string `toml:"nodewallet"`
	VegaWalletPath       string `toml:"vegawallet_path"`
	VegaWalletBase64     string `toml:"vegawallet"`
	EthereumWalletPath   string `toml:"ethereumwallet_path"`
	EthereumWalletBase64 string `toml:"ethereumwallet"`
}
