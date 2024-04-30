#!/usr/bin/env bash 

GENESIS=$HOME/.evoluteaid/config/genesis.json
evoluteaiD_BIN=$(which evoluteaid)

ALICE_ADDRESS=$($evoluteaiD_BIN keys show alice | head -n 1 | cut -f 2 -d ":" | tr -d " ")
if [[ ${#ALICE_ADDRESS} -ne 43 ]] || [[ $ALICE_ADDRESS != allo* ]]; then
    echo "Alice address not found"
    exit 1
fi

BOB_ADDRESS=$($evoluteaiD_BIN keys show bob | head -n 1 | cut -f 2 -d ":" | tr -d " ")
if [[ ${#BOB_ADDRESS} -ne 43 ]] || [[ $BOB_ADDRESS != allo* ]]; then
    echo "Bob address not found"
    exit 1
fi

VALIDATOR_ADDRESS=$(cat $GENESIS | grep "validator_address" | cut -f 2 -d ":" | tr -d " " | tr -d "\"" | tr -d ",")
if [[ ${#VALIDATOR_ADDRESS} -ne 50 ]] || [[ $VALIDATOR_ADDRESS != allovaloper* ]]; then
    echo "Validator address not found in genesis file"
    exit 1
fi
