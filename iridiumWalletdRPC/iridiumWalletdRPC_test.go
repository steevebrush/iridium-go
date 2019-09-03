/*
 * Copyright (c) 2019.
 * by Steve Brush, Iridium Developers
 */

// Iridium payments gateway JSON RPC API for golang test
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
