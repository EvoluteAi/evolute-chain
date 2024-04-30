#!/usr/bin/env bash 

# This script is a quick and dirty integration test to ensure validator rewards are paid out

set -e

source $(dirname $0)/common.sh

# this script expects to be ran AFTER `scripts/init.sh`
if ! test -f $GENESIS; then
  echo "Must run scripts/init.sh first."
  exit 1
fi

evoluteai_REWARDS_ADDRESS=$($evoluteaiD_BIN query auth module-account "evoluteairewards" | grep "address: allo" | cut -f 2 -d ":" | tr -d " " | tr -d "\"")
if [[ ${#evoluteai_REWARDS_ADDRESS} -ne 43 ]] || [[ $evoluteai_REWARDS_ADDRESS != allo* ]]; then
    echo "evoluteai rewards address not found"
    exit 1
fi

# get the current outstanding rewards
DISTRIBUTION_REWARDS_0=$($evoluteaiD_BIN query distribution validator-outstanding-rewards $VALIDATOR_ADDRESS | grep "amount" | cut -f 2 -d ":" | tr -d " " | tr -d "\"")
evoluteai_REWARDS_0=$($evoluteaiD_BIN query bank balances $evoluteai_REWARDS_ADDRESS | grep "amount" | cut -f 2 -d ":" | tr -d " " | tr -d "\"")

# wait for some blocks to get mined
sleep 5

# get the current outstanding rewards
DISTRIBUTION_REWARDS_1=$($evoluteaiD_BIN query distribution validator-outstanding-rewards $VALIDATOR_ADDRESS | grep "amount" | cut -f 2 -d ":" | tr -d " " | tr -d "\"")
evoluteai_REWARDS_1=$($evoluteaiD_BIN query bank balances $evoluteai_REWARDS_ADDRESS | grep "amount" | cut -f 2 -d ":" | tr -d " " | tr -d "\"")

# assert that the rewards have increased
DISTRIBUTION_INCREASED=$(bc <<< "$DISTRIBUTION_REWARDS_1 > $DISTRIBUTION_REWARDS_0")
evoluteai_REWARDS_INCREASED=$(bc <<< "$evoluteai_REWARDS_1 > $evoluteai_REWARDS_0")

if [[ $DISTRIBUTION_INCREASED -ne 1 ]]; then
    echo "Distribution of rewards to validators did not increase"
    exit 1
else 
    echo "Distribution of rewards to validators increased"
fi

if [[ $evoluteai_REWARDS_INCREASED -ne 1 ]]; then
    echo "Distribution of rewards to evoluteai rewards did not increase"
    exit 1
else 
    echo "Distribution of rewards to evoluteai rewards increased"
fi
