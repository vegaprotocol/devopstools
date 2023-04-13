package tools

import (
	"fmt"
	"time"
)

func Retry(retries int, sleepAfterFailure time.Duration, action func() error) error {
	finalErr := ""
	for i := 0; i < retries; i++ {
		err := action()
		if err != nil {
			finalErr = fmt.Sprintf("%s\nError #%d: %s", finalErr, i+1, err.Error())
			time.Sleep(sleepAfterFailure)
			continue
		}

		// Successfull
		return nil
	}

	if finalErr != "" {
		return fmt.Errorf(finalErr)
	}

	return nil
}
