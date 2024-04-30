#!/bin/bash

BINARY=./evoluteaid
CHAIN_DIR=./data
CHAINID=evoluteai_demo

echo "Starting $CHAINID in $CHAIN_DIR..."
echo "Creating log file at $CHAIN_DIR/$CHAINID.log"
$BINARY start --log_level trace --log_format json --home $CHAIN_DIR/$CHAINID --minimum-gas-prices="1uallo" --pruning=nothing > $CHAIN_DIR/$CHAINID.log 2>&1 &
