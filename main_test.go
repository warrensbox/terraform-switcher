package main_test

import (
	"os"
	"os/user"
	"testing"
)

// TestMain : check to see if user exist
func TestMain(t *testing.T) {

	t.Run("User should exist",
		func(t *testing.T) {
			_, errCurr := user.Current()
			if errCurr != nil {
				t.Errorf("Unable to get user %v [unexpected]", errCurr)
			}
			_, errCurr = os.UserHomeDir()
			if errCurr != nil {
				t.Errorf("Unable to get user home directory: %v [unexpected]", errCurr)
			}
		},
	)
}
