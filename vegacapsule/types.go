package vegacapsule

const (
	VegaModeValidator = "validator"
	VegaModeFull      = "full"
)

type NodeDetails struct {
	Name string `json:"Name"`
	Mode string `json:"Mode"`
	Vega struct {
		HomeDir        string `json:"HomeDir"`
		ConfigFilePath string `json:"ConfigFilePath"`
		BinaryPath     string `json:"BinaryPath"`
	} `json:"Vega"`
	Tendermint struct {
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
