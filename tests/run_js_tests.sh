#!/usr/bin/env bash
# ----------------------------------------------------------
# PURPOSE

# This is the test manager for monax jobs. It will run the testing
# sequence for monax jobs referencing test fixtures in this tests directory.

# ----------------------------------------------------------
# REQUIREMENTS

# m

# ----------------------------------------------------------
# USAGE

# run_pkgs_tests.sh [appXX]

# Various required binaries locations can be provided by wrapper
burrow_bin=${burrow_bin:-burrow}

# If false we will not try to start Burrow and expect them to be running
boot=${boot:-true}
debug=${debug:-false}

test_exit=0

script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
chain_dir="${script_dir}/chain"
log_dir="${script_dir}/logs"
mkdir -p ${log_dir}
js_dir="${script_dir}/../legacy-contracts.js"

if [[ "$debug" = true ]]; then
    set -o xtrace
fi
set -e

# ----------------------------------------------------------
# Constants

# Ports etc must match those in burrow.toml
keys_port=48002
rpc_tm_port=48003
burrow_root="${chain_dir}/.burrow"

# Temporary logs
burrow_log="${log_dir}/burrow.log"
#
# ----------------------------------------------------------

# ---------------------------------------------------------------------------
# Needed functionality

account_data(){
  test_account=$(jq -r "." ${chain_dir}/account.json)
}

test_setup(){
  echo "Setting up..."

  echo
  echo "Using binaries:"
  echo "  $(type ${burrow_bin}) (version: $(${burrow_bin} --version))"
  echo

  # start test chain
  if [[ "$boot" = true ]]; then
    echo "Starting Burrow with tendermint RPC port: $rpc_tm_port"
    rm -rf ${burrow_root}

    ${burrow_bin} start -c "${chain_dir}/burrow.toml" -g "${chain_dir}/genesis.json" 2> "$burrow_log" &
    burrow_pid=$!

    sleep 1

  else
    echo "Not booting Burrow, but expecting Burrow to be running with tm RPC on port $rpc_tm_port"
  fi

  account_data
  sleep 4 # boot time

  echo "Setup complete"
  echo ""
}


perform_tests(){
  cd "$js_dir"

  vector=privateKey_vector account=$test_account mocha --recursive --reporter mocha-circleci-reporter ${1}
  vector=address_vector SIGNBYADDRESS=true account=$test_account mocha --recursive --reporter mocha-circleci-reporter ${1}

  test_exit=$?
}

test_teardown(){
  cd "$script_dir"
  echo "Cleaning up..."
  if [[ "$boot" = true ]]; then
    kill ${burrow_pid}
    echo "Waiting for burrow to shutdown..."
    wait ${burrow_pid} 2> /dev/null &
    rm -rf "$burrow_root"
  fi
  echo ""
  if [[ "$test_exit" -eq 0 ]]
  then
    [[ "$boot" = true ]] && rm -rf "$log_dir"
    echo "Tests complete! Tests are Green. :)"
  else
    echo "Tests complete. Tests are Red. :("
    echo "Failure in: $app"
  fi
  exit ${test_exit}
}

# ---------------------------------------------------------------------------
# Setup


echo "Hello! I'm the marmot that tests the $bos_bin jobs tooling."
echo
echo "testing with target $bos_bin"
echo


if [[ "$TEST" == "record" ]] || [[ "$TEST" == "server" ]]; then
    test_setup
    trap test_teardown EXIT
fi

echo "Running js Tests..."
perform_tests "$1"

