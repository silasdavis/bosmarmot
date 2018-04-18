#!/usr/bin/env bash
set -e
# The version of solc we will fetch and install into ./bin/ for integration testsS
# Our custom build of solidity fixing linking issue; https://github.com/monax/solidity/tree/truncate-beginning-v0.4.22
SOLC_URL="https://drive.google.com/uc?export=download&id=1_RHMTxlYHL4zcmZDyXemccR8ENVeAEO2"
SOLC_BIN="$1"

wget -O "$SOLC_BIN" "$SOLC_URL"

chmod +x "$SOLC_BIN"
