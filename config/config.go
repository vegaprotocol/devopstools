package config

type Config struct {
	// Name of the configuration file used. It's extracted from the filename.
	Name string

	// Bridges lists the Ethereum bridges used by the Vega network.
	Bridges Bridges `toml:"bridges"`

	// Nodes lists all the nodes participating to the network.
	Nodes []Node

	Bots Bots

	Network Network
}

// Bridges holds the configuration of all the Ethereum bridges used by the
// Vega network.
type Bridges struct {
	// Primary configures the primary bridge used by the Vega network, usually
	// Ethereum Mainnet (or equivalent), acting as the primary collateral, multisig,
	// staking, and vesting contracts holder.
	Primary PrimaryBridge `toml:"primary"`

	// Secondary configures the secondary bridge used by the Vega network, acting
	// as secondary collateral, and multisig contracts holder.
	Secondary SecondaryBridge `toml:"secondary"`
}

type PrimaryBridge struct {
	ClientURL     string        `toml:"client_url"`
	BlockExplorer BlockExplorer `toml:"block_explorer"`

	// Assets defines the assets on the Ethereum chain.
	Assets []Asset `toml:"assets"`

	Wallets BridgeWallets
}

type SecondaryBridge struct {
	ClientURL     string        `toml:"client_url"`
	BlockExplorer BlockExplorer `toml:"block_explorer"`

	// Assets defines the assets on the Ethereum chain.
	Assets []Asset `toml:"assets"`

	Wallets BridgeWallets `toml:"wallets"`
}

// BlockExplorer describe the block explorer, such as Etherscan to get information
// from an Ethereum blockchain.
type BlockExplorer struct {
	// RESTURL defines the REST endpoint used to query the block explorer API.
	RESTURL string `toml:"rest_url"`

	// APIKey defines the authentication key used to query the block explorer API.
	APIKey string `toml:"api_key"`
}

// Node describes a node on the network.
type Node struct {
	EthereumWallet EthereumWallet `toml:"ethereum_wallet"`
	VegaWallet     VegaWallet     `toml:"vega_wallet"`
	API            NodeAPI        `toml:"api"`
}

type Bots struct {
	// RESTURL defines the REST endpoint used to query the bots API.
	RESTURL string `toml:"rest_url"`

	// APIKey defines the authentication key used to query the bots API.
	APIKey string `toml:"api_key"`
}

// Asset describes an asset from an Ethereum chain.
type Asset struct {
	// Name defines the asset's name.
	Name string `toml:"name"`
	// ContractAddress defines the hexadecimal address of the asset's contract.
	ContractAddress string `toml:"address"`
}

type BridgeWallets struct {
	Main EthereumWallet `toml:"main"`
}

type EthereumWallet struct {
	Address    string `toml:"address"`
	Mnemonic   string `toml:"mnemonic"`
	PrivateKey string `toml:"private_key"`
}

type Network struct {
	Wallets NetworkWallets
}

type NetworkWallets struct {
	VegaTokenWhale VegaWallet `toml:"vega_token_whale"`
}

type VegaWallet struct {
	ID             string `toml:"id"`
	PrivateKey     string `toml:"private_key"`
	PublicKey      string `toml:"public_key"`
	RecoveryPhrase string `toml:"recovery_phrase"`
}

type NodeAPI struct {
	BlockchainRESTURL string `toml:"blockchain_rest_url"`
	VegaGRPCURL       string `toml:"vega_grpc_url"`
	DataNodeRESTURL   string `toml:"datanode_rest_url"`
	DataNodeGRPCURL   string `toml:"datanode_grpc_url"`
	ExplorerRESTURL   string `toml:"explorer_rest_url"`
}
