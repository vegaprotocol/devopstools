package multisigcontrol

import (
	"fmt"
)

type Version string

const (
	V1 Version = "v1"
	V2 Version = "v2"
)

func (n Version) IsValid() error {
	switch n {
	case V1, V2:
		return nil
	}
	return fmt.Errorf("invalid Multisig Control Version %s", n)
}
