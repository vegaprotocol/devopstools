package multisigcontrol

import (
	"fmt"
)

type MultisigControlVersion string

const (
	MultisigControlV1 MultisigControlVersion = "v1"
	MultisigControlV2 MultisigControlVersion = "v2"
)

func (n MultisigControlVersion) IsValid() error {
	switch n {
	case MultisigControlV1, MultisigControlV2:
		return nil
	}
	return fmt.Errorf("Invalid Multisig Control Version %s", n)
}
