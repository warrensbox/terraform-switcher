package main_test

import (
	"os/user"
	"testing"
)

// TestMain : check to see if user exist
func TestMain(t *testing.T) {

	t.Run("User should exist",
		func(t *testing.T) {
			usr, errCurr := user.Current()
			if errCurr != nil {
				t.Errorf("Unable to get user %v [unexpected]", errCurr)
			}

			if usr != nil {
				t.Logf("Current user exist: %v  [expected]\n", usr.HomeDir)
			} else {
				t.Error("Unable to get user [unexpected]")
			}
		},
	)
}
