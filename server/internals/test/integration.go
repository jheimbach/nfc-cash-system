package test

import (
	"os"
	"testing"
)

// IsIntegrationTest will check if envvar RUN_INTEGRATION is non zero.
// if RUN_INTEGRATION is not set or zero, tests that call IsIntegrationTest will be skipped
func IsIntegrationTest(t *testing.T) {
	if EnvWithDefault("RUN_INTEGRATION", "0") == "0" {
		t.Skipf("skipping integration tests")
	}
}

// genEnvWithDefault is a helper function to get envvar or default value
func EnvWithDefault(envName, defaultVal string) string {
	if val, ok := os.LookupEnv(envName); ok {
		return os.ExpandEnv(val)
	}
	return defaultVal
}
