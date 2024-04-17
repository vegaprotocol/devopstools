package vegacapsule

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/vegaprotocol/devopstools/tools"
	vctools "github.com/vegaprotocol/devopstools/vegacapsule"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type GetLiveBinaryArgs struct {
	*Args

	copyTo    string
	overwrite bool
}

var getLiveBinaryArgs GetLiveBinaryArgs

// traderbotCmd represents the traderbot command
var getLiveBInaryCmd = &cobra.Command{
	Use:   "get-live-binary",
	Short: "Find the latest binary which is running on the network.",
	Long:  `Useful when you want to start new node from the protocol upgrade, but you do not know the version of the current vega`,

	Run: func(cmd *cobra.Command, args []string) {
		binaryPath, err := getLiveBinary(getLiveBinaryArgs.Logger, getLiveBinaryArgs.vegacapsuleBinary, getLiveBinaryArgs.networkHomePath)
		if err != nil {
			getLiveBinaryArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}

		fmt.Printf(`{"outputFile": "%s"}`, binaryPath)
		fmt.Println("")

		if getLiveBinaryArgs.copyTo == "" {
			return
		}

		if err := copyBiinaryTo(getLiveBinaryArgs.Logger, binaryPath, getLiveBinaryArgs.copyTo, getLiveBinaryArgs.overwrite); err != nil {
			getLiveBinaryArgs.Logger.Error("Error", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	getLiveBinaryArgs.Args = &vegacapsuleArgs

	Cmd.PersistentFlags().StringVar(
		&getLiveBinaryArgs.copyTo,
		"copy-to",
		"",
		"If not empty binary is copied to given folder")

	Cmd.PersistentFlags().BoolVar(
		&getLiveBinaryArgs.overwrite,
		"overwrite",
		true,
		"If true, binary is removed when exists")

	Cmd.AddCommand(getLiveBInaryCmd)
}

func copyBiinaryTo(logger *zap.Logger, binaryPath, outputFile string, overwrite bool) error {
	if !tools.FileExists(binaryPath) {
		return fmt.Errorf("given binary(%s) does not exists", binaryPath)
	}

	if overwrite {
		logger.Debug("Removing previous output file", zap.String("path", outputFile))
		if err := os.RemoveAll(outputFile); err != nil {
			return fmt.Errorf("failed to remove folder when --overwrite flag is provided: %w", err)
		}
	}

	logger.Debug("Copying file", zap.String("src", binaryPath), zap.String("dst", outputFile))
	if _, err := tools.CopyFile(binaryPath, outputFile); err != nil {
		return fmt.Errorf("failed to copy binary from %s to %s", binaryPath, outputFile)
	}

	logger.Debug("Updating permissions for file", zap.String("path", outputFile))
	if err := os.Chmod(outputFile, os.ModePerm); err != nil {
		return fmt.Errorf("failed to update permissions for %s: %w", outputFile, err)
	}

	return nil
}

func getLiveBinary(logger *zap.Logger, vegacapsuleBinary, networkHomePath string) (string, error) {
	vcNodes, err := vctools.ListNodes(vegacapsuleBinary, networkHomePath)
	if err != nil {
		return "", fmt.Errorf("failed to list vegacapsule nodes: %w", err)
	}

	if len(vcNodes) < 1 {
		return "", fmt.Errorf("no node found for vegacapsule network")
	}

	latestNodeBlock := 0
	selectedVegaBinaryPath := ""

	for _, nodeDetails := range vcNodes {
		coreConfigPath := nodeDetails.Vega.ConfigFilePath
		coreRESTPort, err := tools.ReadStructuredFileValue("toml", coreConfigPath, "API.REST.Port")
		if err != nil {
			return "", fmt.Errorf("failed to read structured file %s: %w", coreConfigPath, err)
		}

		nodeHeight := getBlockHeight(logger, coreRESTPort)
		if nodeHeight > latestNodeBlock {
			latestNodeBlock = nodeHeight
			// if there is visor binary could be replaced during the network lifetime
			// so we have to select path to the binary directly in the <visor_home>/current

			logger.Debug("node is on higher block than last selected binary, looking for binary", zap.String("node", nodeDetails.Name), zap.Bool("vegavisor", nodeDetails.Visor != nil))
			if nodeDetails.Visor != nil {
				selectedVegaBinaryPath, err = findVisorBinary(nodeDetails.Visor.HomeDir)
				if err != nil {
					return "", fmt.Errorf("failed to get binary path for vega node with visor config(%s): %w", nodeDetails.Name, err)
				}
			} else {
				// there is no visor, we can assume binary did not changed during the network life time

				selectedVegaBinaryPath = nodeDetails.Vega.BinaryPath
			}

			logger.Debug("binary found for node", zap.String("binary_path", selectedVegaBinaryPath), zap.String("node", nodeDetails.Name), zap.Bool("vegavisor", nodeDetails.Visor != nil))
		}
	}

	logger.Debug("selected binary for latestNode", zap.String("binary_path", selectedVegaBinaryPath))

	return selectedVegaBinaryPath, nil
}

func findVisorBinary(visorHomePath string) (string, error) {
	runConfigFilePath := filepath.Join(visorHomePath, "current", "run-config.toml")

	binary, err := tools.ReadStructuredFileValue("toml", runConfigFilePath, "vega.binary.path")
	if err != nil {
		return "", fmt.Errorf("failed to read run-config.toml: %w", err)
	}

	if filepath.IsAbs(binary) {
		return binary, nil
	}

	absoluteBinaryPath, err := filepath.Abs(binary)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for binary %s: %w", binary, err)
	}

	return absoluteBinaryPath, nil
}
