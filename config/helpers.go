package config

func ListDatanodeGRPCEndpoints(cfg Config) []string {
	var endpoints []string
	for _, node := range cfg.Nodes {
		if node.API.DataNodeGRPCURL == "" {
			continue
		}
		endpoints = append(endpoints, node.API.DataNodeGRPCURL)
	}

	return endpoints
}

func ListDatanodeRESTEndpoints(cfg Config) []string {
	var endpoints []string
	for _, node := range cfg.Nodes {
		if node.API.DataNodeRESTURL == "" {
			continue
		}
		endpoints = append(endpoints, node.API.DataNodeRESTURL)
	}

	return endpoints
}

func ListBlockchainRESTEndpoints(cfg Config) []string {
	var endpoints []string
	for _, node := range cfg.Nodes {
		if node.API.BlockchainRESTURL == "" {
			continue
		}
		endpoints = append(endpoints, node.API.BlockchainRESTURL)
	}

	return endpoints
}

func ListValidatorRESTEndpoints(cfg Config) []string {
	var endpoints []string
	for _, node := range cfg.Nodes {
		if node.API.VegaRESTURL == "" {
			continue
		}
		endpoints = append(endpoints, node.API.VegaRESTURL)
	}

	return endpoints
}
