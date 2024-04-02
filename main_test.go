package main

import (
	"os/user"
	"testing"

	"github.com/mitchellh/go-homedir"
)

// TestMain : check to see if user exist
func TestMain(t *testing.T) {

	t.Run("User should exist",
		func(t *testing.T) {
			_, errCurr := user.Current()
			if errCurr != nil {
				t.Errorf("Unable to get user %v [unexpected]", errCurr)
			}
			_, errCurr = homedir.Dir()
			if errCurr != nil {
				t.Errorf("Unable to get user home directory: %v [unexpected]", errCurr)
			}
		},
	)
}
