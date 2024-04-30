# Steps to init the graph indexer from scratch

## Create home

```bash
mkdir $HOME/firehose
mkdir $HOME/firehose/onex-bins
```

## Prepare onex binaries with firehose  (figment-networks)  required dependencies.

In the original onex repo the `github.com/tendermint/tendermint` is replaced by `github.com/cometbft/cometbft`.

But for the firehose we need different go.mod replacement:

Example:

```
github.com/tendermint/tendermint => github.com/figment-networks/tendermint v0.34.28-fh
```

Full list of tags is [here](https://github.com/graphprotocol/tendermint/tags).

For each release/upgrade we need the corresponding binary with such go.mod replacement.

### Example on how to download onex binary with curl:

```bash
export ONEXD_VERSION=v1.0.0-dev-fh
curl -L -o $HOME/firehose/onex-bins/onexd-$ONEXD_VERSION https://github.com/onomyprotocol/onex/releases/download/$ONEXD_VERSION/onexd
chmod +x $HOME/firehose/onex-bins/onexd-$ONEXD_VERSION
```

## Init onex chain

### Move genesis onexd to firehose home

```bash
export GENESIS_ONEXD_VERSION=v1.0.0-dev-fh
cp $HOME/firehose/onex-bins/onexd-$GENESIS_ONEXD_VERSION $HOME/firehose/onexd
# check 
$HOME/firehose/onexd version
```

### Init default onexd config

```bash
$HOME/firehose/onexd init $(hostname)
```

### Replace default genesis

```bash
rm $HOME/.onomy_onex/config/genesis.json
curl -L -o $HOME/.onomy_onex/config/genesis.json https://raw.githubusercontent.com/onomyprotocol/onex/dev/chain/onex-testnet-5/genesis.json
```

### Set seeds

```bash
export ONEX_SEEDS="a2be48320ead4280e644107aa1536d94be235e9f@65.109.69.90:2030,2f96d16645fd52dba217fb477a66c7b637fbb3c7@64.71.153.55:26756,e6e0a2fef354c509f31d573305626cc2a5cc9982@64.71.153.54:26756,f80867e8181a07b26a17e4f597b0cfb7408b1b2a@180.131.222.73:26756,eb823e14ff73127ccce3e17bd674046b290416f1@51.250.106.107:36656"
sed -i -e "s/seeds = \"\"/seeds = \"$ONEX_SEEDS\"/g" $HOME/.onomy_onex/config/config.toml
```

### Update node config

```bash
cat << END >> $HOME/.onomy_onex/config/config.toml

#######################################################
###       Extractor Configuration Options     ###
#######################################################
[extractor]
enabled = true
output_file = "stdout"
END
```

## Set-up firehose

### Download binary

```bash
export FIREHOSE_VERSION=v0.7.1
curl -L -o $HOME/firehose/firehose https://github.com/graphprotocol/firehose-cosmos/releases/download/$FIREHOSE_VERSION/firecosmos_linux_amd64 
chmod +x $HOME/firehose/firehose
# check 
$HOME/firehose/firehose help
```

### Create firehose config

```bash
cat << END >> $HOME/firehose/firehose.yml
start:
  args:
    - reader
    - merger
    - firehose
  flags:
    common-first-streamable-block: 1
    common-live-blocks-addr:
    reader-mode: node
    reader-node-path: $HOME/firehose/onexd
    reader-node-args: start --x-crisis-skip-assert-invariants
    reader-node-logs-filter: "module=(p2p|pex|consensus|x/bank)"
    relayer-max-source-latency: 99999h
    common-live-blocks-addr: localhost:15011
    relayer-grpc-listen-addr: :15011
    merger-time-between-store-pruning: 10s
    verbose: 1
END
```

### Check that firehose is set up correctly

#### Run to check setup

```bash
$HOME/firehose/firehose start -c $HOME/firehose/firehose.yml --data-dir $HOME/firehose/fh-data
```

#### In different terminal call (install grpcurl if not installed)

```bash
grpcurl -plaintext 0.0.0.0:9030 sf.firehose.v2.Stream/Blocks
```

If you see the blocks are printing, stop the `firehose`.

#### Creare linux service

##### Set up journald

```bash
nano /etc/systemd/journald.conf
```

And set

```bash
# [Journal]
Storage=persistent
```

Restart the systemd-journald

```bash
systemctl restart systemd-journald
```

##### Run linux service

```bash
echo "
[Unit]
After=network.target

[Service]
User=root
ExecStart=$HOME/firehose/firehose start -c $HOME/firehose/firehose.yml --data-dir $HOME/firehose/fh-data
Restart=always
RestartSec=3
LimitAS=infinity
LimitRSS=infinity
LimitCORE=infinity
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
" > "/etc/systemd/system/firehose.service"

systemctl daemon-reload
systemctl start firehose
systemctl enable firehose
systemctl status firehose --no-pager
journalctl -a -f
```

### Chain upgrade

The onex binary running by the firehose will stop running once the chan upgrade block is reached.
At that time update the `onexd` binary to the upgrade binary in the `$HOME/firehose/onexd ` and restart the
`firehose` service.

### Stop service

```bash
systemctl stop firehose
```

### Backup the firehouse data

Back up data in `$HOME/firehose/fh-data` folder.

### Replace binary

```bash
export NEW_ONEXD_VERSION=v1.0.1-dev-fh # the version here depends on your upgrade
# expect the binary with the upgrade in the `$HOME/firehose/onex-bins`
# if not download it using the "Example on how to download onex binary with curl" in this doc.
rm $HOME/firehose/onexd
cp $HOME/firehose/onex-bins/onexd-$NEW_ONEXD_VERSION $HOME/firehose/onexd
# validate that it works
$HOME/firehose/onexd version
```

### Start service

```bash
systemctl start firehose
journalctl -a -f
```

## Create and prepare the graph env

### Init files

```bash
# !!! Update the values before use !!!
export POSTGRES_USER="you postgres username"
export POSTGRES_PASSWORD="you postgres password"
```

```bash

mkdir $HOME/firehose/compose

cat << END >> $HOME/firehose/compose/graph-node-config.toml

[deployment]
[[deployment.rule]]
shard = "primary"
indexers = [ "default" ]

[store]
[store.primary]
connection = "postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@0.0.0.0:5432/graph-node"
pool_size = 10

[chains]
ingestor = "block_ingestor_node"

[chains.onex-testnet-5]
shard = "primary"
protocol = "cosmos"
provider = [
  { label = "onex-testnet-5", details = { type = "firehose", url = "http://0.0.0.0:9030" }},
]
END


cat << END >> $HOME/firehose/compose/docker-compose.yml
version: '3'
services:
  graph-node:
    image: graphprotocol/graph-node:v0.35.0
    depends_on:
      - ipfs
      - postgres
    network_mode: "host"
    environment:
      postgres_host: postgres
      postgres_user: graph-node
      postgres_pass: let-me-in
      postgres_db: graph-node
      ipfs: '0.0.0.0:5001'
      ethereum: 'mainnet:http://0.0.0.0:8545'
      GRAPH_LOG: info
      GRAPH_NODE_CONFIG: /etc/config.toml
    restart: always
    volumes:
      - ./graph-node-config.toml:/etc/config.toml

  ipfs:
    image: ipfs/kubo:v0.17.0
    network_mode: "host"
    restart: always
    volumes:
      - ./data/ipfs:/data/ipfs:Z

  postgres:
    image: postgres:14-alpine
    network_mode: "host"
    command:
      [
        "postgres",
        "-cshared_preload_libraries=pg_stat_statements",
        "-cmax_connections=200"
      ]
    environment:
      POSTGRES_USER: $POSTGRES_USER
      POSTGRES_PASSWORD: $POSTGRES_PASSWORD
      POSTGRES_DB: graph-node
      PGDATA: "/var/lib/postgresql/data"
      POSTGRES_INITDB_ARGS: "-E UTF8 --locale=C"
    restart: always
    volumes:
      - ./data/postgres:/var/lib/postgresql/data:Z
END
```

### Run docker-compose

```bash
docker compose up -d
```

### Check logs of the

```bash
docker compose logs graph-node
```

There should be no errors.

Proceed with  [app-example](app-example)
