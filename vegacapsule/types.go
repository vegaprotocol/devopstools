package vegacapsule

const ModeValidator = "validator"

type NodeDetails struct {
	Name string `json:"Name"`
	Mode string `json:"Mode"`
	Vega struct {
		HomeDir        string `json:"HomeDir"`
		ConfigFilePath string `json:"ConfigFilePath"`
		BinaryPath     string `json:"BinaryPath"`
		NodeWalletInfo struct {
			VegaWalletID        string `json:"VegaWalletID"`
			EthereumAddress     string `json:"EthereumAddress"`
			VegaWalletPublicKey string `json:"VegaWalletPublicKey"`
		} `json:"NodeWalletInfo"`
	} `json:"Vega"`
	Tendermint struct {
		HomeDir        string `json:"HomeDir"`
		ConfigFilePath string `json:"ConfigFilePath"`
	} `json:"Tendermint"`
	DataNode *struct {
		HomeDir        string `json:"HomeDir"`
		ConfigFilePath string `json:"ConfigFilePath"`
	} `json:"DataNode,omitempty"`
	Visor *struct {
		HomeDir        string `json:"HomeDir"`
		ConfigFilePath string `json:"ConfigFilePath"`
	} `json:"Visor"`
}

type VegacapsuleNodesListOut map[string]NodeDetails
