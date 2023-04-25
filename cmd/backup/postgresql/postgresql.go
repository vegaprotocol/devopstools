package postgresql

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	_ "github.com/lib/pq"
	"github.com/vegaprotocol/devopstools/tools"
	"go.uber.org/zap/buffer"
	"gopkg.in/ini.v1"
)

type Config struct {
	DataDirectory string `ini:"data_directory"`
}

func ReadConfig(location string) (*Config, error) {
	resultConfig := &Config{}

	if _, err := os.Stat(location); err != nil {
		return resultConfig, fmt.Errorf("failed to check if postgresql config(%s) exists: %w", location, err)
	}

	configData, err := ini.Load(location)
	if err != nil {
		return resultConfig, fmt.Errorf("failed to read postgresql file: %w", err)
	}

	if err := configData.MapTo(&resultConfig); err != nil {
		return resultConfig, fmt.Errorf("failed to unmarshal postgresql config: %w", err)
	}

	return resultConfig, nil
}

func IgnoreConfigParams(location string, params []string, ignoreMissingFile bool) (int, error) {
	if !tools.FileExists(location) {
		if ignoreMissingFile {
			return 0, nil
		}

		return 0, fmt.Errorf("the %s file does not exists", location)
	}

	var newContent buffer.Buffer

	file, err := os.Open(location)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	changedLines := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		matched := false
		for _, param := range params {
			re, err := regexp.Compile(fmt.Sprintf(`^\s*%s\s*=.*$`, param))
			if err != nil {
				return changedLines, fmt.Errorf("failed to compile regex for the %s param", param)
			}

			if re.MatchString(line) {
				changedLines += 1
				matched = true
				if _, err := newContent.WriteString(fmt.Sprintf("# %s\n", line)); err != nil {
					return changedLines, fmt.Errorf("failed to append commented line to new content: %w", err)
				}
				break
			}
		}

		if !matched {
			if _, err := newContent.WriteString(fmt.Sprintf("%s\n", line)); err != nil {
				return changedLines, fmt.Errorf("failed to append line to new content: %w", err)
			}
		}
	}

	if err := os.WriteFile(location, newContent.Bytes(), os.ModePerm); err != nil {
		return changedLines, fmt.Errorf("failed to update the %s file: %w", location, err)
	}

	return changedLines, nil
}
