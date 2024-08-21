package stakingbridge

type Version string

const (
	V1 Version = "v1"
	V2 Version = "v2"
)

func StdVersion(version string) Version {
	if version == string(V2) {
		return V2
	}

	return V1
}
