package iridiumWalletdRPC

import (
	"strconv"
	"testing"
)

func TestVersion(t *testing.T) {
	name, major, minor, patch := Version()
	version := strconv.Itoa(major) + "." + strconv.Itoa(minor) + "." + strconv.Itoa(patch)
	t.Logf("Package %s v%s found", name, version)
}
