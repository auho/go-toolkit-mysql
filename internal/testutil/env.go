package testutil

import "github.com/auho/go-toolkit-testutil"

// LoadEnv loads environment variables from the .env.test file at the project
// root. It is a thin wrapper around testutil.LoadEnv; see the upstream package
// for details.
func LoadEnv() error {
	return testutil.LoadEnv()
}
