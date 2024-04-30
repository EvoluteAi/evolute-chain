#!/bin/bash
set -exu

NETWORK="${NETWORK:-testnet}"
GENESIS_URL="https://raw.githubusercontent.com/evoluteai-network/networks/main/${NETWORK}/genesis.json"
SEEDS_URL="https://raw.githubusercontent.com/evoluteai-network/networks/main/${NETWORK}/seeds.txt"
BLOCKLESS_API_URL="${BLOCKLESS_API_URL:-https://heads.${NETWORK}.evoluteai.network}"               #! Replace with your blockless API URL

APP_HOME="${APP_HOME:-/data}"
INIT_FLAG="${APP_HOME}/.initialized"
MONIKER="${MONIKER:-$(hostname)}"
KEYRING_BACKEND=test                              #! Use test for simplicity, you should decide which backend to use !!!
GENESIS_FILE="${APP_HOME}/config/genesis.json"
DENOM="uallo"

echo "To re-initiate the node, remove the file: ${INIT_FLAG}"
if [ ! -f $INIT_FLAG ]; then
    rm -rf ${APP_HOME}/config

    #* Init node
    evoluteaid --home=${APP_HOME} init ${MONIKER} --chain-id=${NETWORK} --default-denom $DENOM

    #* Download genesis
    rm -f $GENESIS_FILE
    curl -Lo $GENESIS_FILE $GENESIS_URL

    #* Import evoluteai account, priv_validator_key.json and node_key.json from the vault here
    #* Here create a new evoluteaid account
    evoluteaid --home $APP_HOME keys add ${MONIKER} --keyring-backend $KEYRING_BACKEND > $APP_HOME/${MONIKER}.account_info 2>&1

    #* Adjust configs
    #* Enable prometheus metrics
    #dasel put -t bool -v true 'instrumentation.prometheus' -f ${APP_HOME}/config/config.toml

    #* Setup evoluteaid client
    evoluteaid --home=${APP_HOME} config set client chain-id ${NETWORK}
    evoluteaid --home=${APP_HOME} config set client keyring-backend $KEYRING_BACKEND

    #* Create symlink for evoluteaid config
    ln -sf . ${APP_HOME}/.evoluteaid

    touch $INIT_FLAG
fi
echo "Node is initialized"

SEEDS=$(curl -s ${SEEDS_URL})

echo "Starting validator node"
evoluteaid \
    --home=${APP_HOME} \
    start \
    --moniker=${MONIKER} \
    --minimum-gas-prices=0${DENOM} \
    --rpc.laddr=tcp://0.0.0.0:26657 \
    --p2p.seeds=$SEEDS

