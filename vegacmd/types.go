package vegacmd

// {"snapshots":[{"height":5648600,"version":18830,"size":71,"hash":"80bedacff88b8069f3abfff49d42930c553632ce48ecc6f675339955edd8f24a"},
type CoreToolsSnapshot struct {
	Height  int    `json:"height"`
	Version int    `json:"version"`
	Size    int    `json:"size"`
	Hash    string `json:"hash"`
}

type CoreToolsSnapshots struct {
	Snapshots []CoreToolsSnapshot `json:"snapshots"`
}
