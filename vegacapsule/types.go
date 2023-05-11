package vegacapsule

type NodeDetails struct {
	Name string `json:"Name"`
	Mode string `json:"Mode"`
	Vega struct {
		HomeDir        string `json:"HomeDir"`
		ConfigFilePath string `json:"ConfigFilePath"`
	} `json:"Vega"`
	Tendermint struct {
		ConfigFilePath string `json:"ConfigFilePath"`
	} `json:"Tendermint"`
	DataNode *struct {
		HomeDir        string `json:"HomeDir"`
		ConfigFilePath string `json:"ConfigFilePath"`
	} `json:"DataNode,omitempty"`
}

type VegacapsuleNodesListOut map[string]NodeDetails
