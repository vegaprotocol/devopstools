package secrets

//
// Wallets
//

type EthereumWalletPrivate struct {
	Address    string `json:"address"`
	Mnemonic   string `json:"mnemonic"`
	Seed       string `json:"seed"`
	PrivateKey string `json:"private_key"`
}

type VegaWalletPrivate struct {
	Id             string `json:"id"`
	PublicKey      string `json:"public_key"`
	PrivateKey     string `json:"private_key"`
	RecoveryPhrase string `json:"recovery_phrase"`
}

type WalletSecretStore interface {
	GetEthereumWallet(secretPath string) (*EthereumWalletPrivate, error)
	StoreEthereumWallet(secretPath string, secretData *EthereumWalletPrivate) error
	DoesEthereumWalletExist(secretPath string) (bool, error)

	GetVegaWallet(secretPath string) (*VegaWalletPrivate, error)
	StoreVegaWallet(secretPath string, secretData *VegaWalletPrivate) error
	DoesVegaWalletExist(secretPath string) (bool, error)
}
