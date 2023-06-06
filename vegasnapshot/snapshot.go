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

//	alidatorUpdate": {
//		"nodeId": "fbbf4f6aacd78d69ccb80bd9e95ac7c27dbefbdff5458c50b38727a5faceb968",
//		"vegaPubKey": "3763b69285af3d59679cf066a89e416005eecf0032dc669a6b80eb8d9a8c0026",
//		"ethereumAddress": "0x8ae64E2a56D7581C848A8CAed7cAf7695B9aD6ad",
//		"tmPubKey": "vzCL3m3mJAvc0wVjp5hv4HOwpNFoeV17Ssha6Y2TUFA=",
type ValidatorPublicData struct {
	NodeId           string // WalletID
	VegaPubKey       string // WalletPubKey
	EthereumAddress  string
	TmPubKey         string // prov_validator_key.json
	ValidatorAddress string // localhost:26607/validators
}

func (vs VegaSnapshot) updateTopology(payload *snapshot.Payload_Topology, newValidators []ValidatorPublicData) (*snapshot.Payload, error) {
	for idx := range payload.Topology.ValidatorData {
		// not enough validators in the vegacapsule network
		if len(newValidators) <= idx {
			break
		}

		newValidator := newValidators[idx]

		payload.Topology.ValidatorData[idx].ValidatorUpdate.NodeId = newValidator.NodeId
		payload.Topology.ValidatorData[idx].ValidatorUpdate.VegaPubKey = newValidator.VegaPubKey
		payload.Topology.ValidatorData[idx].ValidatorUpdate.EthereumAddress = newValidator.EthereumAddress
		payload.Topology.ValidatorData[idx].ValidatorUpdate.TmPubKey = newValidator.TmPubKey
	}

	for idx := range payload.Topology.ValidatorPerformance.ValidatorPerfStats {
		if len(newValidators) <= idx {
			break
		}

		newValidator := newValidators[idx]

		payload.Topology.ValidatorPerformance.ValidatorPerfStats[idx].ValidatorAddress = newValidator.ValidatorAddress
	}

	// TODO: Update the payload.Topology.PendingPubKeyRotations and other stuff here when needed
	return &snapshot.Payload{
		Data: payload,
	}, nil
}

func (vs *VegaSnapshot) UpdateSnapshot(newValidators []ValidatorPublicData, save bool) error {
	stopped, err := vs.root.Iterate(func(key, val []byte) bool {
		p := &snapshot.Payload{}
		if err := proto.Unmarshal(val, p); err != nil {
			return true
		}
		switch dt := p.Data.(type) {
		case *snapshot.Payload_Topology:
			newTopology, err := vs.updateTopology(dt, newValidators)
			if err != nil {
				fmt.Printf("Failed update topology: %s\n", err.Error())
				return true
			}

			topologyBytes, err := proto.Marshal(newTopology)
			if err != nil {
				fmt.Printf("Failed to marshal payload: %s\n", err.Error())
				return true
			}
			vs.root.Set(key, topologyBytes)
		}

		return false
	})

	if err != nil {
		return fmt.Errorf("failed to update validators in the snapshot: %w", err)
	}

	if stopped {
		return fmt.Errorf("update validator has not finished")
	}

	if !save {
		return nil
	}

	if _, _, err := vs.root.SaveVersion(); err != nil {
		return fmt.Errorf("failed to write new version of db snapshot: %w", err)
	}

	return nil
}

func (vs *VegaSnapshot) Save() error {
	if _, _, err := vs.root.SaveVersion(); err != nil {
		return fmt.Errorf("failed to write new version of db snapshot: %w", err)
	}

	return nil
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
