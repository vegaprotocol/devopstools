package snapshotcompatibility

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/vegaprotocol/devopstools/tools"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type DownloadBinaryArgs struct {
	*Args

	Destination      string
	CoreRESTEndpoint string
	Repository       string
}

var downloadBinaryArgs DownloadBinaryArgs

var downloadBinaryCmd = &cobra.Command{
	Use:   "download-binary",
	Short: "Download the binary which is running on given core endpoint",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runDownloadBinary(downloadBinaryArgs.Logger,
			downloadBinaryArgs.Destination,
			downloadBinaryArgs.CoreRESTEndpoint,
			downloadBinaryArgs.Repository,
		); err != nil {
			downloadMainnetSnapshotArgs.Logger.Fatal(
				"failed to prepare for snapshot compatibility pipeline",
				zap.Error(err),
			)

			return
		}
	},
}

func init() {
	downloadBinaryArgs.Args = &snapshotCompatibilityArgs

	downloadBinaryCmd.PersistentFlags().
		StringVar(
			&downloadBinaryArgs.Destination,
			"destination",
			"./vega-mainnet",
			"Path, where the binary is downloaded",
		)
	downloadBinaryCmd.PersistentFlags().
		StringVar(
			&downloadBinaryArgs.CoreRESTEndpoint,
			"core-rest-endpoint",
			"https://api2.vega.community",
			"URL to the vega core REST API",
		)
	downloadBinaryCmd.PersistentFlags().
		StringVar(
			&downloadBinaryArgs.Repository,
			"repository",
			"vegaprotocol/vega",
			"Repository, where are the binaries",
		)
}

type statisticsResponse struct {
	Statistics struct {
		AppVersion string `json:"appVersion"`
	} `json:"statistics"`
}

func getStatistics(endpoint string) (*statisticsResponse, error) {
	if !strings.HasPrefix(endpoint, "http") {
		endpoint = fmt.Sprintf("https://%s", endpoint)
	}

	url := strings.Join([]string{strings.TrimRight(endpoint, "/"), "statistics"}, "/")

	return tools.RetryReturn(3, 3*time.Second, func() (*statisticsResponse, error) {
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to get response from statistic: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body for statistics: %w", err)
		}

		result := &statisticsResponse{}
		if err := json.Unmarshal(body, result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal statistics response: %w", err)
		}

		return result, nil
	})
}

func downloadURL(binaryVersion, repository, artifactName string) (string, error) {
	url := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", repository, binaryVersion, artifactName)

	resp, err := http.Head(url)
	if err == nil && resp.StatusCode == http.StatusOK {
		return url, nil
	}

	return "", fmt.Errorf("failed to get existing url for asset binary %s/%s: %w", binaryVersion, artifactName, err)
}

func downloadVegaBinary(logger *zap.Logger, repository, version, destinationFile string) error {
	artifactName := tools.FormatAssetName("vega", "zip")

	logger.Info(
		"Ensuring binary is available to download",
		zap.String("version", version),
		zap.String("artifact", artifactName),
	)
	url, err := downloadURL(version, repository, artifactName)
	if err != nil {
		return fmt.Errorf("failed to download vega binary: expected version of vega not found: %w", err)
	}

	logger.Info("Downloading vega binary", zap.String("url", url))

	c := http.Client{
		Timeout: time.Second * 120,
	}
	resp, err := c.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get vega binary release from %s: %w", url, err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get vega binary release with bad status: %q", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}

	logger.Info("Extracting vega binary from archive")
	zReader, err := zip.NewReader(bytes.NewReader(body), resp.ContentLength)
	if err != nil {
		return fmt.Errorf("failed to unzip vega package: %w", err)
	}

	file, err := zReader.Open("vega")
	if err != nil {
		return fmt.Errorf("failed to get vega binary from unzipped folder: %w", err)
	}
	defer file.Close()

	if err := os.RemoveAll(destinationFile); err != nil {
		return fmt.Errorf("failed to remove existing destination file: %w", err)
	}

	binFile, err := os.Create(destinationFile)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer binFile.Close()

	err = os.Chmod(destinationFile, 0o755)
	if err != nil {
		return fmt.Errorf("failed to change permission for file %q: %w", destinationFile, err)
	}

	if _, err := io.Copy(binFile, file); err != nil {
		return fmt.Errorf("failed to copy content to file %q: %w", destinationFile, err)
	}

	defer func() {
		if err != nil {
			_ = os.Remove(binFile.Name())
		}
	}()

	logger.Info("Vega binary was successfully extracted")

	return nil
}

func runDownloadBinary(
	logger *zap.Logger,
	destination, coreRESTEndpoint, repository string,
) error {
	logger.Info("Getting version from the core rest api", zap.String("endpoint", coreRESTEndpoint))

	stats, err := getStatistics(coreRESTEndpoint)
	if err != nil {
		return fmt.Errorf("failed to get statistics for endpoint(%s): %w", coreRESTEndpoint, err)
	}
	logger.Info("Vega version got correctly", zap.String("version", stats.Statistics.AppVersion))

	if err := downloadVegaBinary(logger, repository, stats.Statistics.AppVersion, destination); err != nil {
		return fmt.Errorf("failed to download vega version: %w", err)
	}

	stdOut, err := tools.ExecuteBinary(destination, []string{"version"}, nil)
	if err != nil {
		return fmt.Errorf("downloaded vega binary looks broken: %s, %w", stdOut, err)
	}

	logger.Info(fmt.Sprintf("Vega version: %s", stdOut))

	return nil
}
