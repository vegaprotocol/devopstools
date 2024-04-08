package vegacmd

type CoreToolsSnapshot struct {
	Height  int    `json:"height"`
	Version int    `json:"version"`
	Size    int    `json:"size"`
	Hash    string `json:"hash"`
}

type CoreToolsSnapshots struct {
	Snapshots []CoreToolsSnapshot `json:"snapshots"`
}
