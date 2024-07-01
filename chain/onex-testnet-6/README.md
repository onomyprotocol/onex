# Onex consumer chain

## Detail

The `onex-testnet-6` chain will be launched as a consumer chain in Onomy testnet.

- Network information: https://github.com/decentrio/onex/tree/dev/chain/onex-testnet-6
- Chain ID: `onex-testnet-6`
* Spawn time: `July 1st, 2024, 15:00 UTC`
* Genesis file (without CCV): https://github.com/onomyprotocol/onex/blob/dev/chain/onex-testnet-6/genesis_without-ccv.json
* Genesis with CCV: (To be generated)
- Current version: `v1.2.2`
* Binary: 
   * Version: [v1.2.2](https://github.com/onomyprotocol/onex/releases/tag/v1.2.2)
   * SHA256: `823bcf074b1ac5a4136a9a80b4a90c30dd19c5d01371da4043137370c16515f5`
* Onex GitHub repository: https://github.com/onomyprotocol/onex
- Peers: `27c4033e76a7b51d9ffa20de2dc5b12776332509@207.211.188.215:26656, a51d1ac5f0db3359ac8e673a4f33f50ccae20d4e@207.211.187.77:26656, 0ac3dcee8d3b7b1819807e39fcc4a578146016a6@64.71.153.55:26756, 9828e61e909b51d6d01fe4134241f39c128c6aff@222.106.187.14:53100, 13153783ee0ee5f7b787d5ed2f0ce43b0da696f4@180.131.222.73:26756`
- Endpoints: 
    - RPC: `To be added`
    - API: `To be added`
    - gRPC: `To be added`
- Block Explorer: ``

## IBC detail
| | onex-testnet-6 | onomy-testnet-2 |
|-------------|---------------------|-----------------|
|Client |``| ``|
|Connections | `` | `` |
|Channels | `transfer`: `` <br/><br/> `consumer`: `` | `transfer`: `` <br/><br/> `provider`: `` |

## Setup Instruction

### 1. Joining Onomy provider chain (onomy-testnet-2) as a validator
First, validators need to run the Onomy provider chain. To set up the node and join the network, please follow the instructions in [testnet documentation](https://github.com/onomyprotocol/validator/blob/main/testnet/readme.md).

Here is the detail of the Onomy provider chain:
- Chain ID: `onomy-testnet-2`
- Version: [v1.1.4](https://github.com/onomyprotocol/onomy/releases/tag/v1.1.4)
- Genesis: https://raw.githubusercontent.com/onomyprotocol/validator/main/testnet/genesis/genesis-testnet-6.json
- Seeds: 
```

```
- Persistent Peers:
```

```


### 4. Setup Onex consumer chain
The validators also need to set up the `onex-testnet-6` consumer chain. Here are the commands to install the binary and set up the new chain.
```bash
# detail of setup will appear here
cd $HOME/go/bin
wget -O onexd https://github.com/onomyprotocol/onex/releases/download/v1.2.2/onexd && chmod +x onexd
onexd version # v1.2.2
onexd init <moniker> --chain-id onex-testnet-6
cd $HOME/.onomy_onex/
wget -O config/genesis.json https://raw.githubusercontent.com/onomyprotocol/onex/dev/chain/onex-testnet-6/genesis_without-ccv.json
```

The validators **MUST NOT** run the node but wait until the new genesis is published on the Onomy repository, which will be detailed in step **[5. Vote the consumer-addition proposal](#5-vote-the-consumer-addition-proposal)**.

### 5. Vote on the consumer-addition proposal
The proposal to launch `onex-testnet-6` as a consumer chain will be submitted on the Onomy provider testnet and the validators should participate in voting for the proposal. After the proposal is passed, the validators should wait until the `spawn_time` and replace the old genesis file with the new `genesis.json` file from the Onomy repository.

```bash
wget -O /$HOME/.onomy_onex/config/genesis.json https://raw.githubusercontent.com/onomyprotocol/onex/dev/chain/onex-testnet-6/genesis.json
```

### 6. Wait for genesis and run

At the genesis time, validators can start the consumer chain by running
```bash
onexd start
```

> Note: if validators choose to run onex and onomy in the same machine, it is highly recommended to setup separate ports to prevent clashing. These ports are: P2P, RPC, REST, gRPC, gRPC-web

The validators can also use service to run and monitor the node. Here is the example of `/etc/systemd/system/onex.service`:
```
[Unit]
Description=Onex node
After=network.target

[Service]
ExecStart=/$HOME/go/bin/onexd start --p2p.persistent_peers="To be added"
Restart=always
RestartSec=3
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
```

After that, run these commands to enable and start the chain:
```bash
systemctl daemon-reload
systemctl enable onex.service
systemctl restart onex.service
```
and run `journalctl -fu onex -n150` to check the log. 

## Launch Stages
|Step|When?|What do you need to do?|What is happening?|
|----|--------------------------------------------------|----------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------|
|1   |ASAP                                              |Join the Onomy testnet `onomy-testnet-2`  as a full node and sync to the tip of the chain.|Validator machines getting caught up on existing Composable chain's history                                                                         |
|2   | Consumer Addition proposal on provider chain | [PROVIDER] Optional: Vote for the consumer-addition proposal.  | The proposals that provide new details for the launch.                            |
|3   |The proposals passed                                 |Nothing                                                                           | The proposals passed, `spawn_time` is set. After `spawn_time` is reached, the `ccv.json` file containing `ccv` state will be provided from the provider chain.
|4   |`spawn_time` reached                                  |The `genesis.json` file will be provided in the testnets repo. Replace the old `genesis.json` in the `$HOME/.onomy_onex/config` directory with the new `genesis.json`. The new `genesis.json` file with ccv data will be published in [onomyprotocol/onex](https://github.com/onomyprotocol/onex/tree/dev/chain/onex-testnet-6) |
|5   |Genesis reached     | Start your node with the consumer binary | onex-testnet-6 chain will start and become a consumer chain.                                                                                     |
|6   |3 blocks after upgrade height                     |Celebrate! :tada: 🥂                                                |<chain> blocks are now produced by the provider validator set|
