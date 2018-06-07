#!/usr/bin/env bash

script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

address_of() {
    jq -r ".Accounts | map(select(.Name == \"$1\"))[0].Address" genesis.json
}

full_addr=$(address_of "Full_0")

burrow keys export --addr ${full_addr} | jq '{address: .address, pubKey: .pub_key[1], privKey: .priv_key[1]}'
