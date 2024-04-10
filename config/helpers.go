package config

func ListDatanodeGRPCEndpoints(cfg Config) []string {
	var endpoints []string
	for _, node := range cfg.Nodes {
		endpoints = append(endpoints, node.API.DataNodeGRPCURL)
	}

	return endpoints
}
