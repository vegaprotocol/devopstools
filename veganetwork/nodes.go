package veganetwork

import "fmt"

func (network *VegaNetwork) GetNetworkNodes() []string {
	switch network.Name {
	case "mainnet":
		return []string{"mainnet-observer.ops.vega.xyz"}
	}
	hosts := make([]string, 21)
	for i := 0; i < 21; i++ {
		hosts[i] = fmt.Sprintf("n%02d.%s", i, network.DNSSuffix)
	}
	return hosts
}
