package vegacapsule

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vegaprotocol/devopstools/tools"
	vctools "github.com/vegaprotocol/devopstools/vegacapsule"
	"github.com/vegaprotocol/devopstools/vegasnapshot"
	"go.uber.org/zap"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/state"
	"github.com/cometbft/cometbft/store"
)

type LoadMainnetSnapshotArgs struct {
	*VegacapsuleArgs

	snapshotSourcePath string
	workDirPath        string
}

var loadMainnetSnapshotArgs LoadMainnetSnapshotArgs

// traderbotCmd represents the traderbot command
var loadMainnetSnapshotCmd = &cobra.Command{
	Use:   "load-mainnet-snapshot",
	Short: "Load mainnet snapshot into the generated vegacapsule-network",
	Long:  `Snapshot must be downloaded to local file system. To download snapshot see the 'devopstools remote download-snapshot' command.`,

	Run: func(cmd *cobra.Command, args []string) {

		if err := loadSnapshot(
			loadMainnetSnapshotArgs.Logger,
			loadMainnetSnapshotArgs.vegacapsuleBinary,
			loadMainnetSnapshotArgs.networkHomePath,
			loadMainnetSnapshotArgs.snapshotSourcePath,
			loadMainnetSnapshotArgs.workDirPath,
		); err != nil {
			loadMainnetSnapshotArgs.Logger.Fatal("failed to load snapshot", zap.Error(err))
		}

		// checkErr(tendermintState())
	},
}

func init() {
	loadMainnetSnapshotArgs.VegacapsuleArgs = &vegacapsuleArgs

	loadMainnetSnapshotCmd.PersistentFlags().StringVar(
		&loadMainnetSnapshotArgs.snapshotSourcePath,
		"snapshot-source-path",
		"",
		"Path to the snapshot source downloaded from the mainnet")

	loadMainnetSnapshotCmd.PersistentFlags().StringVar(
		&loadMainnetSnapshotArgs.workDirPath,
		"work-dir-path",
		"./",
		"Path to the work dir")

	if err := loadMainnetSnapshotCmd.MarkPersistentFlagRequired("snapshot-source-path"); err != nil {
		panic(err)
	}

	VegacapsuleCmd.AddCommand(loadMainnetSnapshotCmd)
}

type TendermintPriv struct {
	Address string `json:"address"`
	PubKey  struct {
		Value string `json:"value"`
	} `json:"pub_key"`
}

