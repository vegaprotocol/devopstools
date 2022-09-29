package claimcodes

import (
	"fmt"
)

type ClaimCodesVersion string

const (
	ClaimCodesV1 ClaimCodesVersion = "v1"
)

func (n ClaimCodesVersion) IsValid() error {
	switch n {
	case ClaimCodesV1:
		return nil
	}
	return fmt.Errorf("Invalid Claim Codes Version %s", n)
}
