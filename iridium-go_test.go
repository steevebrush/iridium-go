package iridium_go

import (
	"github.com/steevebrush/iridium-go/iridiumWalletdRPC"
	"github.com/steevebrush/iridium-go/iridiumdRPC"
	"strconv"
	"testing"
)

func TestIridiumdRPCVersion(t *testing.T) {
	name, major, minor, patch := iridiumdRPC.Version()
	version := strconv.Itoa(major) + "." + strconv.Itoa(minor) + "." + strconv.Itoa(patch)
	t.Logf("Found %q v%s", name, version)
}

func TestIridiumWalletRPCVersion(t *testing.T) {
	name, major, minor, patch := iridiumWalletdRPC.Version()
	version := strconv.Itoa(major) + "." + strconv.Itoa(minor) + "." + strconv.Itoa(patch)
	t.Logf("Found %q v%s", name, version)
}