func describeValidatorsForSnapshot(vegacapsuleBinary, vegacapsuleHome string) ([]vegasnapshot.ValidatorPublicData, error) {
	nodesDetails, err := vctools.ListNodes(vegacapsuleBinary, vegacapsuleHome)
	if err != nil {
		return nil, fmt.Errorf("failed to list vegacapsule nodes: %w", err)
	}

	result := []vegasnapshot.ValidatorPublicData{}
	for _, node := range nodesDetails {
		// no validator, we are not interested in this node
		if node.Mode != vctools.ModeValidator {
			continue
		}

		tmPrivKeyPath := filepath.Join(node.Tendermint.HomeDir, "config", "priv_validator_key.json")
		if !tools.FileExists(tmPrivKeyPath) {
			return nil, fmt.Errorf("priv_validator_key.json file for note %s does not exists", node.Name)
		}

		privKeyContent, err := os.ReadFile(tmPrivKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read priv_validator_key.json file: %w", err)
		}

		tmPrivKey := TendermintPriv{}
		if err := json.Unmarshal(privKeyContent, &tmPrivKey); err != nil {
			return nil, fmt.Errorf("failed unmarshal priv_validator_key.json file: %w", err)
		}

		// todo: check if node wallet info is not empty

		result = append(result, vegasnapshot.ValidatorPublicData{
			EthereumAddress:  node.Vega.NodeWalletInfo.EthereumAddress,
			NodeId:           node.Vega.NodeWalletInfo.VegaWalletID,
			VegaPubKey:       node.Vega.NodeWalletInfo.VegaWalletPublicKey,
			TmPubKey:         tmPrivKey.PubKey.Value,
			ValidatorAddress: tmPrivKey.Address,
		})
	}

	return result, nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func tendermintState() error {
	dbType := dbm.BackendType("goleveldb")

	if !tools.FileExists(filepath.Join("/home/daniel/www/vega/tm_home/data", "blockstore.db")) {
		return fmt.Errorf("no blockstore found in %v", "/home/daniel/vega/tm_home/data")
	}

	// Get BlockStore
	blockStoreDB, err := dbm.NewDB("blockstore", dbType, "/home/daniel/www/vega/tm_home/data")
	if err != nil {
		return err
	}
	blockStore := store.NewBlockStore(blockStoreDB)

	if !tools.FileExists(filepath.Join("/home/daniel/www/vega/tm_home/data", "state.db")) {
		return fmt.Errorf("no statestore found in %v", "/home/daniel/www/vega/tm_home/data")
	}

	// Get StateStore
	stateDB, err := dbm.NewDB("state", dbType, "/home/daniel/www/vega/tm_home/data")
	if err != nil {
		return err
	}
	stateStore := state.NewStore(stateDB, state.StoreOptions{
		DiscardABCIResponses: false, // config.Storage.DiscardABCIResponses,
	})

	tmState, err := stateStore.Load()

	checkErr(err)
	fmt.Println(tmState)
	_ = blockStore
	_ = stateStore

	checkErr(SaveToJSON(tmState.Copy(), "./tm-state.json"))

	lastBlock := blockStore.Height()
	blockResp := blockStore.LoadBlock(lastBlock)

	fmt.Printf("State.AppHash: %X, Block.AppHash: %v", tmState.AppHash, blockResp.AppHash)
	// tmState.AppHash = []byte("E3D1875A730A0902EB075369103446C74AA0E8DC65FD9A1C574B72606CCC2109")
	checkErr(stateStore.Save(tmState))

	// fmt.Printf("State.AppHash: %X, Block.AppHash: %v", tmState.AppHash, blockResp.AppHash)
	checkErr(SaveToJSON(blockResp, "./tm-block.json"))

	// 	return blockStore, stateStore, nil
	return nil
}

func SaveToJSON[T any](state T, outputFile string) error {
	bytes, err := json.MarshalIndent(state, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal state into JSON: %w", err)
	}

	if err := os.WriteFile(outputFile, bytes, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write content into %s: %w", outputFile, err)
	}

	return nil
}

func loadSnapshot(logger *zap.Logger, vegacapsuleBinary, vegacapsuleHomePath, snapshotSourcePath, workDirPath string) error {
	// TODO: Move it to separate function
	validatorsData, err := describeValidatorsForSnapshot(vegacapsuleBinary, vegacapsuleHomePath)
	if err != nil {
		return fmt.Errorf("failed to get validators details for snapsghot: %w", err)
	}

	snapshot, err := vegasnapshot.OpenSnapshotDB(snapshotSourcePath)
	if err != nil {
		return fmt.Errorf("failed to open db snapshot: %w", err)
	}

	snapshotJSONOutput := filepath.Join(workDirPath, "./snapshot-before-update.json")

	if err := snapshot.WriteSnapshotAsJSON(snapshotJSONOutput); err != nil {
		return fmt.Errorf("failed to write snapshot before update to the JSON file: %w", err)
	}

	if err := snapshot.UpdateSnapshot(validatorsData, true); err != nil {
		return fmt.Errorf("failed update snapshot: %w", err)
	}

	snapshotJSONOutput = filepath.Join(workDirPath, "./snapshot-after-update.json")

	if err := snapshot.WriteSnapshotAsJSON(snapshotJSONOutput); err != nil {
		return fmt.Errorf("failed to write snapshot after update to the JSON file: %w", err)
	}

	nodesDetails, err := vctools.ListNodes(vegacapsuleBinary, vegacapsuleHomePath)
	if err != nil {
		return fmt.Errorf("failed to list nodes for vegacapsule network: %w", err)
	}

	for _, node := range nodesDetails {
		vegaHomePath := node.Vega.HomeDir
		snapshotDirPath := filepath.Join(vegaHomePath, "/state/node/snapshots")

		logger.Info(fmt.Sprintf("Ensuring snapshot directory exists for node %s", node.Name), zap.String("path", snapshotDirPath))
		if err := os.MkdirAll(snapshotDirPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to ensure snapshot dir exists for node %s: %w", node.Name, err)
		}

		if err := tools.CopyDirectory(snapshotSourcePath, snapshotDirPath); err != nil {
			return fmt.Errorf("failed to copy mainnet snapshot for node %s: %w", node.Name, err)
		}
	}

	return nil
}
