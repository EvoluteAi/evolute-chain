#!/usr/bin/env bash

set -e

source $(dirname $0)/common.sh

NEXT_TOPIC_ID=$($evoluteaiD_BIN query emissions next-topic-id | head -n 1 | cut -f 2 -d ":" | tr -d " " | tr -d "\"")
if [[ $(bc <<< "$NEXT_TOPIC_ID > 1") -ne 1 ]]; then
  echo "NEXT_TOPIC_ID not incremented, must run topic.sh first"
  exit 1
fi

ALICE_REGD_0=$($evoluteaiD_BIN query emissions registered-topic-ids "$ALICE_ADDRESS" "true")
if [[ "$ALICE_REGD_0" != "{}" ]] ; then
    echo "Alice already registered as a reputer or some error";
    exit 1
fi

BOB_REGD_0=$($evoluteaiD_BIN query emissions registered-topic-ids "$BOB_ADDRESS" "false")
if [[ "$BOB_REGD_0" != "{}" ]] ; then
    echo "Bob already registered as a worker or some error";
    exit 1
fi


echo "Register Alice as a reputer in topic 1"
REGR_CREATOR="$ALICE_ADDRESS"
REGR_LIBP2P_KEY="reputerkey"
REGR_MULTI_ADDRESS="reputermultiaddress"
REGR_TOPIC_IDS="1"
REGR_INITIAL_STAKE="100000"
REGR_OWNER="$ALICE_ADDRESS"
REGR_IS_REPUTER="true"
$evoluteaiD_BIN tx emissions register \
    $REGR_CREATOR \
    $REGR_LIBP2P_KEY \
    $REGR_MULTI_ADDRESS \
    $REGR_TOPIC_IDS \
    $REGR_INITIAL_STAKE \
    $REGR_OWNER \
    $REGR_IS_REPUTER \
    --yes --keyring-backend=test --chain-id=demo \
    --gas-prices=1uallo --gas=auto --gas-adjustment=1.5;


echo "Checking that Alice now shows up as registered"
ALICE_REGD=false
for COUNT_SLEEP in 1 2 3 4 5
do
  ALICE_REGD_1=$($evoluteaiD_BIN query emissions registered-topic-ids "$ALICE_ADDRESS" "true" | tail -n 2 | cut -f 2 -d "-" | tr -d " " | tr -d "\"")
  if [[ "$ALICE_REGD_1" != "1" ]] ; then
      echo "Alice not registered as a reputer in topic 1";
      COUNT_SLEEP=$((COUNT_SLEEP+1))
      sleep 1
  else
      echo "Alice successfully registered as a reputer in topic 1";
      ALICE_REGD=true
      break 
  fi
done
if [ "$ALICE_REGD" = false ]; then
  echo "The network has not registered Alice."
  exit 1
fi

echo "Register Bob as a worker in topic 1"
REGR_CREATOR="$BOB_ADDRESS"
REGR_LIBP2P_KEY="workerkey"
REGR_MULTI_ADDRESS="workermultiaddress"
REGR_TOPIC_IDS="1"
REGR_INITIAL_STAKE="100000"
REGR_OWNER="$BOB_ADDRESS"
REGR_IS_REPUTER="false"
$evoluteaiD_BIN tx emissions register \
    $REGR_CREATOR \
    $REGR_LIBP2P_KEY \
    $REGR_MULTI_ADDRESS \
    $REGR_TOPIC_IDS \
    $REGR_INITIAL_STAKE \
    $REGR_OWNER \
    $REGR_IS_REPUTER \
    --yes --keyring-backend=test --chain-id=demo \
    --gas-prices=1uallo --gas=auto --gas-adjustment=1.5;



echo "Checking that Bob now shows up as registered"
BOB_REGD=false
for COUNT_SLEEP in 1 2 3 4 5
do
  BOB_REGD_1=$($evoluteaiD_BIN query emissions registered-topic-ids "$BOB_ADDRESS" "false" | tail -n 2 | cut -f 2 -d "-" | tr -d " " | tr -d "\"")
  if [[ "$BOB_REGD_1" != "1" ]] ; then
      echo "Bob not registered as a worker in topic 1";
      COUNT_SLEEP=$((COUNT_SLEEP+1))
      sleep 1
  else
      echo "Bob successfully registered as a worker in topic 1";
      BOB_REGD=true
      break 
  fi
done
if [ "$BOB_REGD" = false ]; then
  echo "The network has not registered Bob."
  exit 1
fi

echo "Setting epoch length to one block"
$evoluteaiD_BIN tx emissions update-params \
    $ALICE_ADDRESS \
    '{"epoch_length": ["1"], "max_inference_request_validity":[],"max_missing_inference_percent":[],"max_request_cadence":[],"max_topics_per_block":[],"min_request_cadence":[],"min_request_unmet_demand":[],"min_topic_unmet_demand":[],"min_weight_cadence":[],"percent_rewards_reputers_workers":[],"remove_stake_delay_window":[],"required_minimum_stake":[],"version":[] }' \
    --yes --keyring-backend=test --chain-id=demo \
    --gas-prices=1uallo --gas=auto --gas-adjustment=1.5;
