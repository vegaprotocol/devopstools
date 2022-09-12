package veganetwork

import (
	"fmt"

	"github.com/vegaprotocol/devopstools/tools"
)

func (network *VegaNetwork) GetNetworkNodes() []string {
	switch network.Name {
	case "mainnet":
		return []string{"mainnet-observer.ops.vega.xyz"}
	}
	hosts := []string{}
	previousMissing := false
	for i := 0; i < 100; i++ {
		host := fmt.Sprintf("n%02d.%s", i, network.DNSSuffix)
		if _, err := tools.GetIP(host); err != nil {
			if previousMissing {
				break
			} else {
				previousMissing = true
			}
		} else {
			hosts = append(hosts, host)
		}
	}
	return hosts
}
