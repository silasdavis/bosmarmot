#!/usr/bin/env bash
set -e
# The version of solc we will fetch and install into ./bin/ for integration testsS
# Our custom build of solidity fixing linking issue; https://github.com/monax/solidity/tree/truncate-beginning-v0.4.22
SOLC_URL="https://pkgs.monax.io/apt/solc/0.4.24/linux-amd64/solc"

SOLC_BIN="$1"

wget -O "$SOLC_BIN" "$SOLC_URL"

chmod +x "$SOLC_BIN"
