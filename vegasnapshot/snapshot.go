package vegasnapshot

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	snapshot "code.vegaprotocol.io/vega/protos/vega/snapshot/v1"
	"github.com/cosmos/iavl"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	db "github.com/tendermint/tm-db"
	"google.golang.org/protobuf/proto"
)

type VegaSnapshot struct {
	dbPath string

	root *iavl.MutableTree
}

func OpenSnapshotDB(dbPath string) (*VegaSnapshot, error) {
	options := &opt.Options{
		ErrorIfMissing: true,
		ReadOnly:       false,
	}
	db, err := db.NewGoLevelDBWithOpts("snapshot", dbPath, options)
	if err != nil {
		return nil, fmt.Errorf("failed to open database located at %s : %w", dbPath, err)
	}

	tree, err := iavl.NewMutableTree(db, 0, false)
	if err != nil {
		return nil, fmt.Errorf("failed to open snapshot db: %w", err)
	}

	if _, err := tree.Load(); err != nil {
		return nil, fmt.Errorf("failed to load goleveldb into memory: %w", err)
	}

	return &VegaSnapshot{
		dbPath: dbPath,
		root:   tree,
	}, nil
}

func (vs *VegaSnapshot) getAllPayloads() ([]*snapshot.Payload, uint64, error) {
	payloads := []*snapshot.Payload{}
	var err error
	var blockHeight uint64
	vs.root.Iterate(func(key []byte, val []byte) (stop bool) {
		p := &snapshot.Payload{}
		err = proto.Unmarshal(val, p)
		if err != nil {
			return true
		}

		// grab block-height while we're here
		switch dt := p.Data.(type) {
		case *snapshot.Payload_AppState:
			blockHeight = dt.AppState.Height
		}

		payloads = append(payloads, p)
		return false
	})

	return payloads, blockHeight, err
}

func (vs VegaSnapshot) WriteSnapshotAsJSON(outputPath string) error {
	// traverse the tree and get the payloads

	payloads, _, err := vs.getAllPayloads()
	if err != nil {
		return err
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	m := jsonpb.Marshaler{Indent: "    "}

	payloadMap := map[string]interface{}{}
	for _, p := range payloads {
		s, err := m.MarshalToString(p)
		if err != nil {
			return fmt.Errorf("failed to unmarshal payload to string: %w", err)
		}

		singlePayload := map[string]interface{}{}

		if err := json.Unmarshal([]byte(s), &singlePayload); err != nil {
			return fmt.Errorf("failed to unmarshal json into map: %w", err)
		}

		for k, v := range singlePayload {
			payloadMap[k] = v
		}
	}

	jsonOut, err := json.MarshalIndent(payloadMap, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal final json: %w", err)
	}

	w.Write(jsonOut)
	w.Flush()

	return nil
}
