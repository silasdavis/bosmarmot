#!/usr/bin/env bash
set -e
# The version of solc we will fetch and install into ./bin/ for integration testsS
SOLC_URL="https://github.com/ethereum/solidity/releases/download/v0.4.20/solc-static-linux"
SOLC_BIN="$1"

wget -O "$SOLC_BIN" "$SOLC_URL"

chmod +x "$SOLC_BIN"
