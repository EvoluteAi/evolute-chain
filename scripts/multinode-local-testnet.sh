#!/bin/bash
set -e

# Clean up any existing data
rm -rf $HOME/.evoluteaid/

# Create four evoluteaid directories for validators
mkdir -p $HOME/.evoluteaid/validator1
mkdir -p $HOME/.evoluteaid/validator2
mkdir -p $HOME/.evoluteaid/validator3
mkdir -p $HOME/.evoluteaid/validator4

# Initialize all four validators
evoluteaid init test --chain-id=demo --default-denom uallo --home=$HOME/.evoluteaid/validator1
evoluteaid init test --chain-id=demo --default-denom uallo --home=$HOME/.evoluteaid/validator2
evoluteaid init test --chain-id=demo --default-denom uallo --home=$HOME/.evoluteaid/validator3
evoluteaid init test --chain-id=demo --default-denom uallo --home=$HOME/.evoluteaid/validator4

# Set up keyring-backend and chain-id for all validators
for i in 1 2 3 4; do
    evoluteaid config --home $HOME/.evoluteaid/validator$i set client keyring-backend test
    evoluteaid config --home $HOME/.evoluteaid/validator$i set client chain-id demo
done

# Create keys for all validators
evoluteaid keys add validator1 --keyring-backend test --home $HOME/.evoluteaid/validator1
evoluteaid keys add validator2 --keyring-backend test --home $HOME/.evoluteaid/validator2
evoluteaid keys add validator3 --keyring-backend test --home $HOME/.evoluteaid/validator3
evoluteaid keys add validator4 --keyring-backend test --home $HOME/.evoluteaid/validator4

# Add genesis accounts for each validator
evoluteaid genesis add-genesis-account $(evoluteaid keys show validator1 -a --keyring-backend test --home=$HOME/.evoluteaid/validator1) 10000000allo --home $HOME/.evoluteaid/validator1
evoluteaid genesis add-genesis-account $(evoluteaid keys show validator2 -a --keyring-backend test --home=$HOME/.evoluteaid/validator2) 10000000allo --home $HOME/.evoluteaid/validator2
evoluteaid genesis add-genesis-account $(evoluteaid keys show validator3 -a --keyring-backend test --home=$HOME/.evoluteaid/validator3) 10000000allo --home $HOME/.evoluteaid/validator3
evoluteaid genesis add-genesis-account $(evoluteaid keys show validator4 -a --keyring-backend test --home=$HOME/.evoluteaid/validator4) 10000000allo --home $HOME/.evoluteaid/validator4

# Create a gentx for each validator
evoluteaid genesis gentx validator1 1000allo --chain-id demo --keyring-backend test --home=$HOME/.evoluteaid/validator1
evoluteaid genesis gentx validator2 1000allo --chain-id demo --keyring-backend test --home=$HOME/.evoluteaid/validator2
evoluteaid genesis gentx validator3 1000allo --chain-id demo --keyring-backend test --home=$HOME/.evoluteaid/validator3
evoluteaid genesis gentx validator4 1000allo --chain-id demo --keyring-backend test --home=$HOME/.evoluteaid/validator4

# Collect gentxs to the first validator's genesis file
evoluteaid genesis collect-gentxs --home=$HOME/.evoluteaid/validator1

# Copy the genesis file from the first validator to others
cp $HOME/.evoluteaid/validator1/config/genesis.json $HOME/.evoluteaid/validator2/config/
cp $HOME/.evoluteaid/validator1/config/genesis.json $HOME/.evoluteaid/validator3/config/
cp $HOME/.evoluteaid/validator1/config/genesis.json $HOME/.evoluteaid/validator4/config/


# Update validator1
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $HOME/.evoluteaid/validator1/config/config.toml
sed -i -E 's|prometheus = false|prometheus = true|g' $HOME/.evoluteaid/validator1/config/config.toml

# Update port configurations for validator2
sed -i -E 's|tcp://127.0.0.1:26658|tcp://127.0.0.1:26655|g' $HOME/.evoluteaid/validator2/config/config.toml # P2P
sed -i -E 's|tcp://127.0.0.1:26657|tcp://127.0.0.1:26654|g' $HOME/.evoluteaid/validator2/config/config.toml # RPC
sed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26653|g' $HOME/.evoluteaid/validator2/config/config.toml # pprof
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $HOME/.evoluteaid/validator2/config/config.toml
sed -i -E 's|prometheus = false|prometheus = true|g' $HOME/.evoluteaid/validator2/config/config.toml

# Update port configurations for validator3
sed -i -E 's|tcp://127.0.0.1:26658|tcp://127.0.0.1:26652|g' $HOME/.evoluteaid/validator3/config/config.toml # P2P
sed -i -E 's|tcp://127.0.0.1:26657|tcp://127.0.0.1:26651|g' $HOME/.evoluteaid/validator3/config/config.toml # RPC
sed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26650|g' $HOME/.evoluteaid/validator3/config/config.toml # pprof
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $HOME/.evoluteaid/validator3/config/config.toml
sed -i -E 's|prometheus = false|prometheus = true|g' $HOME/.evoluteaid/validator3/config/config.toml

# Update port configurations for validator4
sed -i -E 's|tcp://127.0.0.1:26658|tcp://127.0.0.1:26649|g' $HOME/.evoluteaid/validator4/config/config.toml # P2P
sed -i -E 's|tcp://127.0.0.1:26657|tcp://127.0.0.1:26648|g' $HOME/.evoluteaid/validator4/config/config.toml # RPC
sed -i -E 's|tcp://0.0.0.0:26656|tcp://0.0.0.0:26647|g' $HOME/.evoluteaid/validator4/config/config.toml # pprof
sed -i -E 's|allow_duplicate_ip = false|allow_duplicate_ip = true|g' $HOME/.evoluteaid/validator4/config/config.toml
sed -i -E 's|prometheus = false|prometheus = true|g' $HOME/.evoluteaid/validator4/config/config.toml

# Get the node ID of validator1
VALIDATOR1_NODE_ID=$(evoluteaid tendermint show-node-id --home $HOME/.evoluteaid/validator1)

# Configure validator2, validator3, and validator4 to have validator1 as a persistent peer
sed -i -E "s|persistent_peers = \"\"|persistent_peers = \"${VALIDATOR1_NODE_ID}@localhost:26656\"|g" $HOME/.evoluteaid/validator2/config/config.toml
sed -i -E "s|persistent_peers = \"\"|persistent_peers = \"${VALIDATOR1_NODE_ID}@localhost:26656\"|g" $HOME/.evoluteaid/validator3/config/config.toml
sed -i -E "s|persistent_peers = \"\"|persistent_peers = \"${VALIDATOR1_NODE_ID}@localhost:26656\"|g" $HOME/.evoluteaid/validator4/config/config.toml

# Start each validator in separate tmux sessions
tmux new -s validator1 -d evoluteaid start --home=$HOME/.evoluteaid/validator1
tmux new -s validator2 -d evoluteaid start --home=$HOME/.evoluteaid/validator2
tmux new -s validator3 -d evoluteaid start --home=$HOME/.evoluteaid/validator3
tmux new -s validator4 -d evoluteaid start --home=$HOME/.evoluteaid/validator4