sleep 5


echo "Setting weights from alice"
WEIGHT_0=$($evoluteaiD_BIN query emissions weight 1 "$ALICE_ADDRESS" "$BOB_ADDRESS" | cut -f 2 -d ":" | tr -d " " | tr -d "\"")
if [[ "$WEIGHT_0" != "0" ]] ; then
    echo "Weight already set by alice on bob or some error";
    exit 1
fi

$evoluteaiD_BIN tx emissions set-weights\
    "$ALICE_ADDRESS" \
    "{\"topic_id\": 1, \"reputer\": \"$ALICE_ADDRESS\", \"worker\": \"$BOB_ADDRESS\", \"weight\": \"1000\"}" \
    --keyring-backend=test --chain-id=demo --yes \
    --gas-prices=1uallo --gas=auto --gas-adjustment=1.5;

echo "Checking that weights are set"
WEIGHT_SET=false
for COUNT_SLEEP in 1 2 3 4 5
do
  WEIGHT_1=$($evoluteaiD_BIN query emissions weight 1 "$ALICE_ADDRESS" "$BOB_ADDRESS" | cut -f 2 -d ":" | tr -d " " | tr -d "\"")
  if [[ "$WEIGHT_1" != "1000" ]] ; then
      echo "Weights not set";
      COUNT_SLEEP=$((COUNT_SLEEP+1))
      sleep 1
  else
      echo "Weights successfully set";
      WEIGHT_SET=true
      break 
  fi
done
if [ "$WEIGHT_SET" = false ]; then
  echo "The network has not set the weights as expected."
  exit 1
fi

echo "Checking that the staking module is getting paid via inflation via the reward module"
evoluteai_STAKING_ADDRESS=$($evoluteaiD_BIN query auth module-account "evoluteaistaking" | grep "address: allo" | cut -f 2 -d ":" | tr -d " " | tr -d "\"")
if [[ ${#evoluteai_STAKING_ADDRESS} -ne 43 ]] || [[ $evoluteai_STAKING_ADDRESS != allo* ]]; then
    echo "evoluteai rewards address not found"
    exit 1
fi

evoluteai_STAKING_0=$($evoluteaiD_BIN query bank balances $evoluteai_STAKING_ADDRESS | grep "amount" | cut -f 2 -d ":" | tr -d " " | tr -d "\"")
evoluteai_STAKE_SUCCEED=false
for COUNT_SLEEP in 1 2 3 4 5 6 7 8 9 10
do
    evoluteai_STAKING_1=$($evoluteaiD_BIN query bank balances $evoluteai_STAKING_ADDRESS | grep "amount" | cut -f 2 -d ":" | tr -d " " | tr -d "\"")
    evoluteai_STAKING_INCREASED=$(bc <<< "$evoluteai_STAKING_1 > $evoluteai_STAKING_0")
    if [[ $evoluteai_STAKING_INCREASED -ne 1 ]]; then
        echo "Distribution of rewards to evoluteai staking did not increase"
        COUNT_SLEEP=$((COUNT_SLEEP+1))
        sleep 1
    else 
        echo "Distribution of rewards to evoluteai staking increased"
        evoluteai_STAKE_SUCCEED=true
    fi
    sleep 1
done
if [ "$evoluteai_STAKE_SUCCEED" = false ]; then
    echo "The network has not distributed rewards to evoluteai staking as expected."
    exit 1
fi

echo "Checking that bob's stake is going up due to having non-zero weights"
BOB_STAKE_POSITION_0=$($evoluteaiD_BIN query emissions account-stake-list "$BOB_ADDRESS" | grep "amount" | cut -f 2 -d ":" | tr -d " " | tr -d "\"")
BOB_STAKE_SUCCEED=false
for COUNT_SLEEP in 1 2 3 4 5 6 7 8 9 10
do
    BOB_STAKE_POSITION_1=$($evoluteaiD_BIN query emissions account-stake-list "$BOB_ADDRESS" | grep "amount" | cut -f 2 -d ":" | tr -d " " | tr -d "\"")
    BOB_STAKE_POSITION_INCREASED=$(bc <<< "$BOB_STAKE_POSITION_1 > $BOB_STAKE_POSITION_0")
    if [[ $BOB_STAKE_POSITION_INCREASED -ne 1 ]]; then
        echo "Bob did not get rewards for staking"
        COUNT_SLEEP=$((COUNT_SLEEP+1))
        sleep 1
    else 
        echo "Bob got rewards for staking"
        BOB_STAKE_SUCCEED=true
    fi
done
if [ "$BOB_STAKE_SUCCEED" = false ]; then
    echo "The network has not distributed rewards to bob as expected."
    exit 1
fi

echo "Rewards checks complete"
