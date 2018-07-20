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
bos_bin=${bos_bin:-bos}
burrow_bin=${burrow_bin:-burrow}

# currently we must use 'solc' as hardcoded by compilers
solc_bin=solc

# If false we will not try to start Burrow and expect them to be running
boot=${boot:-true}
debug=${debug:-false}

test_exit=0
script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

if [[ "$debug" = true ]]; then
    set -o xtrace
fi

# ----------------------------------------------------------
# Constants

# Ports etc must match those in burrow.toml
grpc_port=20997
tendermint_port=36656
chain_dir="$script_dir/chain"
burrow_root="$chain_dir/.burrow"

# Temporary logs
burrow_log=burrow.log
#
# ----------------------------------------------------------

# ---------------------------------------------------------------------------
# Needed functionality

goto_base(){
  cd ${script_dir}/jobs_fixtures
}

pubkey_of() {
    jq -r ".Accounts | map(select(.Name == \"$1\"))[0].PublicKey.PublicKey" chain/genesis.json
}

address_of() {
    jq -r ".Accounts | map(select(.Name == \"$1\"))[0].Address" chain/genesis.json
}

test_setup(){
  echo "Setting up..."
  cd "$script_dir"

  echo
  echo "Using binaries:"
  echo "  $(type ${solc_bin}) (version: $(${solc_bin} --version))"
  echo "  $(type ${bos_bin}) (version: $(${bos_bin} version))"
  echo "  $(type ${burrow_bin}) (version: $(${burrow_bin} --version))"
  echo
  # start test chain
  if [[ "$boot" = true ]]; then
    echo "Starting Burrow with tendermint port: $tendermint_port, GRPC port: $grpc_port"
    rm -rf ${burrow_root}
    $(cd "$chain_dir" && ${burrow_bin} start -v0 2> "$burrow_log")&
    burrow_pid=$!

  else
    echo "Not booting Burrow, but expecting Burrow to be running with tm RPC on port $grpc_port"
  fi

  key1_addr=$(address_of "Full_0")
  key2_addr=$(address_of "Participant_0")
  key2_pub=$(pubkey_of "Participant_0")

  echo -e "Default Key =>\t\t\t\t$key1_addr"
  echo -e "Backup Key =>\t\t\t\t$key2_addr"
  sleep 4 # boot time

  echo "Setup complete"
  echo ""
}

run_test(){
  # Run the jobs test
  echo ""
  echo -e "Testing $bos_bin jobs using fixture =>\t$1"
  goto_base
  cd $1
  echo
  cat readme.md
  echo

  echo ${bos_bin} --chain-url="localhost:$grpc_port" --address "$key1_addr" \
    --set "addr1=$key1_addr" --set "addr2=$key2_addr" --set "addr2_pub=$key2_pub" #--debug
  ${bos_bin} --chain-url="localhost:$grpc_port" --address "$key1_addr" \
    --set "addr1=$key1_addr" --set "addr2=$key2_addr" --set "addr2_pub=$key2_pub" #--debug
  test_exit=$?

  git clean -fdx ../**/bin ./jobs_output.csv
  rm ./*.output.json

  # Reset for next run
  goto_base
  return $test_exit
}

perform_tests(){
  echo ""
  goto_base
  apps=($1*/)
  echo $apps
  repeats=${2:-1}
  # Useful for soak testing/generating background requests to trigger concurrency issues
  for rep in `seq ${repeats}`
  do
    for app in "${apps[@]}"
    do
      echo "Test: $app, Repeat: $rep"
      run_test ${app}
      # Set exit code properly
      test_exit=$?
      if [ ${test_exit} -ne 0 ]
      then
        break
      fi
    done
  done
}

perform_tests_that_should_fail(){
  echo ""
  goto_base
  apps=($1*/)
  for app in "${apps[@]}"
  do
    run_test ${app}

    # Set exit code properly
    test_exit=$?
    if [ ${test_exit} -ne 0 ]
    then
      # actually, this test is meant to pass
      test_exit=0
    else
      break
    fi
  done
}

test_teardown(){
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
    [[ "$boot" = true ]] && rm -f "$burrow_log"
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
test_setup

# ---------------------------------------------------------------------------
# Go!

if [[ "$1" != "setup" ]]
then
  # Cleanup
  trap test_teardown EXIT
  if ! [ -z "$1" ]
  then
    echo "Running tests beginning with $1..."
    perform_tests "$1" "$2"
  else
    echo "Running tests that should fail"
    perform_tests_that_should_fail expected-failure

    echo "Running tests that should pass"
    perform_tests app
  fi
fi
