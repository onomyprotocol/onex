#Use Ubuntu Latest

#Setting up constants
ONEX_HOME=$HOME/.onomy_onex
ONEX_SRC=$ONEX_HOME/src/onex
COSMOVISOR_SRC=$ONEX_HOME/src/cosmovisor

ONEX_VERSION="v1.0.0-dev"
NODE_EXPORTER_VERSION="0.18.1"
COSMOVISOR_VERSION="cosmovisor-v1.0.1"

mkdir -p $ONEX_HOME
mkdir -p $ONEX_HOME/bin
mkdir -p $ONEX_HOME/contracts
mkdir -p $ONEX_HOME/logs
mkdir -p $ONEX_HOME/cosmovisor/genesis/bin
mkdir -p $ONEX_HOME/cosmovisor/upgrades/

echo "-----------installing dependencies---------------"
sudo apt install build-essential crudini jq

echo "----------------------installing onomy---------------"
curl -LO https://github.com/onomyprotocol/onex/releases/download/$ONEX_VERSION/onomyd
mv onomyd $ONEX_HOME/cosmovisor/genesis/bin/onexd

echo "----------------------installing cosmovisor---------------"
curl -LO https://github.com/onomyprotocol/onomy-sdk/releases/download/$COSMOVISOR_VERSION/cosmovisor
mv cosmovisor $ONEX_HOME/bin/cosmovisor

# echo "----------------installing eth bridge gbt-------------"
# curl -LO https://github.com/onomyprotocol/arc/releases/download/$ETH_BRIDGE_VERSION/gbt
# mv gbt $ONEX_HOME/bin/gbt

echo "-------------------installing node_exporter-----------------------"
curl -LO "https://github.com/prometheus/node_exporter/releases/download/v$NODE_EXPORTER_VERSION/node_exporter-$NODE_EXPORTER_VERSION.linux-amd64.tar.gz"
tar -xvf "node_exporter-$NODE_EXPORTER_VERSION.linux-amd64.tar.gz"
mv "node_exporter-$NODE_EXPORTER_VERSION.linux-amd64/node_exporter" $ONEX_HOME/bin/node_exporter
rm -r "node_exporter-$NODE_EXPORTER_VERSION.linux-amd64"
rm "node_exporter-$NODE_EXPORTER_VERSION.linux-amd64.tar.gz"

echo "-------------------adding binaries to path-----------------------"
chmod +x $ONEX_HOME/bin/*
export PATH=$PATH:$ONEX_HOME/bin
chmod +x $ONEX_HOME/cosmovisor/genesis/bin/*
export PATH=$PATH:$ONEX_HOME/cosmovisor/genesis/bin

echo "export PATH=$PATH" >> ~/.profile

# set the cosmovisor environments
echo "export DAEMON_HOME=$ONEX_HOME/" >> ~/.profile
echo "export DAEMON_NAME=onomyd" >> ~/.profile
echo "export DAEMON_RESTART_AFTER_UPGRADE=true" >> ~/.profile

source $HOME/.profile
ulimit -S -n 65536

echo "Onex binaries are installed successfully."
