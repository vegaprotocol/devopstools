package multisigcontrol

import (
	"fmt"
)

type MultisigControlVersion string

const (
	V1 MultisigControlVersion = "v1"
	V2 MultisigControlVersion = "v2"
)

func (n MultisigControlVersion) IsValid() error {
	switch n {
	case V1, V2:
		return nil
	}
	return fmt.Errorf("Invalid Multisig Control Version %s", n)
}
